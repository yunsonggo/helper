package clould

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/yunsonggo/helper/types"
	"io"
	"log"
	"strings"
)

type AliOss struct {
	Client     *oss.Client
	Bucket     *oss.Bucket
	BucketName string
	Marker     oss.Option // 分页数据
}

func NewAliOss(conf *types.Oss, checkBuckets bool) (*AliOss, error) {
	endpoint := fmt.Sprintf("https://%s", conf.Endpoint)
	accessKeyID := conf.AccessKeyID
	accessKeySecret := conf.AccessKeySecret
	bucketName := conf.BucketName
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}
	aliOss := &AliOss{
		Client:     client,
		BucketName: bucketName,
	}
	if checkBuckets {
		ok, err := aliOss.IsBucketExist()
		if err != nil {
			return nil, err
		}
		if !ok {
			if err = aliOss.CreateBucket(bucketName); err != nil {
				return nil, err
			}
		}
	}
	b, err := aliOss.NewBucket()
	if err != nil {
		log.Fatal(err)
	}
	aliOss.Bucket = b
	return aliOss, nil
}

func (o *AliOss) IsBucketExist() (bool, error) {
	return o.Client.IsBucketExist(o.BucketName)
}
func (o *AliOss) CreateBucket(bucketName string) error {
	return o.Client.CreateBucket(bucketName, oss.StorageClass(oss.StorageStandard), oss.ACL(oss.ACLPrivate), oss.RedundancyType(oss.RedundancyLRS))
}
func (o *AliOss) NewBucket() (*oss.Bucket, error) {
	return o.Client.Bucket(o.BucketName)
}
func (o *AliOss) ObjectName(fileName string) string {
	prefix := "./"
	suffix := ".json"
	if strings.HasPrefix(fileName, prefix) {
		fileName = strings.TrimPrefix(fileName, prefix)
	}
	lastDotIndex := strings.LastIndex(fileName, ".")
	if lastDotIndex < 0 {
		fileName = fileName + ".json"
	} else {
		actualSuffix := fileName[lastDotIndex:]
		if actualSuffix != suffix {
			fileName = fileName + ".json"
		}
	}
	return fileName
}
func (o *AliOss) IsExist(fileName string) (bool, error) {
	path := o.ObjectName(fileName)
	return o.Bucket.IsObjectExist(path)
}
func (o *AliOss) Put(fileName string, body []byte) error {
	path := o.ObjectName(fileName)
	fmt.Printf("path:%s\n", path)
	err := o.Bucket.PutObject(path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	return nil
}
func (o *AliOss) Get(fileName string) ([]byte, error) {
	path := o.ObjectName(fileName)
	body, err := o.Bucket.GetObject(path)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (o *AliOss) Remove(filename string) error {
	path := o.ObjectName(filename)
	return o.Bucket.DeleteObject(path)
}
func (o *AliOss) ObjectList(page, size int) ([]string, error) {
	if page < 1 {
		page = 1
	}
	if size > 100 {
		size = 100
	}
	if page == 1 {
		marker := oss.Marker("")
		o.Marker = marker
	}
	res, err := o.Bucket.ListObjects(oss.MaxKeys(size), o.Marker)
	if err != nil {
		return nil, err
	}
	o.Marker = oss.Marker(res.NextMarker)
	var keys []string
	for _, object := range res.Objects {
		keys = append(keys, object.Key)
	}
	return keys, nil
}
