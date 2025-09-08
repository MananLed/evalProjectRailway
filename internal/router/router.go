package router

import (
	"net/http"

	"github.com/MananLed/evalProjectRailway/internal/handler"
	"github.com/MananLed/evalProjectRailway/internal/middleware"
	"github.com/MananLed/evalProjectRailway/internal/service"
)

func SetupRouter(userService service.UserService, trainService service.TrainService, ticketService service.TicketService) *http.ServeMux {
	r := http.NewServeMux()

	userHandler := handler.NewUserHandler(&userService)
	trainHandler := handler.NewTrainHandler(&trainService, &userService)
	ticketHandler := handler.NewTicketHandler(&ticketService, &userService, &trainService)

	r.HandleFunc("POST /signup", http.HandlerFunc(userHandler.SignUp)) //Content in body
	r.HandleFunc("POST /login", http.HandlerFunc(userHandler.Login))   //content in body

	r.Handle("GET /profile", middleware.AuthMiddleWare(http.HandlerFunc(userHandler.ViewProfile)))            // GET
	r.Handle("PATCH /profile/update", middleware.AuthMiddleWare(http.HandlerFunc(userHandler.UpdateProfile))) // PATCH  //update details in body
	r.Handle("PATCH /profile/password", middleware.AuthMiddleWare(http.HandlerFunc(userHandler.ChangePassword)))
	r.Handle("DELETE /profile", middleware.AuthMiddleWare(http.HandlerFunc(userHandler.DeleteProfile))) // DELETE

	r.Handle("POST /trains", middleware.AuthMiddleWare(http.HandlerFunc(trainHandler.AddTrain)))
	r.Handle("DELETE /trains/{id}", middleware.AuthMiddleWare(http.HandlerFunc(trainHandler.DeleteTrain)))
	r.Handle("GET /trains", middleware.AuthMiddleWare(http.HandlerFunc(trainHandler.ViewTrain)))

	r.Handle("POST /tickets", middleware.AuthMiddleWare(http.HandlerFunc(ticketHandler.BookTicket)))
	r.Handle("DELETE /tickets/{id}", middleware.AuthMiddleWare(http.HandlerFunc(ticketHandler.CancelTicket)))
	r.Handle("GET /tickets/", middleware.AuthMiddleWare(http.HandlerFunc(ticketHandler.GetTicketOfPassenger)))

	return r
}
