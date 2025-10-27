package businesslogic

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	dbinterface "github.com/sohWenMing/portfolio/internal/db_interface"
)

type passwordInterface interface {
	HashPassword(password string) (hash string, err error)
	CheckPassword(attemptedPassword string, hashedPassword string) (err error)
}

// Struct holds all the interfaces that give access to CRUD operations regarding users
type UserService struct {
	defaultTimeout time.Duration
	dbService      dbinterface.UserService
	pwInterface    passwordInterface
}

// Initializes the business logic level User Service. DefaultTimeout that is passed in is used to define the
// timeout for all contexts that are passed into database related operations. dbService should be an interface
// that ties the UserService to the required database crud operations
func InitUserService(
	defaultTimeout time.Duration,
	dbService dbinterface.UserService,
	pwInterface passwordInterface) *UserService {
	return &UserService{
		defaultTimeout: defaultTimeout,
		dbService:      dbService,
		pwInterface:    pwInterface,
	}
}

func (us *UserService) CheckPassword(attemptedPassword string, hashedPassword string) (err error) {
	err = us.pwInterface.CheckPassword(attemptedPassword, hashedPassword)
	return err
}

func (us *UserService) GetUserDetailsById(id int64) (dbinterface.UserDetails, error) {
	ctx, cancel := us.CreateTimeoutContext()
	defer cancel()
	details, err := us.dbService.GetUserDetailsById(ctx, id)
	if err != nil {
		return dbinterface.UserDetails{}, err
	}
	return details, nil
}

func (us *UserService) CreateUser(email, password string) (id int64, err error) {
	err = us.validateEmailAndPasswordNotNull(email, password)
	if err != nil {
		return 0, err
	}
	err = us.checkEmailAlreadyInUse(email)
	if err != nil {
		return 0, err
	}
	hash, err := us.pwInterface.HashPassword(password)
	if err != nil {
		return 0, err
	}
	ctx, cancel := us.CreateTimeoutContext()
	defer cancel()
	userId, err := us.dbService.CreateUser(ctx, dbinterface.CreateUserInterfaceParams{
		Email:          email,
		HashedPassword: hash,
	})
	if err != nil {
		return 0, err
	}
	return userId, err
}
func (us *UserService) DeleteUserById(id int64) (int64, error) {
	ctx, cancel := us.CreateTimeoutContext()
	defer cancel()
	numrows, err := us.dbService.DeleteUserById(ctx, id)
	return numrows, err
}

func (us *UserService) checkEmailAlreadyInUse(email string) error {
	ctx, cancel := us.CreateTimeoutContext()
	defer cancel()
	currEmailCountByEmail, err := us.dbService.GetEmailCountByEmail(ctx, strings.TrimSpace(email))
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

func (us *UserService) validateEmailAndPasswordNotNull(email, password string) error {
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
