package main

import (
	"fmt"
	"log-etl/core"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("D:/go-workspace/workspace/src/log-etl/conf/")
	viper.SetConfigFile("D:/go-workspace/workspace/src/log-etl/conf/core.yaml") // name of config file (without extension)
	err := viper.ReadInConfig()                                                 // Find and read the config file
	if err != nil {                                                             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error when reading %s config file: %s\n", err))
	}
	core.NewEtlMain()

}
