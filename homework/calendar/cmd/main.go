package main

import (
	"flag"
	"log"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendar/internal/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/AcroManiac/otus-go/homework/calendar/internal/calendar"
)

func init() {
	// using standard library "flag" package
	flag.String("config", "../configs/calendar.yaml", "path to configuration flag")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// Reading configuration from file
	configPath := viper.GetString("config") // retrieve value from viper
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Couldn't read configuration file: %s", err.Error())
	}

	// Setting log parameters
	logger.Init(viper.GetString("log.log_level"), viper.GetString("log.log_file"))
}

func main() {
	// Create calendar
	var cal calendar.Calendar = calendar.NewCalendar()

	// Create and add event
	if _, err := cal.CreateEvent(time.Now(), time.Now().Add(time.Hour)); err != nil {
		logger.Error("Error adding event", "error", err.Error())
	}

	logger.Info("Calendar was created. Bye!")

	//log.Println(viper.GetString("http_listen.ip"))
}
