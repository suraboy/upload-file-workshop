package main

import (
	googleStorage "cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sendgrid/sendgrid-go"
	"github.com/spf13/viper"
	"github.com/suraboy/upload-file-worksho/app/core/config"
	"github.com/suraboy/upload-file-worksho/app/core/handler"
	"github.com/suraboy/upload-file-worksho/app/core/repository/mail"
	"github.com/suraboy/upload-file-worksho/app/core/repository/storage"
	"github.com/suraboy/upload-file-worksho/app/core/service"
	"google.golang.org/api/option"
	"io/ioutil"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
)

func main() {
	// setup project
	ctx := context.Background()
	conf := new(config.AppConfig)
	if err := MapConfig(conf); err != nil {

		panic("load app config error")

	}
	router := fiber.New()

	// Initialize Google Cloud Storage client
	credsJSON := []byte(conf.Storage.Services[0].Config.JSONCredential)
	client, err := googleStorage.NewClient(ctx, option.WithCredentialsJSON(credsJSON))
	if err != nil {
		log.Fatalf("Failed to initialize GCS storage: %v", err)
	}

	storagerepo := storage.NewGCSStorage(storage.Config{
		AppConfig: conf,
		Storage:   client,
	})

	log.Info(conf.Email.ApiKey)
	sg := sendgrid.NewSendClient(conf.Email.ApiKey)

	emailrepo := mail.NewSendGridClient(mail.Config{
		AppConfig:  conf,
		MailClient: sg,
	})

	svc := service.NewService(service.Config{
		AppConfig: conf,
		Storage:   storagerepo,
		Email:     emailrepo,
	})
	hdl := handler.SetupProcess(svc)

	router.Post("/upload", hdl.UploadFilesHandler)
	router.Get("/:file_url", hdl.GetFilesHandler)

	go func() {
		address := fmt.Sprintf(":%d", conf.Server.Port)
		if err := router.Listen(address); err != nil {
			log.Error(fmt.Sprintf("server error: %v", err))
			os.Exit(1)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		syscall.SIGTERM, // kill -SIGTERM XXXX
	)

	select {
	case <-signalChan:
		log.Info("terminating: received signal")
	}
	log.Info("shutting down...")
}

const (
	t_CONFIG_DIRECTORY = "config"
)

func MapConfig(conf interface{}, ignoreMergeConfig ...string) error {

	viper.SetDefault("config.path", "./config")
	err := viper.BindEnv("config.path", "CONFIG_PATH")
	if err != nil {
		log.Info("warning: %s \n", err)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(viper.GetString("config.path"))

	err = viper.ReadInConfig()
	if err != nil {
		log.Info("warning: %s \n", err)
		return err
	}

	ignoreConfig := map[string]string{}
	for _, name := range ignoreMergeConfig {
		ignoreConfig[name] = name
	}

	files, err := ioutil.ReadDir(t_CONFIG_DIRECTORY)
	if err != nil {
		log.Info("warning read Config Directory: %s \n", err)
		return err
	}

	newMap := map[string]any{}
	for _, file := range files {
		filename := file.Name()
		configServiceName, _ := strings.CutSuffix(filename, ".yaml")

		if _, ok := ignoreConfig[configServiceName]; !ok {
			v := viper.New()
			v.SetConfigName(configServiceName)
			v.AddConfigPath(viper.GetString("config.path"))
			err = v.ReadInConfig()
			if err != nil {
				log.Info("warning: %s \n", err)
			}
			newMap = mergeMapConfig(newMap, v.AllSettings())
		}
	}

	if len(newMap) > 0 {
		if err = viper.MergeConfigMap(newMap); err != nil {
			log.Info("merge warning: %s \n", err)
			return err
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err = viper.Unmarshal(conf); err != nil {
		return err
	}
	return nil
}

func mergeMapConfig(source map[string]any, input map[string]any) (output map[string]any) {
	output = source

	for k, v := range input {
		pValue := reflect.ValueOf(v)
		if pValue.Kind() == reflect.Ptr {
			pValue = pValue.Elem()
		}

		switch pValue.Kind() {
		case reflect.Map:
			if sourceValue, ok := source[k]; ok {
				newslice := []any{}
				newslice = append(newslice, v)
				newslice = append(newslice, sourceValue)
				output[k] = newslice
			} else {
				output[k] = v
			}
		default:
			output[k] = v
		}
	}

	return
}
