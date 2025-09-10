package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MananLed/evalProjectRailway/internal/config"
	"github.com/MananLed/evalProjectRailway/internal/db"
	"github.com/MananLed/evalProjectRailway/internal/middleware"
	"github.com/MananLed/evalProjectRailway/internal/repository"
	"github.com/MananLed/evalProjectRailway/internal/router"
	"github.com/MananLed/evalProjectRailway/internal/service"
	"github.com/MananLed/evalProjectRailway/pkg/logger"
)

func main() {

	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
		logger.LogToFile(fmt.Sprintf("Error: %v", err))
	}
	defer database.Close()
	if err := db.RunInitialSetup(database); err != nil {
		log.Fatal(err)
		logger.LogToFile(fmt.Sprintf("Error: %v", err))
	}
	userRepo := repository.NewUserRepository(database)
	trainRepo := repository.NewTrainRepository(database)
	ticketRepo := repository.NewTicketRepository(database)

	userService := service.NewUserService(userRepo)
	trainService := service.NewTrainService(trainRepo)
	ticketService := service.NewTicketService(ticketRepo)

	router := router.SetupRouter(*userService, *trainService, *ticketService)

	handler := middleware.LoggingMiddleWare(router)

	log.Println("Server starting on the Port 8080...")
	fmt.Println("Server starting on the Port 8080...")

	err = http.ListenAndServe(config.PortForAPI, handler)
	if err != nil {
		log.Fatal(err)
	}
}
