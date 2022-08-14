package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/privatesquare/bkst-go-utils/utils/config"
	"github.com/privatesquare/bkst-go-utils/utils/errors"
	"github.com/privatesquare/bkst-go-utils/utils/httputils"
	"github.com/privatesquare/bkst-go-utils/utils/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	url       = "http://localhost:8000"
	apiV1Path = "/api/v1"
)

type Config struct {
	AnkiConnectURL string `mapstructure:"ANKI_CONNECT_URL" required:"true"`
	AnkiDeckName   string `mapstructure:"ANKI_DECK_NAME" required:"true"`
	AnkiDeckModel  string `mapstructure:"ANKI_DECK_MODEL" required:"true"`
}

func SetupRoutes() *gin.Engine {
	cnf := new(Config)
	if restErr := loadConfig(cnf); restErr != nil {
		logger.Error(restErr.Message, errors.New(restErr.Error))
		os.Exit(1)
	}
	router := httputils.NewRouter()
	server := &server{
		cnf: cnf,
	}
	apiBaseUrl := url + apiV1Path
	router = RegisterHandlersWithOptions(router, server, GinServerOptions{BaseURL: apiV1Path})
	router.StaticFS(apiV1Path+"/swagger", gin.Dir("api", false))
	router.GET(apiV1Path+httputils.ApiHealthPath, httputils.Health)
	router.GET(apiV1Path+"/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(apiBaseUrl+"/swagger/eng-vocab-builder_v1.yaml")))
	logger.Infof(httputils.ServerStartupSuccessMsg, apiBaseUrl)
	return router
}

func loadConfig(cnf *Config) *errors.RestErr {
	config.Load(cnf)
	err := config.Validate(cnf)
	if err != nil {
		return errors.InternalServerError(err.Error())
	}
	return nil
}
