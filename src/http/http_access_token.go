package http

import (
	"net/http"

	atDomain "github.com/dzikrisyafi/kursusvirtual_oauth-api/src/domain/access_token"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/services/access_token"
	"github.com/dzikrisyafi/kursusvirtual_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
)

type AccessTokenHandler interface {
	GetById(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

type accessTokenHandler struct {
	service access_token.Service
}

func NewAccessTokenHandler(service access_token.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (handler *accessTokenHandler) GetById(c *gin.Context) {
	accessToken, err := handler.service.GetById(c.Param("access_token_id"))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, accessToken)
}

func (handler *accessTokenHandler) Create(c *gin.Context) {
	var request atDomain.AccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	accessToken, err := handler.service.Create(request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, accessToken)
}

func (handler *accessTokenHandler) Update(c *gin.Context) {
	var at atDomain.AccessToken
	at.AccessToken = c.Param("access_token_id")
	if err := handler.service.UpdateExpirationTime(at); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "success updated access token", "status": http.StatusOK})
}

func (handler *accessTokenHandler) Delete(c *gin.Context) {
	if err := handler.service.DeleteAccessToken(c.Param("access_token_id")); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "success deleted access token", "status": http.StatusOK})
}
