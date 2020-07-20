package app

import (
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/http"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/repository/db"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/repository/rest"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/services/access_token"
	"github.com/dzikrisyafi/kursusvirtual_utils-go/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	atHandler := http.NewAccessTokenHandler(access_token.NewService(rest.NewRestUsersRepository(), db.NewRepository()))
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"OPTIONS", "GET", "POST", "PUT", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Content-Length"},
	}))

	router.POST("/oauth/access_token", atHandler.Create)
	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.PUT("/oauth/access_token/:access_token_id", atHandler.Update)
	router.DELETE("/oauth/access_token/:access_token_id", atHandler.Delete)

	logger.Info("start the application...")
	router.Run(":8000")
}
