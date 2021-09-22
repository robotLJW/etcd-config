package config

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

var once sync.Once

var WatchConfigMsg *WatchConfig
var CalculateConfigMsg *CalculateConfig

func ReadWatchConfig(configPath string, configName string, configType string) {
	once.Do(func() {
		viper.SetConfigName(configName)
		viper.SetConfigType(configType)
		viper.AddConfigPath(configPath)
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Sprintf("error %v\n", err))
		}
		WatchConfigMsg = &WatchConfig{
			EndPoints:   viper.GetString("watchConfig.endpoints"),
			DialTimeout: viper.GetDuration("watchConfig.dialTimeout"),
			WatchKey:    viper.GetString("watchConfig.watchKey"),
			PostAddress: viper.GetString("watchConfig.postAddress"),
		}
	})
}

func ReadCalculateConfig(configPath string, configName string, configType string) {
	once.Do(func() {
		viper.SetConfigName(configName)
		viper.SetConfigType(configType)
		viper.AddConfigPath(configPath)
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Sprintf("error %v\n", err))
		}
		CalculateConfigMsg = &CalculateConfig{
			EndPoints:      viper.GetString("calculateConfig.endpoints"),
			DialTimeout:    viper.GetDuration("calculateConfig.dialTimeout"),
			UpdateInterval: viper.GetDuration("calculateConfig.updateInterval"),
			ChannelCount:   viper.GetInt("calculateConfig.channelCount"),
			DataKey:        viper.GetString("calculateConfig.dataKey"),
			DataValue:      viper.GetString("calculateConfig.dataValue"),
		}
	})
}
