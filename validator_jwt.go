package aumpi_core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/cristalhq/jwt/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// JWTValidator permite validar el token de acceso
func JWTValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		// PASE COMPLETO SI ES UN WEBHOOK
		var split = strings.Split(c.FullPath(), "/")
		if len(split) > 1 && split[1] == "webhook" {
			c.Next()
			return
		} else {
			db := c.MustGet("db").(*gorm.DB)
			key := []byte(os.Getenv("JWT_KEY"))
			verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
			tokenString := ExtractToken(c.Request)

			if err != nil {
				fmt.Println(err.Error())
				c.AbortWithStatusJSON(500, gin.H{"success": false, "message": err.Error()})
				return
			}

			// VALIDATE SIGNATURE AND ALGORITHM
			token, errParseVerify := jwt.ParseAndVerifyString(tokenString, verifier)
			if errParseVerify != nil {
				fmt.Println(errParseVerify)
				c.AbortWithStatusJSON(401, gin.H{"success": false, "message": "No se pudo verificar"})
				return
			}

			// UNMARHAL CLAIMS
			var claims jwt.StandardClaims
			errClaims := json.Unmarshal(token.RawClaims(), &claims)
			if errClaims != nil {
				fmt.Println(errClaims)
				c.AbortWithStatusJSON(401, gin.H{"success": false, "message": "No se pudo verificar"})
				return
			}

			// VALIDATE AUDIENCE
			if !claims.IsForAudience(os.Getenv("JWT_AUD")) {
				c.AbortWithStatusJSON(401, gin.H{"success": false, "message": "No se pudo verificar"})
				return
			}

			// VALIDATE IF IS AGENT
			if db.First(&Agents{}, "uid = ?", claims.Subject).RowsAffected == 0 {
				c.AbortWithStatusJSON(401, gin.H{"success": false, "message": "No autorizado como agente"})
				return
			}

			c.Set("uid", claims.Subject)
			c.Next()
		}
	}
}
