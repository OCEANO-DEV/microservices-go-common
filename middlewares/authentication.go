package middlewares

import (
	"encoding/json"
	"net/http"

	//security "github.com/oceano-dev/microservices-go-common/security/jwt"

	"github.com/oceano-dev/microservices-go-common/httputil"
	"github.com/oceano-dev/microservices-go-common/security"

	"github.com/oceano-dev/microservices-go-common/helpers"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Authentication struct {
	logger       *logrus.Logger
	managerToken *security.ManagerToken
}

func NewAuthentication(
	logger *logrus.Logger,
	managerToken *security.ManagerToken,
) *Authentication {
	return &Authentication{
		logger:       logger,
		managerToken: managerToken,
	}
}

func (auth *Authentication) Verify() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := auth.managerToken.ReadCookieAccessToken(c)
		if err != nil {
			auth.logger.Error(err.Error())
			httputil.NewResponseAbort(c, http.StatusUnauthorized, "token is not valid")
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			// 	"status": http.StatusUnauthorized,
			// 	"error":  []string{"invalid cookie"},
			// })
			return
		}

		id := claims.Sub
		if !helpers.IsValidID(id) {
			auth.logger.Error("ID is not valid")
			httputil.NewResponseAbort(c, http.StatusUnauthorized, "ID is not valid")
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			// 	"status": http.StatusUnauthorized,
			// 	"error":  []string{"ID is not valid"},
			// })
			return
		}
		c.Set("user", id)

		var claimsList []interface{}
		data, _ := json.Marshal(claims.Claims)
		json.Unmarshal(data, &claimsList)

		c.Set("claims", claimsList)

		c.Next()
	}
}
