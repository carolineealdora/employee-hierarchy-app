package utils

import (
	"fmt"
)

func GenerateMessageValidatorError(field, validationTag, valueParam string) string {
	switch validationTag {
	case "required":
		return fmt.Sprintf("field %s is required", field)
	case "min":
		return fmt.Sprintf("field %s needs to meet minimum length(%s)", field, valueParam)
	case "max":
		return fmt.Sprintf("field %s needs to meet minimum length(%s)", field, valueParam)
	case "email":
		return fmt.Sprintf("field %s needs to have the correct email format", field)
	case "contains":
		return fmt.Sprintf("field %s needs to contains : %s", field, valueParam)
	case "containsany":
		return fmt.Sprintf("field %s needs to contains any of these characters : %s", field, valueParam)
	}
	return ""
}
