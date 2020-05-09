package access_token

import (
	"strings"

	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/domain/access_token"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/repository/db"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/repository/rest"
	"github.com/dzikrisyafi/kursusvirtual_utils-go/rest_errors"
)

type Service interface {
	GetById(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RestErr
}

type service struct {
	restUsersRepo rest.RestUsersRepository
	dbRepository  db.DbRepository
}

func NewService(usersRepo rest.RestUsersRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUsersRepo: usersRepo,
		dbRepository:  dbRepo,
	}
}

func (s *service) GetById(accessTokenID string) (*access_token.AccessToken, rest_errors.RestErr) {
	accessTokenID = strings.TrimSpace(accessTokenID)
	if len(accessTokenID) == 0 {
		return nil, rest_errors.NewBadRequestError("invalid access token id")
	}

	accessToken, err := s.dbRepository.GetById(accessTokenID)
	if err != nil {
		return nil, err
	}

	if accessToken.IsExpired() {
		return nil, rest_errors.NewUnauthorizedError("access token is expired")
	}

	return accessToken, nil
}

func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	// TODO: Support both grant type: client_credentials and password

	// Authenticate the user against the Users API
	user, err := s.restUsersRepo.LoginUser(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	// Generate a new access token
	at := access_token.GetNewAccessToken(user.ID)
	at.Generate()

	// Save the new access token in mysql
	if err := s.dbRepository.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken) rest_errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}

	return s.dbRepository.UpdateExpirationTime(at)
}
