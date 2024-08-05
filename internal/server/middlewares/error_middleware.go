package middlewares

import (
	"errors"
	"net/http"

	"github.com/carolineealdora/employee-hierarchy-app/internal/apperror"
	"github.com/carolineealdora/employee-hierarchy-app/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		listErrors := c.Errors
		if len(listErrors) > 0 {
			e := listErrors.Last().Err
			switch e.(type) {
			case validator.ValidationErrors:
				var ve validator.ValidationErrors
				if errors.As(c.Errors.Last(), &ve) {
					listErrors := make([]apperror.ValidatorError, len(ve))
					for i, val := range ve {
						msg := utils.GenerateMessageValidatorError(val.Field(), val.Tag(), val.Param())
						listErrors[i] = *apperror.NewValidatorError(val.Field(), msg)
					}

					resp := utils.GenerateResponse(apperror.ErrValidator.Error(), listErrors)
					c.AbortWithStatusJSON(http.StatusBadRequest, resp)
				}
			case *apperror.ErrorWrap:
				unwrappedError := e.(*apperror.ErrorWrap).Err.(*apperror.CustomError)
				switch unwrappedError.Code {
				case http.StatusBadRequest:
					Abort(c, unwrappedError)
				case http.StatusInternalServerError:
					Abort(c, unwrappedError)
				case http.StatusNotFound:
					Abort(c, unwrappedError)
				case http.StatusUnauthorized:
					Abort(c, unwrappedError)
				default:
					Abort(c, apperror.InternalServerError())
				}
			}
		}
	}
}

func Abort(c *gin.Context, e *apperror.CustomError) {
	resp := utils.GenerateResponse(e.Message, e.Details)
	c.AbortWithStatusJSON(http.StatusBadRequest, resp)
}
