package security

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"sync"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type ManagerToken struct {
	config              *config.Config
	managerSecurityKeys *ManagerSecurityKeys
}

func NewManagerToken(
	config *config.Config,
	managerSecurityKeys *ManagerSecurityKeys,
) *ManagerToken {
	return &ManagerToken{
		config:              config,
		managerSecurityKeys: managerSecurityKeys,
	}
}

func (m *ManagerToken) ReadCookieAccessToken(c *gin.Context) (*models.TokenClaims, error) {
	var err error
	tokenString, err := c.Cookie("accessToken")
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	var keyFunc = m.getKeyFunc()
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, keyFunc)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	claims, ok := token.Claims.(*models.TokenClaims)
	if !ok || !token.Valid || claims.Iss != m.config.Token.Issuer {
		return nil, errors.New("JWT failed validation")
	}

	return claims, nil
}

func (m *ManagerToken) ReadRefreshToken(c *gin.Context, tokenString string) (string, error) {
	var keyFunc = m.getKeyFunc()
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, keyFunc)
	if err != nil {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(*models.TokenClaims)
	if !ok || !token.Valid || token.Header["typ"] != "refresh" || claims.Iss != m.config.Token.Issuer {
		return "", errors.New("invalid token")
	}

	return claims.Sub, nil
}

func (m *ManagerToken) getKeyFunc() jwt.Keyfunc {
	var mux sync.Mutex
	mux.Lock()
	keys := m.managerSecurityKeys.GetAllPublicKeys()
	mux.Unlock()

	var keyFunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("expecting JWT header to have string kid")
		}

		var key *ecdsa.PublicKey
		for i := range keys {
			if keys[i].Kid == keyID {
				key = keys[i].Key
				break
			}
		}

		if key == nil {
			return nil, fmt.Errorf("unable to parse public key")
		}

		return key, nil
	}

	return keyFunc
}
