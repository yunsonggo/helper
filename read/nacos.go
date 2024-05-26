package read

import (
	"encoding/json"
	"fmt"
	"github.com/yunsonggo/helper/v3/types"

	"github.com/nacos-group/nacos-sdk-go/clients"
	_ "github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func InNacos(nc *types.Nacos, app *types.App, configType vo.ConfigType, changeChan chan *types.HotUpdate) error {
	nacosServerConfig := []constant.ServerConfig{
		{
			IpAddr: nc.NacosAddr,
			Port:   nc.NacosPort,
		},
	}
	nacosClientConfig := constant.ClientConfig{
		NamespaceId:         nc.NamespaceId,
		TimeoutMs:           nc.TimeoutMs,
		NotLoadCacheAtStart: nc.LoadCache,
		LogDir:              nc.LogDir,
		CacheDir:            nc.CacheDir,
		LogLevel:            nc.LogLevel,
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": nacosServerConfig,
		"clientConfig":  nacosClientConfig,
	})
	if err != nil {
		return err
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: nc.DataId,
		Group:  nc.Group,
		Type:   configType,
	})
	err = json.Unmarshal([]byte(content), app)
	if err != nil {
		return err
	}
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: nc.DataId,
		Group:  nc.Group,
		OnChange: func(namespace, group, dataId, data string) {
			name := fmt.Sprintf("namespace:%s,group:%s,dataId:%s,data:%s\n", namespace, group, dataId, data)
			value := &types.HotUpdate{
				Name: name,
				App:  app,
			}
			changeChan <- value
		},
	})
	return err
}

func WithNacos(nc *types.Nacos, configType vo.ConfigType, changeChan chan *types.HotUpdateSetting) error {
	nacosServerConfig := []constant.ServerConfig{
		{
			IpAddr: nc.NacosAddr,
			Port:   nc.NacosPort,
		},
	}
	nacosClientConfig := constant.ClientConfig{
		NamespaceId:         nc.NamespaceId,
		Username:            nc.Username,
		Password:            nc.Password,
		TimeoutMs:           nc.TimeoutMs,
		NotLoadCacheAtStart: nc.LoadCache,
		LogDir:              nc.LogDir,
		CacheDir:            nc.CacheDir,
		LogLevel:            nc.LogLevel,
	}
	voConfig := vo.ConfigParam{
		DataId: nc.DataId,
		Group:  nc.Group,
		Type:   configType,
		OnChange: func(namespace, group, dataId, data string) {
			name := fmt.Sprintf("namespace:%s,group:%s,dataId:%s", namespace, group, dataId)
			value := &types.HotUpdateSetting{
				Name: name,
				Data: []byte(data),
			}
			changeChan <- value
		},
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": nacosServerConfig,
		"clientConfig":  nacosClientConfig,
	})
	if err != nil {
		return err
	}
	content, err := configClient.GetConfig(voConfig)
	if err != nil {
		return err
	}
	data := &types.HotUpdateSetting{
		Name: fmt.Sprintf("namespace:%s,group:%s,dataId:%s", nc.NamespaceId, nc.Group, nc.DataId),
		Data: []byte(content),
	}
	changeChan <- data

	_ = configClient.CancelListenConfig(voConfig)

	err = configClient.ListenConfig(voConfig)
	return err
}
