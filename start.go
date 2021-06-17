package aumpi_core

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

func Start(cfg Configuration) {
	gin.SetMode(os.Getenv("APP_MODE"))
	log.SetLevel(log.TraceLevel)

	log.Trace("CREATING ROUTER")
	r := gin.Default()

	log.Trace("INIT AND INJECT DATABASE")
	db := SetupModels(cfg)
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Disable CORS
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "*",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	// Endpoint abierto para healtcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Validador de JWT
	r.Use(JWTValidator())
	// Validador de permisos
	r.Use(PermissionsValidator())

	// Agregar rutas dinamicamente
	for _, route := range cfg.Routes {
		if route.Method == "GET" {
			r.GET(route.Path, route.Function)
		} else if route.Method == "POST" {
			r.POST(route.Path, route.Function)
		} else if route.Method == "PUT" {
			r.PUT(route.Path, route.Function)
		} else if route.Method == "PATCH" {
			r.PATCH(route.Path, route.Function)
		} else if route.Method == "DELETE" {
			r.DELETE(route.Path, route.Function)
		}
	}

	log.Debug("Iniciando API")
	r.Run(os.Getenv("APP_ADDR"))
}
