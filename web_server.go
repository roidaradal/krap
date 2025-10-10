package krap

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type WebConfig struct {
	Base     string   // API endpoint prefix
	Port     uint     // port number
	CORSList []string // list of allowed sites for CORS
}

// Validates WebConfig
func (c WebConfig) FindError() error {
	if c.Port == 0 {
		return errors.New("invalid API port")
	}
	if c.Base == "" {
		return errors.New("invalid API base")
	}
	return nil
}

// Creates a new Gin web server
func WebServer(cfg *WebConfig, appEnv string) (*gin.Engine, string) {
	isProdEnv := appEnv == envProd
	if isProdEnv {
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
	if isProdEnv {
		corsCfg.AllowOrigins = cfg.CORSList
	} else {
		corsCfg.AllowAllOrigins = true
	}

	server := gin.Default()
	server.Use(cors.New(corsCfg))
	address := fmt.Sprintf(":%d", cfg.Port)
	return server, address
}
