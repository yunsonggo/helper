package read_test

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/yunsonggo/helper/read"
	"github.com/yunsonggo/helper/types"
	"gopkg.in/yaml.v3"
	"log"
	"testing"
)

type TestConfig struct {
	Title string `yaml:"title"`
}

func TestRead(t *testing.T) {
	filename := "nacos"
	fileChangeChan := make(chan *types.HotUpdateSetting)
	nacosChangeChan := make(chan *types.HotUpdateSetting)
	go func() {
		err := read.WithFile(filename, "yaml", []string{"./conf"}, fileChangeChan)
		if err != nil {
			log.Fatal(err)
		}
	}()

	var nacosConf types.Nacos

	go func() {
		for {
			select {
			case data := <-fileChangeChan:
				if data.Error != nil {
					fmt.Printf("fileChange error:%v\n", data.Error.Error())
				} else if data.Data == nil {
					fmt.Printf("fileChange nil")
				} else {
					if err := yaml.Unmarshal(data.Data, &nacosConf); err != nil {
						fmt.Printf("yaml error:%v\n", err.Error())
					} else {
						fmt.Printf("fileChange success! id:%s\n", data.MsgID)
						if readNacosErr := read.WithNacos(&nacosConf, vo.YAML, nacosChangeChan); readNacosErr != nil {
							fmt.Printf("read nacos error:%v\n", readNacosErr.Error())
						} else {
							fmt.Println("readNacos success")
						}
					}
				}
			default:
				break
			}
		}
	}()

	var testConf TestConfig

	for {
		select {
		case data := <-nacosChangeChan:
			fmt.Printf("nacos data:%+v\n", data)
			if err := yaml.Unmarshal(data.Data, &testConf); err != nil {
				fmt.Printf("yaml unmarshal error:%v\n", err)
			} else {
				fmt.Printf("test config:%v\n", testConf)
			}
		default:
			break
		}
	}
}
