package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	configFilePrefix              = "config"
	configFileType                = "json"
	configPath                    = "."
	aspNetCorePortVariable        = "ASPNETCORE_PORT"
	aspNetCoreEnvironmentVariable = "ASPNETCORE_ENVIRONMENT"
	defaultEnvironment            = "Development"
	defaultPort                   = 3001
)

var (
	okStatusCode   = 200
	badStatusCode  = 500
	delayInSeconds = 10.0
)

func main() {
	CreateLog()
	LogHeader("Health Check")
	LogApplicationStart()
	config, err := readConfig()
	if err != nil {
		panic(err)
	}

	endpointName := config.GetString("endpointName")
	portNumber := config.GetString(aspNetCorePortVariable)
	portString := ":" + portNumber

	http.HandleFunc(endpointName, upDown)
	LogContentService(portString)
	log.Fatal(http.ListenAndServe(portString, nil))
}

func readConfig() (*viper.Viper, error) {
	coreEnv := getAspNetCoreEnvironment()
	v := viper.New()
	v.SetConfigName(configFilePrefix)
	v.SetConfigType(configFileType)
	v.AddConfigPath(configPath)
	v.AutomaticEnv()
	v.SetDefault(aspNetCorePortVariable, defaultPort)
	v.SetDefault(aspNetCoreEnvironmentVariable, coreEnv)
	err := v.ReadInConfig()
	return v, err
}

func getAspNetCoreEnvironment() string {
	env := os.Getenv(aspNetCoreEnvironmentVariable)
	if env == "" {
		env = defaultEnvironment
	}
	return strings.ToLower(env)
}

func upDown(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	config, err := readConfig()
	if err != nil {
		panic(err)
	}
	downString := config.GetString("downString")
	upString := config.GetString("upString")
	duration := time.Now().Sub(started)
	if duration.Seconds() > delayInSeconds {
		w.WriteHeader(badStatusCode)
		formatString := fmt.Sprintf(downString, duration.Seconds())
		LogInfo("Bad Duration: " + formatString)
		_, err := w.Write([]byte(formatString))
		if err != nil {
		}
	} else {
		w.WriteHeader(okStatusCode)
		_, err := w.Write([]byte(upString))
		if err != nil {
			w.WriteHeader(badStatusCode)
			formatString := fmt.Sprintf(downString, err)
			LogInfo("Bad Request: " + formatString)
			LogError(err)
			_, err := w.Write([]byte(formatString))
			if err != nil {
			}
		}
	}
}
