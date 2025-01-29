package handlers

// newUserForm represents the structure of the data required to create a new user during signup.
// It includes the user's name, surname, username, and password. All fields are required during validation.
type newUserForm struct {
	Name     string `json:"name" binding:"required"`     // The user's first name; must be provided in the request body
	Surname  string `json:"surname" binding:"required"`  // The user's last name; must be provided in the request body
	Username string `json:"username" binding:"required"` // The desired username; must be provided in the request body
	Password string `json:"password" binding:"required"` // The password for the user; must be provided in the request body
}

// userLoginForm represents the structure of the data required for user login.
// It includes the user's username and password for authentication. Both fields are required during validation.
type userLoginForm struct {
	Username string `json:"username" binding:"required"` // The user's username for login; must be provided in the request body
	Password string `json:"password" binding:"required"` // The user's password for login; must be provided in the request body
}
