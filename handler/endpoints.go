package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"unicode"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// Login defines model for Login.
type LoginResponse struct {
	Id int32 `json:"id"`
}

type CustomClaims struct {
	Id          int32  `json:"id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	jwt.StandardClaims
}

// Registration implements generated.ServerInterface.
func (s *Server) Registration(ctx echo.Context) error {
	requestBody := new(generated.Registration)
	ctx.Bind(requestBody)

	err := validateRegistrationRequestBody(requestBody)
	if err != "" {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var resp generated.RegistrationResponse
	var id int32 = 321
	resp.Id = &id
	resp.FullName = &requestBody.FullName
	return ctx.JSON(http.StatusOK, resp)
}

// Login implements generated.ServerInterface.
func (*Server) Login(ctx echo.Context) error {
	requestBody := new(generated.Login)
	ctx.Bind(requestBody)

	err := validateLoginRequestBody(requestBody)
	if err != "" {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	response := new(LoginResponse)
	response.Id = 321
	return ctx.JSON(http.StatusOK, response)
}

// GetProfile implements generated.ServerInterface.
func (*Server) GetProfile(ctx echo.Context) error {
	reg := regexp.MustCompile(`^[B|b]earer\s`)
	tokenString := reg.ReplaceAllString(ctx.Request().Header["Authorization"][0], "")
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			error_message := fmt.Sprintf("Expected token algorithm '%v' but got '%v'", jwt.SigningMethodRS256.Name, token.Header)
			return nil, errors.New(error_message)
		}

		verifyBytes, err := os.ReadFile("rsakey.pem.pub")
		if err != nil {
			return nil, err
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			return nil, err
		}

		return verifyKey, nil
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err)
	}

	response := new(generated.ProfileDataResponse)
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		response.FullName = &claims.FullName
		response.PhoneNumber = &claims.PhoneNumber
	}

	return ctx.JSON(http.StatusOK, response)
}

func validateRegistrationRequestBody(requestBody *generated.Registration) (error_message string) {
	errorMessageMap := map[int]string{
		1: "All of the form data must be filled in",
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
		error_message = errorMessageMap[1]
	} else if !validPhoneNumber.MatchString(requestBody.PhoneNumber) {
		error_message = errorMessageMap[2]
	} else if phoneNumberLength < 11 || phoneNumberLength > 14 {
		error_message = errorMessageMap[3]
	} else if fullNameLength < 3 || fullNameLength > 60 {
		error_message = errorMessageMap[4]
	} else if !validLength {
		error_message = errorMessageMap[5]
	} else if !haveNumber {
		error_message = errorMessageMap[6]
	} else if !haveUppercase {
		error_message = errorMessageMap[7]
	} else if !haveSpecialChar {
		error_message = errorMessageMap[8]
	}

	return error_message
}

func validateLoginRequestBody(requestBody *generated.Login) (error_message string) {
	if requestBody.Password == "" || requestBody.PhoneNumber == "" {
		error_message = "All of the form data must be filled in"
	}

	return error_message
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
