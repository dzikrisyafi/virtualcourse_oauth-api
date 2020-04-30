package db

import (
	"strings"

	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/clients/mysql"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/domain/access_token"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/utils/errors"
)

const (
	queryGetAccessToken    = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	queryCreateAccessToken = "INSERT INTO access_tokens(access_token, user_id, client_id, expires) VALUES(?, ?, ?, ?);"
	queryUpdateExpires     = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessToken) *errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetById(ID string) (*access_token.AccessToken, *errors.RestErr) {
	stmt, err := mysql.DbConn().Prepare(queryGetAccessToken)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	var result access_token.AccessToken
	row := stmt.QueryRow(ID)
	if getErr := row.Scan(&result.AccessToken, &result.UserID, &result.ClientID, &result.Expires); getErr != nil {
		if strings.Contains(getErr.Error(), "no rows") {
			return nil, errors.NewNotFoundError("no access token found with given id")
		}
		return nil, errors.NewInternalServerError(getErr.Error())
	}

	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) *errors.RestErr {
	stmt, err := mysql.DbConn().Prepare(queryCreateAccessToken)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	if _, err = stmt.Exec(at.AccessToken, at.UserID, at.ClientID, at.Expires); err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) *errors.RestErr {
	stmt, err := mysql.DbConn().Prepare(queryUpdateExpires)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	if _, err = stmt.Exec(at.Expires, at.AccessToken); err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}