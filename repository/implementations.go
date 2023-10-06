package repository

import "context"

func (r *Repository) GetUserById(ctx context.Context, input GetUserByIdInput) (output GetUserByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number FROM users WHERE id = $1", input.Id).Scan(&output.Id, &output.FullName, &output.PhoneNumber)
	if err != nil {
		return
	}
	return
}

func (r *Repository) CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "INSERT INTO Users (full_name, phone_number, password, successful_login) VALUES ($1, $2, $3, $4) RETURNING id", input.FullName, input.PhoneNumber, input.Password, 0).Scan(&output.Id)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, input GetUserByPhoneNumberInput) (output GetUserByPhoneNumberOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, password FROM users WHERE phone_number = $1", input.PhoneNumber).Scan(&output.Id, &output.Password)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UpdateUserSuccessfulLogin(ctx context.Context, input GetUserByIdInput) (output UpdateUserSuccessfulLoginOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "UPDATE users SET successful_login = successful_login + 1 WHERE id = $1 RETURNING successful_login", input.Id).Scan(&output.SuccessfulLogin)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UpdateUserFullNameOrPhoneNumber(ctx context.Context, input UpdateUserFullNameOrPhoneNumberInput) (output UpdateUserFullNameOrPhoneNumberOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "UPDATE users SET full_name = $1, phone_number = $2 WHERE id = $3 RETURNING id", input.FullName, input.PhoneNumber, input.Id).Scan(&output.Id)
	if err != nil {
		return
	}
	return
}
