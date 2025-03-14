package http

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func formatErrs(errs validator.ValidationErrors) string {
	msgs := make([]string, 0, len(errs))

	for _, err := range errs {
		field := strings.ToLower(err.Field())

		switch err.ActualTag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("Field '%s' is missing", field))
		case "url":
			msgs = append(msgs, fmt.Sprintf("Field '%s' is not a valid URL", field))
		default:
			msgs = append(msgs, fmt.Sprintf("Field '%s' is not valid", field))
		}
	}

	return strings.Join(msgs, ", ")
}
