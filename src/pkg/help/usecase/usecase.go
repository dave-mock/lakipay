package usecase

import (
	"auth/src/pkg/help/core/entity"
	"log"
)

type Error struct {
	Type    string
	Message string
}

func (err Error) Error() string {
	return err.Message
}

type Usecase struct {
	log  *log.Logger
	repo Repo
}

func New(log *log.Logger, repo Repo) Interactor {
	return Usecase{log: log, repo: repo}
}

func (uc Usecase) CreatePhonePrefix(prefix, pattern string) (*entity.PhonePrefix, error) {
	// Errors
	var ErrFailedToCreatePhonePrefix string = "FAILED_TO_CREATE_PHONE_PREFIX"

	// [TODO] Validate

	uc.log.Println("CREATE-LOG -0-")
	// Create
	phonePrefix := &entity.PhonePrefix{
		Prefix:  prefix,
		Pattern: pattern,
	}

	uc.log.Println("CREATE-LOG -1-")
	// Store
	err := uc.repo.StorePhonePrefix(*phonePrefix)
	uc.log.Println("CREATE-LOG -2-")
	if err != nil {
		uc.log.Println("CREATE-LOG -3-")
		return nil, &Error{
			Type:    ErrFailedToCreatePhonePrefix,
			Message: err.Error(),
		}
	}

	uc.log.Println("CREATE-LOG -4-")
	// Return
	return phonePrefix, nil
}

func (uc Usecase) GetPhonePrefixByPrefix(prefix string) (*entity.PhonePrefix, error) {
	// Error
	var ErrPhonePrefixNotFound string = "PHONE_PREFIX_NOT_FOUND"

	uc.log.Println("SUB-LOG -0-")
	phonePrefix, err := uc.repo.FindPhoneprefixByPrefix(prefix)
	uc.log.Println("SUB-LOG -1-")
	if err != nil {
		uc.log.Println("SUB-LOG -2-")
		return nil, &Error{
			Type:    ErrPhonePrefixNotFound,
			Message: err.Error(),
		}
	}

	return phonePrefix, nil
}
