package dbinterface

import (
	"context"

	postgressqlcgenerated "github.com/sohWenMing/portfolio/internal/db_infra/postgres_sqlc_generated"
)

// Interface type to be imported into database modules, has to be arg input param of the CreateUser
// method that allows object to fulfil UserCreater interface
type CreateUserInterfaceParams struct {
	Email          string
	HashedPassword string
}

// Interface type to be imported into database modules, has to be the UserDetails return value of the
// UserDetailGetter method that allows object to fulfil UserDetailGetter interface
type UserDetails struct {
	ID             int64
	Email          string
	HashedPassword string
}

// UserService interface unified all CRUD operations that are tied to a user operations
type UserService interface {
	CreateUser(ctx context.Context, arg CreateUserInterfaceParams) (int64, error)
	GetUserDetailsById(ctx context.Context, id int64) (UserDetails, error)
	DeleteUserById(ctx context.Context, id int64) (int64, error)
	GetEmailCountByEmail(ctx context.Context, email string) (int64, error)
}

type PGUserService struct {
	queries *postgressqlcgenerated.Queries
}

func (p *PGUserService) CreateUser(ctx context.Context, arg CreateUserInterfaceParams) (int64, error) {
	id, err := p.queries.CreateUser(ctx, postgressqlcgenerated.CreateUserParams{
		Email:          arg.Email,
		HashedPassword: arg.HashedPassword,
	})
	return id, err
}
func (p *PGUserService) GetUserDetailsById(ctx context.Context, id int64) (UserDetails, error) {
	details, err := p.queries.GetUserDetailsById(ctx, id)
	return UserDetails{
		ID:             details.ID,
		Email:          details.Email,
		HashedPassword: details.HashedPassword,
	}, err
}

func (p *PGUserService) DeleteUserById(ctx context.Context, id int64) (int64, error) {
	numrows, err := p.queries.DeleteUserById(ctx, id)
	return numrows, err
}

func (p *PGUserService) GetEmailCountByEmail(ctx context.Context, email string) (int64, error) {
	count, err := p.queries.GetEmailCountByEmail(ctx, email)
	return count, err
}

func InitUserServiceWithPostgres(dbtx postgressqlcgenerated.DBTX) *PGUserService {
	return &PGUserService{
		queries: postgressqlcgenerated.New(dbtx),
	}
}

// Interface to allow database wrapper function with CreateUser method to fulfil interface and be used
type UserCreater interface {
	CreateUser(ctx context.Context, arg CreateUserInterfaceParams) (int64, error)
}

// Interface to allow database wrapper function with GetUserDetaislById method to fulfil interface and be used
type UserDetailGetter interface {
	GetUserDetailsById(ctx context.Context, id int64) (UserDetails, error)
}

// Interface that allows database wrapper function with DeleteUserById method to fulfil interface and be used
type UserDeleter interface {
	DeleteUserById(ctx context.Context, id int64) error
}
