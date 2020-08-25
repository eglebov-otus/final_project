package main

import (
	"errors"
	"fmt"
	"image-previewer/internal"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	EnvDevelopment = "dev"
	EnvProduction  = "prod"
)

func main() {
	if err := initConfig(); err != nil {
		panic(fmt.Sprintf("failed to init configuration: %s", err))
	}

	if err := initLogger(); err != nil {
		panic(fmt.Sprintf("failed to init logger: %s", err))
	}

	app := internal.NewApp()

	if err := app.Run(); err != nil {
		panic(fmt.Sprintf("failed to start application: %s", err))
	}
}

func initLogger() error {
	var logger *zap.Logger
	var err error

	switch viper.GetString("app.environment") {
	case EnvDevelopment:
		logger, err = zap.NewDevelopment()
	case EnvProduction:
		logger, err = zap.NewProduction()
	default:
		err = errors.New("unsupported env")
	}

	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)

	return nil
}

func initConfig() error {
	var configFile string

	pflag.StringVar(&configFile, "config", "./configs/config.yml", "path to config")
	pflag.Parse()
	viper.SetConfigFile(configFile)

	return viper.ReadInConfig()
}
