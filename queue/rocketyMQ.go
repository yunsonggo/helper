package queue

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/yunsonggo/helper/algorithms"
	"github.com/yunsonggo/helper/types"
	"time"
)

type Rocket struct {
	Hosts         []string
	Auth          string
	Pass          string
	Broker        string
	RetryTimes    int
	BatchCount    int
	Topics        []string
	Groups        []string
	Producers     map[string]rocketmq.Producer
	PushConsumers map[string]map[string]rocketmq.PushConsumer
}

func NewRocket(conf *types.Mq, topics, groups []string) (*Rocket, error) {
	r := &Rocket{
		Hosts:         conf.Host,
		Auth:          conf.User,
		Pass:          conf.Pass,
		Broker:        conf.Broker,
		RetryTimes:    conf.RetryTimes,
		BatchCount:    conf.BatchCount,
		Topics:        topics,
		Groups:        groups,
		Producers:     make(map[string]rocketmq.Producer),
		PushConsumers: make(map[string]map[string]rocketmq.PushConsumer),
	}
	if len(topics) > 0 {
		if err := r.CheckTopics(context.Background()); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func (rs *Rocket) GroupList(ctx context.Context) (*admin.SubscriptionGroupWrapper, error) {
	adminClient, err := admin.NewAdmin(
		admin.WithResolver(primitive.NewPassthroughResolver(rs.Hosts)),
		admin.WithCredentials(primitive.Credentials{
			AccessKey: rs.Auth,
			SecretKey: rs.Pass,
		}),
	)
	if err != nil {
		return nil, err
	}
	defer adminClient.Close()
	return adminClient.GetAllSubscriptionGroup(ctx, rs.Broker, 30*time.Second)
}

func (rs *Rocket) TopicList(ctx context.Context) (*admin.TopicList, error) {
	adminClient, err := admin.NewAdmin(
		admin.WithResolver(primitive.NewPassthroughResolver(rs.Hosts)),
		admin.WithCredentials(primitive.Credentials{
			AccessKey: rs.Auth,
			SecretKey: rs.Pass,
		}),
	)
	if err != nil {
		return nil, err
	}
	defer adminClient.Close()
	return adminClient.FetchAllTopicList(ctx)
}

func (rs *Rocket) CreateTopic(ctx context.Context, topicName string) error {
	adminClient, err := admin.NewAdmin(
		admin.WithResolver(primitive.NewPassthroughResolver(rs.Hosts)),
		admin.WithCredentials(primitive.Credentials{
			AccessKey: rs.Auth,
			SecretKey: rs.Pass,
		}),
	)
	if err != nil {
		return err
	}
	defer adminClient.Close()
	return adminClient.CreateTopic(
		ctx,
		admin.WithTopicCreate(topicName),
		admin.WithBrokerAddrCreate(rs.Broker),
	)
}

func (rs *Rocket) CheckTopics(ctx context.Context) error {
	if len(rs.Topics) == 0 {
		return nil
	}

	list, err := rs.TopicList(ctx)
	if err != nil {
		return err
	}

	newTopics := algorithms.SliceUnion(list.TopicList, rs.Topics)

	if len(newTopics) > 0 {
		for _, topic := range newTopics {
			if err = rs.CreateTopic(ctx, topic); err != nil {
				return err
			}
		}
	}

	return nil
}

func (rs *Rocket) AddTopic(topic string) error {
	rs.Topics = append(rs.Topics, topic)
	return rs.CheckTopics(context.Background())
}

func (rs *Rocket) NewProducer(group string) (rocketmq.Producer, error) {
	if rs.Producers == nil {
		rs.Producers = make(map[string]rocketmq.Producer)
	}

	if !algorithms.SliceExist(rs.Groups, group) {
		rs.Groups = append(rs.Groups, group)
	}

	if pro, ok := rs.Producers[group]; ok {
		return pro, nil
	}
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(rs.Hosts)),
		producer.WithCredentials(primitive.Credentials{AccessKey: rs.Auth, SecretKey: rs.Pass}),
		producer.WithGroupName(group),
		producer.WithRetry(rs.RetryTimes),
	)
	if err != nil {
		return nil, err
	}
	err = p.Start()
	if err != nil {
		return nil, err
	}
	rs.Producers[group] = p
	return p, nil
}

type PushConsumerFunc func(conf *types.Mq, topic, group, logLevel string) (rocketmq.PushConsumer, error)

func (rs *Rocket) NewPushConsumer(topic, group, logLevel string, fn PushConsumerFunc) error {
	if rs.PushConsumers == nil {
		rs.PushConsumers = make(map[string]map[string]rocketmq.PushConsumer)
		rs.PushConsumers[topic] = make(map[string]rocketmq.PushConsumer)
	}

	if rs.PushConsumers[topic] == nil {
		rs.PushConsumers[topic] = make(map[string]rocketmq.PushConsumer)
	}

	if consumers, ok := rs.PushConsumers[topic]; ok {
		if _, has := consumers[group]; has {
			return nil
		}
	}

	conf := types.Mq{
		Host:       rs.Hosts,
		User:       rs.Auth,
		Pass:       rs.Pass,
		Broker:     rs.Broker,
		BatchCount: rs.BatchCount,
		RetryTimes: rs.RetryTimes,
	}
	consumer, err := fn(&conf, topic, group, logLevel)
	if err != nil {
		return err
	}
	rs.PushConsumers[topic][group] = consumer
	return nil
}

type ConsumersData struct {
	Topic    string
	Group    string
	LogLevel string
	Fun      PushConsumerFunc
}

func (rs *Rocket) Start(consumers []ConsumersData) error {
	if len(rs.Groups) > 0 {
		for _, group := range rs.Groups {
			if _, err := rs.NewProducer(group); err != nil {
				return err
			}
		}
	}
	if len(consumers) > 0 {
		for _, consumer := range consumers {
			if err := rs.NewPushConsumer(consumer.Topic, consumer.Group, consumer.LogLevel, consumer.Fun); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rs *Rocket) Shutdown() {
	if len(rs.Producers) > 0 {
		for _, p := range rs.Producers {
			_ = p.Shutdown()
		}
	}
	if len(rs.PushConsumers) > 0 {
		for _, consumers := range rs.PushConsumers {
			if len(consumers) > 0 {
				for _, consumer := range consumers {
					_ = consumer.Shutdown()
				}
			}
		}
	}
}
