package controller

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ganyacc/ap.git/model"
	"github.com/ganyacc/ap.git/util"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
}

// NewUserHttpHandler returns the instance of UserRepository
func NewUserHttpHandler() UserRepository {
	return &UserHandler{}
}

// Logger Initialization
var logger = echo.New().Logger

// Initialization of In-memory storage and mutex ensure that concurrent operations on the StoreUsers map do not lead to race conditions.
var (
	StoreUsers = make(map[string]model.User)
	mutex      = &sync.Mutex{}
)

// SignUp registers the new user
func (u *UserHandler) SignUp(c echo.Context) error {

	var user *model.User

	//parse the request body
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//validate the request body
	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			logger.Error("Field '%s' failed validation: %s\n", err.Field(), err.Tag())
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Field '%s' failed validation: %s", err.Field(), err.Tag()))
		}
	}

	mutex.Lock()
	defer mutex.Unlock()
	//check if user is already registered
	if _, exists := StoreUsers[user.Email]; exists {
		logger.Infof("user with email '%s' is already registered.", user.Email)
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("user with email '%s' is already registered.", user.Email))
	}

	//hash the password

	password, err := util.HashPassword(user.Password)
	if err != nil {
		logger.Error("error while hashing the password: %v\n", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	//
	user.Password = password

	//store new user in map
	StoreUsers[user.Email] = *user
	return c.JSON(http.StatusCreated, "User registered successfully.")
}

// SignIn handler returns new jwt token and refresh token upon successful
func (u *UserHandler) SignIn(c echo.Context) error {

	user := &model.User{}
	//parse the request body
	err := c.Bind(&user)
	if err != nil {
		logger.Error("error while parsing request body: %v\n", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	if _, exists := StoreUsers[user.Email]; !exists {
		logger.Info("Not a registered Email Id")
		return c.JSON(http.StatusBadRequest, "User not found. Please register to access this service.")
	}

	registeredUser := StoreUsers[user.Email]

	//check password
	if !util.CheckPassword(registeredUser.Password, user.Password) {
		logger.Info("Invalid password for %s user", user.Email)
		return c.JSON(http.StatusBadRequest, "incorrect password.")
	}

	//generate jwt token
	token, err := util.GenerateJwtToken(user.Email)
	if err != nil {
		logger.Error("error generating new token: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	refreshToken, err := util.GenerateRefreshToken(user.Email)
	if err != nil {
		logger.Error("error generating refresh token: ", err)
		return c.JSON(http.StatusInternalServerError, "Error generating refresh token")
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message":       "Signed in successfully.",
		"token":         token,
		"refresh_token": refreshToken,
	})
}

// RefreshToken handler to issue a new JWT token using the refresh token
func (u *UserHandler) RefreshToken(c echo.Context) error {
	refreshToken := c.QueryParam("refresh_token")
	if refreshToken == "" {
		logger.Error("error getting refresh token")
		return c.JSON(http.StatusBadRequest, "Refresh token is required")
	}

	// Parse and validate the refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &util.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return util.RefreshSecretKey, nil
	})

	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, "Invalid or expired refresh token")
	}

	// Extract the email from the claims and generate a new JWT token
	claims, ok := token.Claims.(*util.JwtClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "Invalid token claims")
	}

	newToken, err := util.GenerateJwtToken(claims.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error generating new token")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": newToken,
	})
}
