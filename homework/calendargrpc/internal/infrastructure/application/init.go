package application

import (
	"flag"
	"log"
	"strings"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Init(defaultConfigPath string) {

	// using standard library "flag" package
	flag.String("config", defaultConfigPath, "path to configuration flag")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	// Reading configuration from file
	configPath := viper.GetString("config") // retrieve value from viper
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Couldn't read configuration file: %s", err.Error())
	}

	// Setting log parameters
	logger.Init(viper.GetString("log.level"), viper.GetString("log.file"))
}
