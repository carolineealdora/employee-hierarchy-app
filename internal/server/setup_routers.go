package app

import (
	"net/http"

	"github.com/carolineealdora/employee-hierarchy-app/internal/apperror"
	"github.com/carolineealdora/employee-hierarchy-app/internal/constants"
	middlewares "github.com/carolineealdora/employee-hierarchy-app/internal/server/middlewares"
	"github.com/gin-gonic/gin"
)

func Routers(h *HandlerOpts) http.Handler {
	g := gin.New()
	g.ContextWithFallback = true

	g.NoRoute(func(c *gin.Context) {
		c.Error(apperror.NewError(
			apperror.ErrPageNotFoundError(),
			constants.SetupRoutersFile,
			"setup_routers.Routers",
		))
	})

	g.Use(gin.Recovery(), middlewares.ErrorMiddleware())

	g.POST("/search-employee", h.EmployeeHandler.GetEmployee)

	return g
}
