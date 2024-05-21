package usecase

import (
	"auth/src/pkg/auth/core/entity"
	"auth/src/pkg/jwt"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Initiate Authentication
// Responsible for creating a unique sharable pre session token which will be used in auth processes
func (uc Usecase) InitPreSession() (entity.PreSession, error) {

	uc.log.SetPrefix("[AUTH] [USECASE] [InitPreSession] ")

	// Errors
	var ErrFailedToInitiateAuth string = "FAILED_TO_CREATE_PRE_SESSION"

	id := uuid.New()

	token := jwt.Encode(jwt.Payload{
		Exp:    time.Now().Unix() + 1800,
		Public: id,
	}, "pre_session_secret")

	preSession := entity.PreSession{
		Id:        id,
		Token:     token,
		CreatedAt: time.Now(),
	}

	uc.log.Println("created pre session")

	// Store pre session record
	err := uc.repo.StorePreSession(preSession)
	uc.log.Println("storing pre session")
	if err != nil {
		uc.log.Printf("failed to store pre session : %s\n", err)
		return preSession, Error{
			Type:    ErrFailedToInitiateAuth,
			Message: err.Error(),
		}
	}

	uc.log.Println("initiated authentication")
	// Return presession
	return preSession, nil
}

func (uc Usecase) CheckPreSession(token string) error {
	uc.log.SetPrefix("[AUTH] [USECASE] [CheckPreSession] ")
	var ErrInvalidPreSessionToken string = "INVALID_PRE_SESSION_TOKEN"
	uc.log.Println("checking presession")
	_, err := jwt.Decode(token, "pre_session_secret")
	if err != nil {
		uc.log.Printf("failed checking presession : %s\n", err.Error())
		return Error{
			Type:    ErrInvalidPreSessionToken,
			Message: err.Error(),
		}
	}

	uc.log.Println("checked presession")
	return nil
}

// Authenticate Device

// Responsible for authenticating a device
func (uc Usecase) AuthDevice(token string, ip net.IPAddr, name string, agent string) error {
	uc.log.SetPrefix("[AUTH] [USECASE] [AuthDevice] ")
	// Error
	var ErrFailedToAuthenticateDevice string = "FAILED_TO_AUTH_DEVICE"

	// Check token
	err := uc.CheckPreSession(token)
	if err != nil {
		return err
	}

	// Do device validation
	// Create device
	var device entity.Device

	id := uuid.New()
	device = entity.Device{
		Id:        id,
		IP:        ip,
		Name:      name,
		Agent:     agent,
		CreatedAt: time.Now(),
	}

	// Store device info
	err = uc.repo.StoreDevice(device)
	if err != nil {
		return Error{
			Type:    ErrFailedToAuthenticateDevice,
			Message: err.Error(),
		}
	}

	// Create device auth
	var deviceAuth entity.DeviceAuth

	id = uuid.New()
	deviceAuth = entity.DeviceAuth{
		Id:        id,
		Device:    device,
		Token:     token,
		CreatedAt: time.Now(),
	}

	// Store device auth
	err = uc.repo.StoreDeviceAuth(deviceAuth)
	if err != nil {
		return Error{
			Type:    ErrFailedToAuthenticateDevice,
			Message: err.Error(),
		}
	}

	// Update device auth
	err = uc.repo.UpdateDeviceAuthStatus(deviceAuth.Id, true)
	if err != nil {
		return Error{
			Type:    ErrFailedToAuthenticateDevice,
			Message: err.Error(),
		}
	}

	return nil
}

// Check device auth
func (uc Usecase) CheckDeviceAuth(token string) error {
	var ErrFailedToAuthenticateDevice string = "FAILED_TO_AUTH_DEVICE"

	// Check token
	err := uc.CheckPreSession(token)
	if err != nil {
		return Error{
			Type:    ErrFailedToAuthenticateDevice,
			Message: err.Error(),
		}
	}

	// Get device auth from db
	deviceAuth, err := uc.repo.FindDeviceAuth(token)
	if err != nil {
		return Error{
			Type:    ErrFailedToAuthenticateDevice,
			Message: err.Error(),
		}
	}
	// Check the status
	if !deviceAuth.Status {
		return Error{
			Type:    ErrFailedToAuthenticateDevice,
			Message: "device is not authenticated",
		}
	}
	return nil
}

// Authenticate Phone
func (uc Usecase) InitPhoneAuth(token, prefix, number string) (*entity.PhoneAuth, error) {
	// Error
	var ErrFailedToInitPhoneAuth string = "FAILED_TO_INITIATE_PHONE_AUTH"

	// Validate token
	err := uc.CheckPreSession(token)
	if err != nil {
		return nil, err
	}

	// [TODO] Validate Phone

	// Find / create phone
	phone, err := uc.repo.FindPhone(prefix, number)
	if err != nil {
		// Create phone
		return nil, Error{
			Type:    ErrFailedToInitPhoneAuth,
			Message: err.Error(),
		}
	}

	if phone == nil {
		id := uuid.New()

		phone = &entity.Phone{
			Id:        id,
			Prefix:    prefix,
			Number:    number,
			CreatedAt: time.Now(),
		}

		err := uc.repo.StorePhone(*phone)
		if err != nil {
			return nil, Error{
				Type:    ErrFailedToInitPhoneAuth,
				Message: err.Error(),
			}
		}
	}

	// Create phone auth
	var phoneAuth entity.PhoneAuth
	id := uuid.New()
	method := "SMS"
	codeLength := 5
	timeout := 120
	// Generate Code
	// otp := rand.Intn(99999-10000) + 10000
	otp := 12345

	phoneAuth = entity.PhoneAuth{
		Id:      id,
		Token:   token,
		Phone:   *phone,
		Method:  method,
		Length:  int64(codeLength),
		Timeout: int64(timeout),
		Code: jwt.Encode(jwt.Payload{
			Exp: time.Now().Unix() + 30*60,
		}, fmt.Sprint(otp)),
	}

	err = uc.repo.StorePhoneAuth(phoneAuth)
	if err != nil {
		return nil, Error{
			Type:    ErrFailedToInitPhoneAuth,
			Message: err.Error(),
		}
	}

	// Send code
	go func() {
		uc.sms.SendSMS(phone.String(), fmt.Sprintf("Your LakiPay verification code is: %s. Do not share this code with anyone. Our employees will never ask for the code.", strconv.Itoa(otp)))
	}()

	return &phoneAuth, nil
}

func (uc Usecase) AuthPhone(token, prefix, number, otp string) error {
	// Error
	var ErrFailedToAuthPhone string = "FAILED_TO_AUTH_PHONE"

	// Validate token
	err := uc.CheckPreSession(token)
	if err != nil {
		log.Println("0")
		return err
	}

	// Get Phone Auth
	phoneAuth, err := uc.repo.FindPhoneAuth(token)
	if err != nil {
		log.Println("1")
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: err.Error(),
		}
	}

	// Validate phone
	if phoneAuth.Phone.Prefix != prefix || phoneAuth.Phone.Number != number {
		log.Println("2")
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: "Phone is not valid",
		}
	}

	// Check the status
	if phoneAuth.Status {
		return nil
	}

	// Validate otp
	_, err = jwt.Decode(phoneAuth.Code, otp)
	if err != nil {
		if err.Error() == "invalid token" {
			return Error{
				Type:    ErrFailedToAuthPhone,
				Message: "Incorrect OTP",
			}
		}
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: err.Error(),
		}
	}

	// Update Status
	err = uc.repo.UpdatePhoneAuthStatus(phoneAuth.Id, true)
	if err != nil {
		log.Println("5")
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: err.Error(),
		}
	}

	return nil
}

func (uc Usecase) CheckPhoneAuth(token string) error {
	// Error
	var ErrFailedToAuthPhone = "FAILED_TO_AUTH_PHONE"

	// Validate token
	err := uc.CheckPreSession(token)
	if err != nil {
		return err
	}

	// Get Phone Auth
	phoneAuth, err := uc.repo.FindPhoneAuth(token)
	if err != nil {
		log.Println(err)
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: err.Error(),
		}
	}

	// Check the status
	if !phoneAuth.Status {
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: "Phone unauthenticated",
		}
	}

	return nil
}

// Password

func (uc Usecase) InitPasswordAuth(token string, password string, hint string) (*entity.PasswordAuth, error) {
	// Validate token
	err := uc.CheckPreSession(token)
	if err != nil {
		return nil, Error{
			Type:    "",
			Message: err.Error(),
		}
	}

	// Find user
	phoneAuth, err := uc.repo.FindPhoneAuth(token)
	if err != nil {
		return nil, Error{
			Type:    "ERRAUTHPASS",
			Message: err.Error(),
		}
	}

	user, err := uc.repo.FindUserUsingPhoneIdentity(phoneAuth.Phone.Id)
	if err != nil {
		return nil, Error{
			Type:    "ERRAUTHPASS",
			Message: err.Error(),
		}
	}

	uc.log.Println("user.Id")
	uc.log.Println(user.Id)

	pass, err := uc.CreatePasswordIdentity(user.Id, password, hint)
	if err != nil {
		return nil, Error{
			Type:    "ERRAUTHPASS",
			Message: err.Error(),
		}
	}

	passAuth := entity.PasswordAuth{
		Id:        uuid.New(),
		Token:     token,
		Password:  *pass,
		Status:    true,
		CreatedAt: time.Now(),
	}

	// Store pass auth
	err = uc.repo.StorePasswordAuth(passAuth)

	return &passAuth, err
}

func (uc Usecase) AuthPassword(token string, password string) error {
	// Validate token
	err := uc.CheckPreSession(token)
	if err != nil {
		return Error{
			Type:    "",
			Message: err.Error(),
		}
	}

	// Find user
	phoneAuth, err := uc.repo.FindPhoneAuth(token)
	if err != nil {
		uc.log.Println(err)
		return Error{
			Type:    "ERRAUTHPASS",
			Message: err.Error(),
		}
	}

	user, err := uc.repo.FindUserUsingPhoneIdentity(phoneAuth.Phone.Id)
	if err != nil {
		uc.log.Println(err)
		return Error{
			Type:    "ERRAUTHPASS",
			Message: err.Error(),
		}
	}

	// Get user password
	pass, err := uc.repo.FindPasswordIdentityByUser(user.Id)
	if err != nil {
		uc.log.Println(err)
		return Error{
			Type:    "ERRAUTHPASS",
			Message: err.Error(),
		}
	}

	if pass == nil {
		uc.log.Println("err")
		return Error{
			Type:    "ERRAUTHPASS",
			Message: "No password found",
		}
	}

	hasher := sha256.New()
	_, err = hasher.Write([]byte(password))
	uc.log.Println(err)
	if err != nil {
		return Error{
			Type:    "ERRCRATINGPASSHASH",
			Message: err.Error(),
		}
	}

	// Compare passwords
	if base64.URLEncoding.EncodeToString(hasher.Sum(nil)) != pass.Password {
		return Error{
			Type:    "INCORRECT_PASSWORD",
			Message: "Password is incorrect",
		}
	}

	passAuth := entity.PasswordAuth{
		Id:        uuid.New(),
		Token:     token,
		Password:  *pass,
		Status:    true,
		CreatedAt: time.Now(),
	}

	// Store pass auth
	err = uc.repo.StorePasswordAuth(passAuth)

	return err

}

func (uc Usecase) CheckPasswordAuth(userId uuid.UUID, token string) error {
	// Error
	var ErrFailedToAuthPassword = "FAILED_TO_AUTH_PASSWORD"

	// Validate token
	err := uc.CheckPreSession(token)
	if err != nil {
		return err
	}

	// Get Password Auth
	passAuth, err := uc.repo.FindPasswordAuth(token)
	if err != nil {
		uc.log.Println(err)
		return Error{
			Type:    ErrFailedToAuthPassword,
			Message: err.Error(),
		}
	}

	uc.log.Println(passAuth)

	if passAuth == nil {

		uc.log.Println("pass auth nil")

		// Get user password
		pass, err := uc.repo.FindPasswordIdentityByUser(userId)
		if err != nil {
			return Error{
				Type:    "ERRAUTHPASS",
				Message: err.Error(),
			}
		}

		if pass == nil {
			return Error{
				Type:    "SET_PASSWORD",
				Message: "No password found for the requested user",
			}
		} else {
			return Error{
				Type:    "CHECK_PASSWORD",
				Message: "Password is set and must be verified before authenticating",
			}
		}
	}

	// Check the status
	if !passAuth.Status {
		return Error{
			Type:    ErrFailedToAuthPassword,
			Message: "Password unauthenticated",
		}
	}

	return nil
}

func (uc Usecase) CreateSession(token string) (*entity.Session, string, error) {

	// Error
	var (
		ErrCreatingSession string = "FAILED_TO_CREATE_SESSION"
		ErrSignUp          string = "SIGN_UP"
	)

	var session entity.Session
	var activeToken string

	// Validate token
	// Do auth steps
	// Token check
	err := uc.CheckPreSession(token)
	if err != nil {
		return &session, activeToken, err
	}

	// Device auth check
	err = uc.CheckDeviceAuth(token)
	if err != nil {
		return &session, activeToken, err
	}

	// Phone auth check
	err = uc.CheckPhoneAuth(token)
	if err != nil {
		return &session, activeToken, err
	}

	// Check user existence
	// Find user
	phoneAuth, err := uc.repo.FindPhoneAuth(token)
	if err != nil {
		return &session, token, Error{
			Type:    ErrCreatingSession,
			Message: err.Error(),
		}
	}

	user, err := uc.repo.FindUserUsingPhoneIdentity(phoneAuth.Phone.Id)
	if err != nil {
		return &session, token, Error{
			Type:    ErrSignUp,
			Message: err.Error(),
		}
	}

	if user == nil {
		return &session, token, Error{
			Type:    "SIGN_UP",
			Message: "there is no associated user with the provided phone",
		}
	}

	// Password (2FA) check
	err = uc.CheckPasswordAuth(user.Id, token)
	if err != nil {
		return &session, activeToken, err
	}

	// Find device
	deviceAuth, err := uc.repo.FindDeviceAuth(token)
	if err != nil {
		return &session, activeToken, Error{
			Type:    ErrCreatingSession,
			Message: err.Error(),
		}
	}

	device := deviceAuth.Device

	//

	id := uuid.New()

	// Generate tokens
	active := jwt.Encode(jwt.Payload{
		Exp:    time.Now().Unix() + (3 * 24 * 60 * 60),
		Public: id,
	}, "active")

	refresh := jwt.Encode(jwt.Payload{
		Exp:    time.Now().Unix() + (30 * 24 * 60 * 60),
		Public: id,
	}, "active")

	session = entity.Session{
		Id:        id,
		Device:    device,
		User:      *user,
		Token:     refresh,
		CreatedAt: time.Now(),
	}

	// Store Session
	err = uc.repo.StoreSession(session)
	if err != nil {
		return &session, activeToken, Error{
			Type:    ErrCreatingSession,
			Message: err.Error(),
		}
	}

	return &session, active, nil
}

// Check Session
func (uc Usecase) CheckSession(token string) (*entity.Session, error) {

	// Check session token
	fmt.Println("||| check session")
	pld, err := jwt.Decode(token, "active")
	if err != nil {
		return nil, Error{
			Type:    "UNAUTHORIZED",
			Message: err.Error(),
		}
	}

	fmt.Println("////////// one ", pld.Public)
	fmt.Println("////////// two ", pld)

	// Find Session by Id
	userId, err := uuid.Parse(pld.Public.(string))
	if err != nil {
		return nil, Error{
			Type:    "UNAUTHORIZED",
			Message: err.Error(),
		}
	}
	session, err := uc.repo.FindSessionById(userId)
	if err != nil {
		return nil, Error{
			Type:    "UNAUTHORIZED",
			Message: err.Error(),
		}
	}

	return session, nil
}

// Get User
func (uc Usecase) GetUserById(id uuid.UUID) (*entity.User, error) {
	var user *entity.User

	user, err := uc.repo.FindUserById(id)

	return user, err
}
