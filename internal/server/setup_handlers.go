package app

import (
	"github.com/carolineealdora/employee-hierarchy-app/internal/handlers"
	"github.com/carolineealdora/employee-hierarchy-app/internal/pkg/utils/logger"
	"github.com/carolineealdora/employee-hierarchy-app/internal/repositories"
	"github.com/carolineealdora/employee-hierarchy-app/internal/services"
)

type HandlerOpts struct {
	EmployeeHandler *handlers.EmployeeHandler
}

func SetupHandlers() *HandlerOpts {
	logrusLogger := logger.NewLogger()
	logger.SetLogger(logrusLogger)
	empRepo := repositories.NewEmployeeRepository()

	empService := services.NewEmployeeService(empRepo)
	empHandler := handlers.NewEmployeeHandler(empService)

	return &HandlerOpts{
		EmployeeHandler: empHandler,
	}
}
