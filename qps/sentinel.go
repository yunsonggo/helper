package qps

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/yunsonggo/helper/free"
	"github.com/yunsonggo/helper/types"
)

func InitQps(logPath string) error {
	if logPath == "" {
		if current, err := free.CurrentPath(); err == nil {
			logPath = current
		} else {
			return err
		}
	}
	logsDir := logPath + "/logs"
	entry := config.NewDefaultConfig()
	entry.Sentinel.Log.Dir = logsDir
	if err := sentinel.InitWithConfig(entry); err != nil {
		return err
	}
	return nil
}

func InitRules(c *types.Qps) (*flow.Rule, error) {
	var ts flow.TokenCalculateStrategy
	var cb flow.ControlBehavior

	switch c.TokenCalculateStrategy {
	case "Direct":
		ts = flow.Direct
	default:
		ts = flow.WarmUp
	}

	switch c.ControlBehavior {
	case "Throttling":
		cb = flow.Throttling
	default:
		cb = flow.Reject
	}

	switch c.WarmOrQps {
	case "warm":
		return &flow.Rule{
			Resource:               c.Resource,
			TokenCalculateStrategy: ts,
			ControlBehavior:        cb,
			Threshold:              c.Threshold,
			WarmUpPeriodSec:        c.WarmUpPeriodSec,
			WarmUpColdFactor:       c.WarmUpColdFactor,
		}, nil
	default:
		return &flow.Rule{
			Resource:               c.Resource,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              c.Threshold,
			StatIntervalInMs:       c.StatIntervalMs,
		}, nil
	}
}

func LoadRules(flowRule []*flow.Rule) error {
	if _, err := flow.LoadRules(flowRule); err != nil {
		return err
	}
	return nil
}

func Entry(resource string) (*base.SentinelEntry, *base.BlockError, bool) {
	if e, b := sentinel.Entry(resource, sentinel.WithTrafficType(base.Inbound)); b != nil {
		return e, b, false
	} else {
		return e, b, true
	}
}
