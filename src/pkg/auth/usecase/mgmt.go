package usecase

import (
	"auth/src/pkg/auth/core/entity"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
)

// / User
func (uc Usecase) CreateUser(token string, sirName string, firstName string, lastName string) (*entity.User, error) {
	// Error
	var (
		ErrCreatingUser string = "FAILED_TO_CREATE_USER"
		ErrSignIn       string = "SIGN_IN"
	)

	// Validate token
	err := uc.CheckPreSession(token)
	if err != nil {
		return nil, err
	}

	// [TODO] - Validate input

	// Find user
	phoneAuth, err := uc.repo.FindPhoneAuth(token)
	if err != nil {
		return nil, Error{
			Type:    ErrCreatingUser,
			Message: err.Error(),
		}
	}

	// Check if user exists
	usr, err := uc.repo.FindUserUsingPhoneIdentity(phoneAuth.Phone.Id)
	if err != nil {
		return nil, Error{
			Type:    ErrCreatingUser,
			Message: err.Error(),
		}
	}

	if usr != nil {
		return nil, Error{
			Type:    ErrSignIn,
			Message: "user is already registered, please try to sign in",
		}
	}

	// Create user
	var user entity.User
	id := uuid.New()

	user = entity.User{
		Id:        id,
		SirName:   sirName,
		FirstName: firstName,
		LastName:  lastName,
	}

	// Store User
	err = uc.repo.StoreUser(user)
	if err != nil {
		return nil, Error{
			Type:    ErrCreatingUser,
			Message: err.Error(),
		}
	}

	// Create Phone Identity
	var phoneIdentity entity.PhoneIdentity
	id = uuid.New()
	phoneIdentity = entity.PhoneIdentity{
		Id:    id,
		Phone: phoneAuth.Phone,
		User:  user,
	}

	// Store phone identity
	err = uc.repo.StorePhoneIdentity(phoneIdentity)
	if err != nil {
		return nil, Error{
			Type:    ErrCreatingUser,
			Message: err.Error(),
		}
	}

	return &user, nil
}

func (uc Usecase) GetUserByPhoneNumber(phoneId uuid.UUID) (entity.User, error) {

	var ErrUserNotFound string = "USER_NOT_FOUND"

	user, err := uc.repo.FindUserUsingPhoneIdentity(phoneId)
	if err != nil {
		return entity.User{}, Error{
			Type:    ErrUserNotFound,
			Message: err.Error(),
		}
	}
	return *user, nil
}

func (uc Usecase) CreatePasswordIdentity(userId uuid.UUID, password string, hint string) (*entity.PasswordIdentity, error) {
	var identity entity.PasswordIdentity

	hasher := sha256.New()
	_, err := hasher.Write([]byte(password))
	if err != nil {
		return nil, Error{
			Type:    "ERRCRATINGPASS",
			Message: err.Error(),
		}
	}

	identity = entity.PasswordIdentity{
		Id:        uuid.New(),
		User:      entity.User{Id: userId},
		Password:  base64.URLEncoding.EncodeToString(hasher.Sum(nil)),
		Hint:      hint,
		CreatedAt: time.Now(),
	}

	// Store password
	err = uc.repo.StorePasswordIdentity(identity)
	if err != nil {
		return nil, Error{
			Type:    "FAILED_TO_CREATE_2FA",
			Message: err.Error(),
		}
	}

	return &identity, nil
}
