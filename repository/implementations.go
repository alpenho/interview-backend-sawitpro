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
