package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/AcroManiac/otus-go/homework/calendar/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

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
	_ = viper.BindPFlags(pflag.CommandLine)

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
	cal := calendar.NewCalendar(storage.NewStorage())

	// Create and add event
	if _, err := cal.CreateEvent(time.Now(), time.Now().Add(time.Hour)); err != nil {
		logger.Error("Error adding event", "error", err.Error())
	}
	logger.Info("Calendar was created")

	// Initialize and start HTTP server
	router := mux.NewRouter()
	router.HandleFunc("/hello", handlerHello).Methods("GET")

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("http_listen.ip"),
			viper.GetInt("http_listen.port")),
		Handler: router,
	}

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Error while starting HTTP server", "error", err)
		}
	}()
	logger.Info("HTTP server started")

	<-done
	logger.Info("HTTP server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Make HTTP server graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server shutdown failed", "error", err)
	}
	logger.Info("HTTP server exited properly")
}

func handlerHello(w http.ResponseWriter, r *http.Request) {
	// Log message params
	logger.Debug("Incoming message",
		"host", r.Host,
		"url", r.URL.Path)

	// Make JSON response to incoming http request
	response := HelloResponse{
		Message: "Hello world!",
	}
	jsonResponse(w, response, http.StatusOK)
}

// JSON representation of message
type HelloResponse struct {
	Message string `json:"message"`
}

// Obtained from http-boilerplate project:
// https://github.com/jordan-wright/http-boilerplate/blob/master/server/api/v1/api.go
func jsonResponse(w http.ResponseWriter, data interface{}, c int) {
	dj, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		logger.Error("Error creating JSON response", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	_, _ = fmt.Fprintf(w, "%s", dj)
}
