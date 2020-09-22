package access_token

import (
	"strings"
	"time"

	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/domain/access_token"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/repository/db"
	"github.com/dzikrisyafi/kursusvirtual_oauth-api/src/repository/rest"
	"github.com/dzikrisyafi/kursusvirtual_utils-go/rest_errors"
)

// Service mocking every functions in access token service
type Service interface {
	GetByID(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RestErr
	DeleteAccessToken(string) rest_errors.RestErr
}

type service struct {
	restUsersRepo rest.RestUsersRepository
	dbRepository  db.DbRepository
}

// NewService handling repository for data transfer
func NewService(usersRepo rest.RestUsersRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUsersRepo: usersRepo,
		dbRepository:  dbRepo,
	}
}

// Get access token id from database server
func (s *service) GetByID(accessTokenID string) (*access_token.AccessToken, rest_errors.RestErr) {
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

	// Support both grant type: client_credentials and password
	var at access_token.AccessToken
	if request.GrantType == "password" {
		// Authenticate the user against the Users API
		user, err := s.restUsersRepo.LoginUser(request.Username, request.Password)
		if err != nil {
			return nil, err
		}

		// Generate a new access token by user id
		at = access_token.GetNewAccessToken(user.ID, request.ClientID)
		at.GenerateByUserID()

		if at.Validate(); err != nil {
			return nil, err
		}
	} else {
		// Generate a new access token by client id
		at = access_token.GetNewAccessToken(0, request.ClientID)
		at.GenerateByClientID()

		if at.ClientID == "" {
			restErr := rest_errors.NewBadRequestError("invalid client id")
			return nil, restErr
		}
	}

	// Save the new access token in mysql
	if err := s.dbRepository.Create(at); err != nil {
		return nil, err
	}

	return &at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken) rest_errors.RestErr {
	if at.AccessToken == "" {
		restErr := rest_errors.NewBadRequestError("invalid access token id")
		return restErr
	}

	at.Expires = int(time.Now().UTC().Add(24 * time.Hour).Unix())

	return s.dbRepository.UpdateExpirationTime(at)
}

func (s *service) DeleteAccessToken(accessTokenID string) rest_errors.RestErr {
	accessTokenID = strings.TrimSpace(accessTokenID)
	if len(accessTokenID) == 0 {
		return rest_errors.NewBadRequestError("invalid access token id")
	}

	return s.dbRepository.DeleteAccessToken(accessTokenID)
}
