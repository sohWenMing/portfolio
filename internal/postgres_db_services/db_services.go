package postgresdbservices

import (
	"context"
	"database/sql"

	sqlqueries "github.com/sohWenMing/portfolio/internal/db_infra/postgres_sqlc_generated"
	dbinterface "github.com/sohWenMing/portfolio/internal/db_interface"
)

type DBServices struct {
	UserService *UserService
}

type UserService struct {
	queries *sqlqueries.Queries
}

func (u *UserService) CreateUser(ctx context.Context, arg dbinterface.CreateUserInterfaceParams) (int64, error) {
	paramArg := sqlqueries.CreateUserParams(arg)
	returnedId, err := u.queries.CreateUser(ctx, paramArg)
	if err != nil {
		return 0, err
	}
	return returnedId, nil
}
func (u *UserService) GetUserDetailsById(ctx context.Context, id int64) (dbinterface.UserDetails, error) {
	userDetails, err := u.queries.GetUserDetailsById(ctx, id)
	if err != nil {
		return dbinterface.UserDetails{}, nil
	}
	return dbinterface.UserDetails(userDetails), nil
}
func (u *UserService) DeleteUserById(ctx context.Context, id int64) error {
	err := u.queries.DeleteUserById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func InitDBServices(db *sql.DB) *DBServices {
	initQueries := sqlqueries.New(db)
	return &DBServices{
		UserService: &UserService{queries: initQueries},
	}
}

/*
Takes in initial db connection, returns pointer to DBServices struct that has CRUD actions categorized by service type.
Eg DBServices.UserService - has CRUD actions that are related to Users
*/
