package app

import (
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/clients/mysql"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/domain/access_token"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/http"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/repository/db"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	err := mysql.DbConn().Ping()
	if err != nil {
		panic(err)
	}
	defer mysql.DbConn().Close()

	atHandler := http.NewHandler(access_token.NewService(db.NewRepository()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)
	router.Run(":8001")
}
