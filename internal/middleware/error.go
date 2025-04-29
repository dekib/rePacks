package middleware

import (
	"github.com/dekib/rePacks/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err

			var errorResp *errors.ErrorResponse
			switch e := err.(type) {
			case *errors.ErrorResponse:
				errorResp = e
			case validator.ValidationErrors:
				errorResp = errors.NewBadRequestError(
					errors.ErrValidation,
					e,
					nil,
				)
			default:
				errorResp = errors.NewInternalServerError(
					errors.ErrInternalServer,
					e,
					nil,
				)
			}

			c.JSON(errorResp.StatusCode, gin.H{
				"code":       errorResp.Code,
				"message":    errorResp.Message,
				"details":    errorResp.Variables,
				"validation": errorResp.Validation,
			})
		}
	}
}
