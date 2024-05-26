package read

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/yunsonggo/helper/free"
	"github.com/yunsonggo/helper/types"
	"gopkg.in/yaml.v3"
)

func InFile(fileName, fileType string, paths []string, app *types.App, changeChan chan *types.HotUpdate) error {
	viper.SetConfigType(fileType)
	viper.SetConfigName(fileName)
	viper.AddConfigPath(".")
	if len(paths) > 0 {
		for _, path := range paths {
			viper.AddConfigPath(path)
		}
	}
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(app); err != nil {
		return err
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		data := &types.HotUpdate{
			Name: e.Name,
			App:  app,
		}
		changeChan <- data
	})

	viper.WatchConfig()
	return nil
}

func WithFile(fileName, fileType string, paths []string, changeChan chan *types.HotUpdateSetting) error {
	viper.SetConfigType(fileType)
	viper.SetConfigName(fileName)
	viper.AddConfigPath(".")
	if len(paths) > 0 {
		for _, path := range paths {
			viper.AddConfigPath(path)
		}
	}
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	v := viper.GetViper()
	data, err := toByteSetting(v)
	if err != nil {
		return err
	}
	setting := &types.HotUpdateSetting{
		MsgID: free.RandomNowCode(6, "20060102_150405_", "_", true, []string{"nacos"}),
		Name:  fileName,
		Data:  data,
	}

	changeChan <- setting

	viper.OnConfigChange(func(e fsnotify.Event) {
		var newSetting *types.HotUpdateSetting
		if newData, newErr := toByteSetting(viper.GetViper()); err != nil {
			newSetting = &types.HotUpdateSetting{
				MsgID: free.RandomNowCode(6, "20060102_150405_", "_", true, []string{fileName}),
				Name:  e.Name,
				Error: newErr,
			}
		} else {
			newSetting = &types.HotUpdateSetting{
				MsgID: free.RandomNowCode(6, "20060102_150405_", "_", true, []string{fileName}),
				Name:  e.Name,
				Data:  newData,
			}
		}
		changeChan <- newSetting
	})
	viper.WatchConfig()
	return nil
}

func toByteSetting(v *viper.Viper) ([]byte, error) {
	return yaml.Marshal(v.AllSettings())
}
