package aumpi_core

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	schema "gorm.io/gorm/schema"
)

//SetupModels is
func SetupModels(cfg Configuration) *gorm.DB {
	dbHost := os.Getenv("PG_HOST")
	username := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbName := os.Getenv("PG_DATABASE")
	dbPort := os.Getenv("PG_PORT")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s port=%s  sslmode=require password=%s", dbHost, username, dbName, dbPort, password)
	log.Debug(dbUri)

	db, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{

		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Instalando extension uuid-ossp")
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	cfg.BeforeSetup(db)

	log.Debug("Migrando tabla agentes")
	db.AutoMigrate(&SystemAgents{})

	log.Debug("Migrando tabla de roles")
	db.AutoMigrate(&SystemRoles{})

	log.Debug("Migrando tabla enviroment")
	db.AutoMigrate(&SystemConfigVars{})
	createEnvVars(cfg.Variables, db)

	log.Debug("Migrando tabla de permisos")
	db.AutoMigrate(&SystemPermissions{})
	createPermissions(cfg.Routes, db)

	// Create New Tables of config
	for _, table := range cfg.Tables {
		db.AutoMigrate(table)
	}

	cfg.AfterSetup(db)

	return db
}

func createEnvVars(env []SystemConfigVars, db *gorm.DB) {
	for _, _var := range env {
		if db.First(&SystemConfigVars{}, "key = ?", _var.Key).RowsAffected == 0 {
			log.Debug("Creando variable: " + _var.Key)
			db.Create(&SystemConfigVars{
				Key:         _var.Key,
				Value:       _var.Value,
				Description: _var.Description,
				Type:        _var.Type,
			})
		}
	}
}

func createPermissions(routes []SystemRoutes, db *gorm.DB) {
	for _, route := range routes {
		var split_route = strings.Split(route.Path, "/")

		if split_route[1] == "webhook" {
			return
		}

		if db.First(&SystemPermissions{}, "path = ? AND method = ?", route.Path, route.Method).RowsAffected == 0 {
			log.Debug("Creando permiso: " + route.Description)
			pid := uuid.New()
			db.Create(&SystemPermissions{
				Pid:         pid,
				Description: route.Description,
				Category:    route.Category,
				Self:        route.Self,
				Path:        route.Path,
				Method:      route.Method,
			})
		}
	}
}
