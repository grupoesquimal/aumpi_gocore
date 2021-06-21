package aumpi_core

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PermissionsValidator permite validar el acceso en base a los permisos
func PermissionsValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		// PASE COMPLETO SI ES UN WEBHOOK
		var split = strings.Split(c.FullPath(), "/")
		if len(split) > 1 && split[1] == "webhook" {
			c.Next()
			return
		}

		// OBTENER VARIABLES DE LA INYECCION
		db := c.MustGet("db").(*gorm.DB)
		uid := c.MustGet("uid").(string)

		// OBTENER DATOS DEL AGENTE
		var agent Agents
		db.First(&agent, "uid = ?", uid)

		// OBTENER ROL DEL AGENTE
		var role Roles
		db.First(&role, "rid = ?", agent.Rid)

		// OBTENER PERMISO DE LA RUTA
		var permission Permissions
		res_perm := db.First(&permission, "path = ? AND method = ?", c.FullPath(), c.Request.Method)

		if res_perm.Error != nil {
			c.AbortWithStatusJSON(404, gin.H{"success": false, "message": "Recurso no encontrado"})
			return
		}

		// VALIDAR SI ES UN PERMISO DE TIPO ROOT
		if role.Permissions == "*" {
			c.Next()
			return
		}

		// VALIDAR SI LA RUTA ES DE TIPO SELF
		if permission.Self {
			c.Next()
			return
		}

		// VALIDAR SI EL USUARIO TIENE ACCESO A LA RUTA SOLICITADA
		if strings.Contains(role.Permissions, permission.Pid.String()) {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(401, gin.H{"success": false, "message": "No cuentas con los permisos suficientes"})
	}
}
