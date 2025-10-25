package businesslogic

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	dbinterface "github.com/sohWenMing/portfolio/internal/db_interface"
)

// Struct holds all the interfaces that give access to CRUD operations regarding users
type UserService struct {
	defaultTimeout time.Duration
	dbService      dbinterface.UserService
}

// Initializes the business logic level User Service. DefaultTimeout that is passed in is used to define the
// timeout for all contexts that are passed into database related operations. dbService should be an interface
// that ties the UserService to the required database crud operations
func InitUserService(defaultTimeout time.Duration, dbService dbinterface.UserService) *UserService {
	return &UserService{
		defaultTimeout: defaultTimeout,
		dbService:      dbService,
	}
}

func (us *UserService) CreateUser(email, password string) (id int64, err error) {
	err = us.ValidateEmailAndPassword(email, password)
	if err != nil {
		return 0, err
	}
	err = us.checkEmailAlreadyInUse(email)
	if err != nil {
		return 0, err
	}
	//TODO: Finish up CreateUser
	return 0, err
}

func (us *UserService) checkEmailAlreadyInUse(email string) error {
	ctx, cancel := us.CreateTimeoutContext()
	currEmailCountByEmail, err := us.dbService.GetEmailCountByEmail(ctx, strings.TrimSpace(email))
	defer cancel()
	if err != nil {
		//TODO: Centralise to logging function
		fmt.Println("error occured in CheckEmailAlreadyInUser: ", err)
		return err
	}
	if currEmailCountByEmail != 0 {
		return errors.New("email is already in use, cannot be used to create user")
	}
	return nil
}

func (us *UserService) CreateTimeoutContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), us.defaultTimeout)
	return ctx, cancel
}

// Function takes in interfaces regarding user created CRUD operations, to return UserService struct

func (us *UserService) ValidateEmailAndPassword(email, password string) error {
	err := validateEmailNotNull(email)
	if err != nil {
		return err
	}
	err = validatePasswordNotNull(password)
	if err != nil {
		return err
	}
	return nil
}

func validateEmailNotNull(email string) error {
	if strings.TrimSpace(email) == "" {
		//TODO: will need to eventually unify all error handling in an error handling package
		return errors.New("email cannot be null")
	}
	return nil
}
func validatePasswordNotNull(password string) error {
	if strings.TrimSpace(password) == "" {
		//TODO: will need to eventually unify all error handling in an error handling package
		return errors.New("password cannot be null")
	}
	return nil
}
