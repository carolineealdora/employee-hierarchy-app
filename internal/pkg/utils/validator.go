package utils

import (
	"fmt"
)

func GenerateMessageValidatorError(field, validationTag, valueParam string) string {
	switch validationTag {
	case "required":
		return fmt.Sprintf("field %s is required", field)
	}
	return ""
}
