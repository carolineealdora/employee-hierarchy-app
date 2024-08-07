package main

import app "github.com/carolineealdora/employee-hierarchy-app/internal/server"

func main() {
	handlers := app.SetupHandlers()

	routers := app.Routers(handlers)
	
	app.RunServer(routers)
}
