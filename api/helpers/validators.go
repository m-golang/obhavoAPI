package helpers

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-playground/validator/v10"
)

// RespondWithValidationErrors handles and formats validation errors from request data.
// It takes in a gin.Context, the error from validation, and the struct type for reflecting field names.
// If the error format is invalid, it returns a generic error message to the client.
func RespondWithValidationErrors(c *gin.Context, err error, structType interface{}) {
	// Assert the error as a slice of ValidationErrors (from the validator package)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// If the error format is not a ValidationErrors type, return a bad request response
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid validation error format, please make sure all requested formats are correct",
		})
		return
	}

	var errorMessages []string

	// Reflect over the struct type to access the struct's field names and JSON tags
	structValue := reflect.TypeOf(structType)

	// Iterate over the validation errors and format them into user-friendly error messages
	for _, err := range errs {
		// Get the field info based on the name of the field that caused the validation error
		field, _ := structValue.FieldByName(err.Field())

		fieldName := err.Field()

		// If a JSON tag exists for the field, use it instead of the raw field name
		if tag := field.Tag.Get("json"); tag != "" {
			fieldName = tag
		}

		// Format the error message by combining field name and validation tag (e.g. "field is required")
		errorMessages = append(errorMessages, fmt.Sprintf("'%s' is %s", fieldName, err.Tag()))
	}

	// Return the formatted error messages in the JSON response
	c.JSON(http.StatusBadRequest, gin.H{
		"error": errorMessages,
	})
}

// ValidatePassword ensures that the password meets the following criteria:
// - At least 8 characters long
// - Contains at least one lowercase letter
// - Contains at least one uppercase letter
// - Contains at least one number
// - Contains at least one special character
// It also ensures that the password and confirm password match.
func ValidatePassword(password string) error {

	if len(password) == 0 || len(strings.TrimSpace(password)) == 0 {
		return fmt.Errorf("password cannot be empty or just spaces")
	}

	// Validate password using various rules
	err := validation.Validate(password,
		validation.Length(8, 0).Error("password must be at least 8 characters long"),
		validation.Match(regexp.MustCompile(`[a-z]`)).Error("password must contain at least one lowercase letter"),
		validation.Match(regexp.MustCompile(`[A-Z]`)).Error("password must contain at least one uppercase letter"),
		validation.Match(regexp.MustCompile(`\d`)).Error("password must contain at least one number"),
		validation.Match(regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)).Error("password must contain at least one special character"),
	)
	// Return any validation error
	if err != nil {
		return err
	}

	return nil
}

// GetParametersFromUrl extracts the API key and query parameters from the URL.
// It returns the API key, query parameter, and an error if either is missing or invalid.
func GetParametersFromUrl(c *gin.Context) (string, string, error) {
	// Extract the 'key' parameter (API key) from the URL query string
	apiKey := c.Query("key")
	if len(apiKey) == 0 || len(strings.TrimSpace(apiKey)) == 0 {
		// If the API key is missing or invalid, return an error
		return "", "", fmt.Errorf("api key is missing or invalid. Please include a valid API key in your request")
	}

	// Extract the 'q' parameter (query) from the URL query string
	query := c.Query("q")
	if len(query) == 0 || len(strings.TrimSpace(query)) == 0 {
		// If the query is missing or invalid, return an error
		return "", "", fmt.Errorf("parameter q is missing")
	}

	// Return the API key and query if both are valid
	return apiKey, query, nil
}

// GetParametersFromUrlForBulk extracts the API key and checks if the 'q' parameter is set to 'bulk'.
// It returns the API key and an error if either condition is violated.
func GetParametersFromUrlForBulk(c *gin.Context) (string, error) {
	// Extract the 'key' parameter (API key) from the URL query string
	apiKey := c.Query("key")
	if len(apiKey) == 0 || len(strings.TrimSpace(apiKey)) == 0 {
		// If the API key is missing or invalid, return an error
		return "", fmt.Errorf("api key is missing or invalid. Please include a valid API key in your request")
	}

	// Extract the 'q' parameter (query) and check if it equals 'bulk'
	query := c.Query("q")
	if query != "bulk" {
		// If 'q' is not set to 'bulk', return an error
		return "", fmt.Errorf("parameter q='bulk' is required")
	}

	// Return the API key if it is valid
	return apiKey, nil
}

// FilterValidQValues filters the valid 'q' values from a LocationsForm or similar struct.
// It extracts the 'Q' field from each location and returns a slice of valid non-empty strings.
func FilterValidQValues(data interface{}) []string {
	var qValues []string

	// Get the reflect value of the input data
	val := reflect.ValueOf(data)

	// Check if the input data is a struct, as we're expecting a struct type with a 'Locations' field
	if val.Kind() != reflect.Struct {
		return qValues // Return an empty slice if the data is not a struct
	}

	// Extract the 'Locations' field, which should be a slice of Location structs
	locationsField := val.FieldByName("Locations")
	if locationsField.IsValid() && locationsField.Kind() == reflect.Slice {
		// Iterate over each location in the 'Locations' slice
		for i := 0; i < locationsField.Len(); i++ {
			location := locationsField.Index(i)

			// Extract the 'Q' field from each location (the query string)
			qField := location.FieldByName("Q")
			if qField.IsValid() && qField.Kind() == reflect.String {
				// If the 'Q' field is a non-empty string, append it to the result slice
				qValue := qField.String()
				if qValue != "" && qValue != " " {
					qValues = append(qValues, qValue)
				}
			}
		}
	}

	// Return the slice of valid 'q' values
	return qValues
}
