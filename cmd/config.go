package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Config = config{Host: "localhost", Port: 80}

type config struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName(".ohctl")
	viper.AddConfigPath("$HOME")
	viper.SetEnvPrefix("ohctl")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("error reading config: ", err)
		os.Exit(1)
	}

	err := viper.Unmarshal(&Config)

	if err != nil {
		fmt.Println("error loading config: ", err)
		os.Exit(1)
	}
}
