package rest

import (
	"net/http"
	"os"
	"testing"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8000/users/login",
		ReqBody:      `{"email": "email@gmail.com", "password": "the-password"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid rest client response when trying to login user", err.Message)
}
func TestLoginUserInvalidUserInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8000/users/login",
		ReqBody:      `{"email": "email@gmail.com", "password": "the-password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "invalid login credentials", "status": "404", "error": "not_found"}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interfaces when trying to login user", err.Message)
}
func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8000/users/login",
		ReqBody:      `{"email": "email@gmail.com", "password": "the-password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "invalid login credentials", "status": 404, "error": "not_found"}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid error interfaces when trying to login user", err.Message)
}
func TestLoginUserInvalidJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8000/users/login",
		ReqBody:      `{"email": "email@gmail.com", "password": "the-password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": "13", "username": "dzikri", "firstname": "Dzikri", "surname": "Syafi Auliya", "email": "dzikrisyafi@lepkom.com"}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when trying to unmarshal users login response", err.Message)
}
func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8000/users/login",
		ReqBody:      `{"email": "email@gmail.com", "password": "the-password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 13, "username": "dzikri", "firstname": "Dzikri", "surname": "Syafi Auliya", "email": "dzikrisyafi@lepkom.com"}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 13, user.ID)
	assert.EqualValues(t, "dzikri", user.Username)
	assert.EqualValues(t, "Dzikri", user.Firstname)
	assert.EqualValues(t, "Syafi Auliya", user.Surname)
	assert.EqualValues(t, "dzikrisyafi@lepkom.com", user.Email)
}
