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
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Instalando extension uuid-ossp")
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	cfg.BeforeSetup(db)

	log.Debug("Migrando tabla agentes")
	db.AutoMigrate(&Agents{})

	log.Debug("Migrando tabla de roles")
	db.AutoMigrate(&Roles{})

	log.Debug("Migrando tabla enviroment")
	db.AutoMigrate(&ConfigVars{})
	createEnvVars(cfg.Variables, db)

	log.Debug("Migrando tabla de permisos")
	db.AutoMigrate(&Permissions{})
	createPermissions(cfg.Routes, db)

	// Create New Tables of config
	for _, table := range cfg.Tables {
		db.AutoMigrate(table)
	}

	cfg.AfterSetup(db)

	return db
}

func createEnvVars(env []ConfigVars, db *gorm.DB) {
	for _, _var := range env {
		if db.First(&ConfigVars{}, "key = ?", _var.Key).RowsAffected == 0 {
			log.Debug("Creando variable: " + _var.Key)
			db.Create(&ConfigVars{
				Key:         _var.Key,
				Value:       _var.Value,
				Description: _var.Description,
			})
		}
	}
}

func createPermissions(routes []Routes, db *gorm.DB) {
	for _, route := range routes {
		var split_route = strings.Split(route.Path, "/")

		if split_route[1] == "webhook" {
			return
		}

		if db.First(&Permissions{}, "path = ? AND method = ?", route.Path, route.Method).RowsAffected == 0 {
			log.Debug("Creando permiso: " + route.Name)
			pid := uuid.New()
			db.Create(&Permissions{
				Pid:         pid,
				Name:        route.Name,
				Description: route.Description,
				Category:    route.Category,
				Self:        route.Self,
				Path:        route.Path,
				Method:      route.Method,
			})
		}
	}
}
