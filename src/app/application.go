package app

import (
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/http"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/repository/db"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/repository/rest"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/services/access_token"
	"github.com/dzikrisyafi/kursusvirtual_utils-go/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	atHandler := http.NewAccessTokenHandler(access_token.NewService(rest.NewRestUsersRepository(), db.NewRepository()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)

	logger.Info("start the application...")
	router.Run(":8001")
}
