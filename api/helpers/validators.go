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

func GetParametersFromUrl(c *gin.Context) (string, string, error) {
	apiKey := c.Query("key")
	if len(apiKey) == 0 || len(strings.TrimSpace(apiKey)) == 0 {
		return "", "", fmt.Errorf("api key is missing or invalid. Please include a valid API key in your request")
	}

	query := c.Query("q")
	if len(query) == 0 || len(strings.TrimSpace(query)) == 0 {

		return "", "", fmt.Errorf("parameter q is missing")
	}

	return apiKey, query, nil
}

func GetParametersFromUrlForBulk(c *gin.Context) (string, error) {
	apiKey := c.Query("key")
	if len(apiKey) == 0 || len(strings.TrimSpace(apiKey)) == 0 {
		return "", fmt.Errorf("api key is missing or invalid. Please include a valid API key in your request")
	}

	query := c.Query("q")
	if query != "bulk" {
		return "", fmt.Errorf("parameter q='bulk' is required")
	}

	return apiKey, nil
}

func FilterValidQValues(data interface{}) []string {
	var qValues []string

	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Struct {
		return qValues
	}

	locationsField := val.FieldByName("Locations")
	if locationsField.IsValid() && locationsField.Kind() == reflect.Slice {
		for i := 0; i < locationsField.Len(); i++ {
			location := locationsField.Index(i)

			qField := location.FieldByName("Q")
			if qField.IsValid() && qField.Kind() == reflect.String {
				qValue := qField.String()
				if qValue != "" && qValue != " " {
					qValues = append(qValues, qValue)
				}
			}
		}
	}

	return qValues
}
