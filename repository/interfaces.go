// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	GetUserById(ctx context.Context, input GetUserByIdInput) (output GetUserByIdOutput, err error)
	CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error)
}
