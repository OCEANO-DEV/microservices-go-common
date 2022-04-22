package middlewares

import (
	"net/http"
	"strings"

	"github.com/oceano-dev/microservices-go-common/httputil"

	"github.com/gin-gonic/gin"
)

func Authorization(claimName string, claimValue string) gin.HandlerFunc {
	return func(c *gin.Context) {
		getClaims, permissionOk := c.Get("claims")
		if permissionOk {
			var claims = getClaims.([]interface{})
			permissionOk = verifyClaimsPermission(claims, claimName, claimValue)
		}

		if !permissionOk {
			httputil.NewResponseAbort(c, http.StatusUnauthorized, "you do not have permission")
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			// 	"status": http.StatusUnauthorized,
			// 	"error":  []string{"you do not have permission"},
			// })
			return
		}

		c.Next()
	}
}

func verifyClaimsPermission(claims []interface{}, claimType string, claimValue string) bool {
	sClaimType := strings.TrimSpace(claimType)
	sClaimValue := strings.TrimSpace(claimValue)
	if len(sClaimType) == 0 || len(sClaimValue) == 0 {
		return false
	}
	for _, interator := range claims {
		values := interator.(map[string]interface{})
		if values["type"] == sClaimType && strings.Contains(values["value"].(string), sClaimValue) {
			return true
		}
	}

	return false
}
