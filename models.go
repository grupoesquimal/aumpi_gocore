package aumpi_core

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//Configuration is
type Configuration struct {
	Tables      []interface{}
	Routes      []Routes
	Variables   []ConfigVars
	BeforeSetup func(db *gorm.DB)
	AfterSetup  func(db *gorm.DB)
}

//Routes is
type Routes struct {
	Name        string          // Nombre del permiso
	Description string          // Descripcion corta de lo que hace la ruta
	Category    string          // Categoria de la ruta de separado por >
	Self        bool            // Poner en true si el controlador tiene acceso unicamente a informacion del usuario que solicito la peticion por ejemplo sus permisos, sus leads, su actividad
	Path        string          // Ruta en la que respondera el controlador
	Method      string          // Metodo en el que respondera el controlador
	Function    gin.HandlerFunc // Controlador asociado a la ruta
}

//Permissions is
type Permissions struct {
	Pid         uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name        string    // Nombre del permiso
	Description string    // Descripcion corta de lo que hace la ruta
	Category    string    // Categoria de la ruta de separado por >
	Self        bool      // Poner en true si el controlador tiene acceso unicamente a informacion del usuario que solicito la peticion por ejemplo sus permisos, sus leads, su actividad
	Path        string    // Ruta en la que respondera el controlador
	Method      string    // Metodo en el que respondera el controlador
}

//Roles is
type Roles struct {
	Rid            uuid.UUID `gorm:"primaryKey;type:uuid"`
	Name           string    `gorm:"type:varchar(25)"`
	Description    string    `gorm:"type:varchar(70)"`
	Permissions    string    `gorm:"type:text"`
	PermissionsWeb string    `gorm:"type:text"`
	Editable       bool
}

//ConfigVars is
type ConfigVars struct {
	ID          uint64 `gorm:"primaryKey"`
	Description string
	Key         string
	Value       string
	Type        string
}

//Agents is
type Agents struct {
	Aid       uuid.UUID `gorm:"primaryKey;type:uuid"`
	Uid       uuid.UUID `gorm:"unique;type:uuid"`
	Rid       uuid.UUID `gorm:"type:uuid"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
