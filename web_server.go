package krap

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const envProd string = "prod"

type WebConfig struct {
	Base string
	Port uint
}

func (c WebConfig) FindError() error {
	if c.Port == 0 {
		return errors.New("invalid API port")
	}
	if c.Base == "" {
		return errors.New("invalid API base")
	}
	return nil
}

func WebServer(cfg *WebConfig, appEnv string) (*gin.Engine, string) {
	if appEnv == envProd {
		gin.SetMode(gin.ReleaseMode)
	}

	corsCfg := cors.DefaultConfig()
	corsCfg.MaxAge = 12 * time.Hour
	corsCfg.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Authorization",
		"Accept",
		"User-Agent",
		"Cache-Control",
	}
	corsCfg.ExposeHeaders = []string{
		"Content-Length",
	}
	corsCfg.AllowMethods = []string{
		"GET",
		"POST",
		"PATCH",
		"DELETE",
	}
	corsCfg.AllowAllOrigins = true

	server := gin.Default()
	server.Use(cors.New(corsCfg))
	address := fmt.Sprintf(":%d", cfg.Port)
	return server, address
}
