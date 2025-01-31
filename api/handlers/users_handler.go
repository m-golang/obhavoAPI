package handlers

import (
	"errors"
	"fmt"
	"havoAPI/api/helpers"
	"havoAPI/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler is a struct that holds the service for user-related operations.
type UserHandler struct {
	user services.UsersServiceInterface // Interface to interact with the user service layer
}

// NewUsersHandler creates a new instance of UserHandler with the provided user service.
// This is typically called when setting up the handler for routing.
func NewUsersHandler(user services.UsersServiceInterface) *UserHandler {
	return &UserHandler{user: user}
}

// Signup handles the user signup process.
// It expects a JSON body with user details and performs validation, password checks, and user creation.
// Responds with appropriate errors or success message based on the signup outcome.
func (service *UserHandler) Signup(c *gin.Context) {
	var newUser newUserForm

	// Bind incoming JSON data to the newUser form
	if err := c.ShouldBindJSON(&newUser); err != nil {
		// If binding fails, respond with validation errors
		helpers.RespondWithValidationErrors(c, err, newUser)
		return
	}

	// Validate the password (e.g., length, complexity)
	if err := helpers.ValidatePassword(newUser.Password); err != nil {
		// If the password is invalid, respond with a client error
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	// Attempt to insert the new user into the database
	err := service.user.InsertNewUser(newUser.Name, newUser.Surname, newUser.Username, newUser.Password)
	if err != nil {
		// Handle case when the username already exists
		if errors.Is(err, services.ErrUsernameExists) {
			helpers.ClientError(c, http.StatusConflict, "Username already exists. Consider using a different one or check if you already have an account.")
			return
		}
		// If another error occurs, respond with a server error
		helpers.ServerError(c, err)
		return
	}

	// Return a success response after successful user creation
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully! Please login now",
	})
}

// Login handles the user login process.
// It expects a JSON body with username and password, then validates and authenticates the user.
// If successful, a JWT token is generated and sent as a cookie.
func (service *UserHandler) Login(c *gin.Context) {
	var userLogin userLoginForm

	// Bind incoming JSON data to the userLogin form
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		// If binding fails, respond with validation errors
		helpers.RespondWithValidationErrors(c, err, userLogin)
		return
	}

	// Authenticate the user by checking the username and password
	userID, err := service.user.UserAuthentication(userLogin.Username, userLogin.Password)
	if err != nil {
		// Handle cases for user not found or invalid credentials
		if errors.Is(err, services.ErrUserNotFound) {
			helpers.ClientError(c, http.StatusNotFound, "User not found")
			return
		}
		if errors.Is(err, services.ErrInvalidUserCredentials) {
			helpers.ClientError(c, http.StatusUnauthorized, "Invalid user credentials")
			return
		}

		// For any other errors, respond with a server error
		helpers.ServerError(c, err)
		return
	}

	// Create and sign a JWT token for the authenticated user
	tokenString, err := helpers.CreateAndSignJWT(userID)
	if err != nil {
		// Respond with a server error if JWT creation fails
		helpers.ServerError(c, err)
		return
	}

	// Set the JWT token as a cookie in the response
	helpers.SetCookie(c, tokenString)

	// Return a success response after successful login
	c.JSON(http.StatusOK, gin.H{
		"message": "Login complete! Explore what's new!",
	})
}

// Logout handles user logout by clearing the JWT token from the client's cookies.
// It sends a success message once the token is removed.
func (service *UserHandler) Logout(c *gin.Context) {
	// Clear the JWT token stored in the "u_auth" cookie
	c.SetCookie("u_auth", "", -1, "", "", false, true)

	// Return a success response after logout
	c.JSON(http.StatusOK, gin.H{
		"message": "You are now logged out. Have a great day!",
	})
}

// UserDashboard fetches the user's API key and returns it in the response.
// The user must be authenticated and the ID is extracted from the context.
func (service *UserHandler) UserDashboard(c *gin.Context) {
	// Get the userID from the context (which should have been set during authentication)
	userID, _ := c.Get("userID")
	user_id := int(userID.(float64))

	// Fetch the API key for the authenticated user
	apiKey, err := service.user.FetchUserAPIKey(user_id)
	if err != nil {
		helpers.ServerError(c, err)
		return
	}

	// Return the API key in the response
	c.JSON(http.StatusOK, gin.H{
		"Your API key": apiKey,
	})
}
