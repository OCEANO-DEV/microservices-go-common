package httputil

import (
	"github.com/gin-gonic/gin"
	"github.com/oceano-dev/microservices-go-common/helpers"
)

type ResponseError struct {
	Status int           `json:"status"`
	Error  []interface{} `json:"error"`
}

func NewResponseError(c *gin.Context, statusCode int, err interface{}) {
	response := &ResponseError{
		Status: statusCode,
		Error: []interface{}{
			err,
		},
	}

	c.JSON(statusCode, response)
}

func NewResponseAbort(c *gin.Context, statusCode int, err interface{}) {
	response := &ResponseError{
		Status: statusCode,
		Error: []interface{}{
			err,
		},
	}

	c.AbortWithStatusJSON(statusCode, response)
}

type ResponseSuccess struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewResponseSuccess(c *gin.Context, statusCode int, message string) {
	response := &ResponseSuccess{
		Status:  statusCode,
		Message: message,
	}

	c.JSON(statusCode, response)
}

type ResponseCredentials struct {
	AccessToken  string           `json:"accessToken"`
	RefreshToken string           `json:"refreshToken"`
	User         *UserCredentials `json:"user"`
}

type UserCredentials struct {
	Id     helpers.ID `json:"id"`
	Email  string     `json:"email"`
	Claims []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"claims"`
	Version uint `json:"version"`
}

type User struct {
	ID     helpers.ID `json:"id"`
	Email  string     `json:"email,omitempty"`
	Claims []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"claims,omitempty"`
	Version uint `json:"version" validate:"required"`
}

func NewResponseCredentials(c *gin.Context, statusCode int, user *User, accessToken string, refreshToken string) {
	response := &ResponseCredentials{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &UserCredentials{
			Id:    user.ID,
			Email: user.Email,
			Claims: []struct {
				Type  string "json:\"type\""
				Value string "json:\"value\""
			}(user.Claims),
			Version: user.Version,
		},
	}

	c.JSON(statusCode, response)
}
