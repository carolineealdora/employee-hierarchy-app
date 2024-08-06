package middlewares

import (
	"errors"
	"log"
	"net/http"

	"github.com/carolineealdora/employee-hierarchy-app/internal/pkg/apperror"
	"github.com/carolineealdora/employee-hierarchy-app/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		listErrors := c.Errors
		log.Println(listErrors, "listss")
		if len(listErrors) > 0 {
			// e := listErrors.Last().Err
			e := listErrors[0].Err
			log.Println(e, "nas")
			switch e.(type) {
			case validator.ValidationErrors:
				var ve validator.ValidationErrors
				log.Println(listErrors, "lll")
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
