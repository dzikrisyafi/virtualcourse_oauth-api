package db

import (
	"errors"
	"strings"

	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/clients/mysql"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/domain/access_token"
	"github.com/dzikrisyafi/kursusvirtual_utils-go/logger"
	"github.com/dzikrisyafi/kursusvirtual_utils-go/rest_errors"
)

const (
	queryGetAccessToken    = `SELECT access_token, user_id, client_id, expires, token_type FROM access_tokens WHERE access_token=?;`
	queryCreateAccessToken = `INSERT INTO access_tokens(access_token, user_id, client_id, expires, token_type) VALUES(?, ?, ?, ?, ?);`
	queryUpdateExpires     = `UPDATE access_tokens SET expires=? WHERE access_token=?;`
	queryDeleteAccessToken = `DELETE FROM access_tokens WHERE access_token=?;`
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessToken) rest_errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RestErr
	DeleteAccessToken(string) rest_errors.RestErr
}

type dbRepository struct {
}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetById(ID string) (*access_token.AccessToken, rest_errors.RestErr) {
	stmt, err := mysql.DbConn().Prepare(queryGetAccessToken)
	if err != nil {
		logger.Error("error when trying to prepare get token by id statement", err)
		return nil, rest_errors.NewInternalServerError("error when trying to get access token", errors.New("database error"))
	}
	defer stmt.Close()

	var result access_token.AccessToken
	row := stmt.QueryRow(ID)
	if getErr := row.Scan(&result.AccessToken, &result.UserID, &result.ClientID, &result.Expires, &result.TokenType); getErr != nil {
		if strings.Contains(getErr.Error(), "no rows") {
			return nil, rest_errors.NewNotFoundError("no access token found with given id")
		}
		logger.Error("error when trying to get token by id", err)
		return nil, rest_errors.NewInternalServerError("error when trying to get access token", errors.New("database error"))
	}

	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) rest_errors.RestErr {
	stmt, err := mysql.DbConn().Prepare(queryCreateAccessToken)
	if err != nil {
		logger.Error("error when trying to prepare create access token statement", err)
		return rest_errors.NewInternalServerError("error when trying to create access token", errors.New("database error"))
	}
	defer stmt.Close()

	if _, err = stmt.Exec(at.AccessToken, at.UserID, at.ClientID, at.Expires, at.TokenType); err != nil {
		logger.Error("error when trying to create access token", err)
		return rest_errors.NewInternalServerError("error when trying to create access token", errors.New("database error"))
	}

	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) rest_errors.RestErr {
	stmt, err := mysql.DbConn().Prepare(queryUpdateExpires)
	if err != nil {
		logger.Error("error when trying to prepare update access token statement", err)
		return rest_errors.NewInternalServerError("error when trying to update access token", errors.New("database error"))
	}
	defer stmt.Close()

	if _, err = stmt.Exec(at.Expires, at.AccessToken); err != nil {
		logger.Error("error when trying to update access token", err)
		return rest_errors.NewInternalServerError("error when trying to update access token", errors.New("database error"))
	}

	return nil
}

func (r *dbRepository) DeleteAccessToken(ID string) rest_errors.RestErr {
	stmt, err := mysql.DbConn().Prepare(queryDeleteAccessToken)
	if err != nil {
		logger.Error("error when trying to delete access token", err)
		return rest_errors.NewInternalServerError("error when trying to delete access token", errors.New("database error"))
	}
	defer stmt.Close()

	if _, err := stmt.Exec(ID); err != nil {
		logger.Error("error when trying to delete access token", err)
		return rest_errors.NewInternalServerError("error when trying to delete access token", errors.New("database error"))
	}

	return nil
}
