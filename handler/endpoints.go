package handler

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
	"unicode"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
)

type UserData struct {
	Id          int32
	FullName    string
	PhoneNumber string
	Password    string
}

// Registration implements generated.ServerInterface.
func (s *Server) Registration(ctx echo.Context) error {
	requestBody := new(generated.Registration)
	ctx.Bind(requestBody)

	err := validateRegistrationRequestBody(requestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	hashPass, err := HashAndSalt(requestBody.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	input := repository.CreateUserInput{
		FullName:    requestBody.FullName,
		PhoneNumber: requestBody.PhoneNumber,
		Password:    hashPass,
	}
	output, err := s.Repository.CreateUser(ctx.Request().Context(), input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	user := UserData{Id: output.Id}
	tok, err := createNewJWTToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	resp := generated.RegistrationResponse{
		Id:          &output.Id,
		AccessToken: &tok,
	}
	return ctx.JSON(http.StatusOK, resp)
}

// Login implements generated.ServerInterface.
func (s *Server) Login(ctx echo.Context) error {
	requestBody := new(generated.Login)
	ctx.Bind(requestBody)

	err := validateLoginRequestBody(requestBody)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	input := repository.GetUserByPhoneNumberInput{
		PhoneNumber: requestBody.PhoneNumber,
	}
	output, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), input)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	user := UserData{Id: output.Id, Password: output.Password}
	correct := comparePasswords(user.Password, requestBody.Password)
	if !correct {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect Password")
	}

	otherInput := repository.GetUserByIdInput{
		Id: output.Id,
	}
	_, err = s.Repository.UpdateUserSuccessfulLogin(ctx.Request().Context(), otherInput)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	tok, err := createNewJWTToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	resp := generated.RegistrationResponse{
		Id:          &output.Id,
		AccessToken: &tok,
	}
	return ctx.JSON(http.StatusOK, resp)
}

// GetProfile implements generated.ServerInterface.
func (s *Server) GetProfile(ctx echo.Context) error {
	pubKey, err := os.ReadFile("rsakey.pem.pub")
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err)
	}

	reg := regexp.MustCompile(`^[B|b]earer\s`)
	tokenString := reg.ReplaceAllString(ctx.Request().Header["Authorization"][0], "")
	content, err := NewJWT(nil, pubKey).Validate(tokenString)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err)
	}

	var input repository.GetUserByIdInput
	input.Id = content.UserId
	data, err := s.Repository.GetUserById(ctx.Request().Context(), input)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	response := new(generated.ProfileDataResponse)
	response.FullName = &data.FullName
	response.PhoneNumber = &data.PhoneNumber

	return ctx.JSON(http.StatusOK, response)
}

func validateRegistrationRequestBody(requestBody *generated.Registration) error {
	errorMessageMap := map[int]string{
		1: "all of the form data must be filled in",
		2: "phone number doesn't have country code",
		3: "phone number is below 10 number or more then 13 number",
		4: "full name is below 3 character or more then 60 character",
		5: "password is below 6 character or more then 64 character",
		6: "password doesn't have number in it",
		7: "password doesn't have uppercase character in it",
		8: "password doesn't have special character in it",
	}
	validPhoneNumber := regexp.MustCompile(`^\+62`)
	phoneNumberLength := len(requestBody.PhoneNumber)
	fullNameLength := len(requestBody.FullName)
	validLength, haveNumber, haveUppercase, haveSpecialChar := verifyPassword(requestBody.Password)

	if requestBody.FullName == "" || requestBody.Password == "" || requestBody.PhoneNumber == "" {
		return fmt.Errorf(errorMessageMap[1])
	} else if !validPhoneNumber.MatchString(requestBody.PhoneNumber) {
		return fmt.Errorf(errorMessageMap[2])
	} else if phoneNumberLength < 11 || phoneNumberLength > 14 {
		return fmt.Errorf(errorMessageMap[3])
	} else if fullNameLength < 3 || fullNameLength > 60 {
		return fmt.Errorf(errorMessageMap[4])
	} else if !validLength {
		return fmt.Errorf(errorMessageMap[5])
	} else if !haveNumber {
		return fmt.Errorf(errorMessageMap[6])
	} else if !haveUppercase {
		return fmt.Errorf(errorMessageMap[7])
	} else if !haveSpecialChar {
		return fmt.Errorf(errorMessageMap[8])
	}
	return nil
}

func validateLoginRequestBody(requestBody *generated.Login) error {
	if requestBody.Password == "" || requestBody.PhoneNumber == "" {
		return fmt.Errorf("all of the form data must be filled in")
	}

	return nil
}

func verifyPassword(s string) (sixOrSixtyFour, number, upper, special bool) {
	letters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
		}
	}
	sixOrSixtyFour = letters >= 6 || letters <= 64
	return
}

func createNewJWTToken(content UserData) (string, error) {
	prvKey, err := os.ReadFile("rsakey.pem")
	if err != nil {
		return "", fmt.Errorf("error read private key: %w", err)
	}

	tok, err := NewJWT(prvKey, nil).Create(time.Hour, content)
	if err != nil {
		return "", fmt.Errorf("error creating jwt token: %w", err)
	}

	return tok, nil
}
