package helpers

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
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

// ServerError logs unexpected server errors and returns a generic internal server error response.
// It ensures sensitive information about the error is not exposed to the client.
func ServerError(c *gin.Context, err error) {
	// Log the error on the server for further inspection
	log.Println(err)
	// Send a generic error response to the client
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "An unexpected server error occurred. Please try again later.",
	})
}

// ClientError is used to respond to client errors with a specific HTTP status code and message.
// This can be used for known client-side errors (e.g., 400 for bad requests, 404 for not found).
func ClientError(c *gin.Context, code int, message string) {
	// Return the client error response with the provided status code and message
	c.JSON(code, gin.H{
		"error": message,
	})
}
