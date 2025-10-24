package dbinterface

import "context"

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
