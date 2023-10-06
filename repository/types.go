// This file contains types that are used in the repository layer.
package repository

type GetUserByIdInput struct {
	Id int32
}

type GetUserByIdOutput struct {
	Id          int32
	FullName    string
	PhoneNumber string
}

type CreateUserInput struct {
	FullName    string
	PhoneNumber string
	Password    string
}

type CreateUserOutput struct {
	Id int32
}
