package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/yahialm/GoReserve/api"
	"github.com/yahialm/GoReserve/api/middleware"
	"github.com/yahialm/GoReserve/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi = "mongodb://localhost:27017"
)

// TODO: This config is useless till this time since we don't make use of it. Fix that!!
// var config = fiber.Config{
// 	ErrorHandler: func(c *fiber.Ctx, err error) error {
// 		return c.JSON(map[string]string{"error": err.Error()})
// 	},
// }

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		panic(err)
	}

	// Handlers initialization
	var (
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookingStore(client)
		store        = db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		hotelHandler   = api.NewHotelHandler(store)
		userHandler    = api.NewUserHandler(store.User)
		authHandler    = api.NewAuthHandler(store.User)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New()
		apiV1          = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		auth           = app.Group("/api")
		admin          = apiV1.Group("/admin", middleware.AdminAuth)
	)

	// Choose the port in which you run your application using: -port=:<port_number>
	listenAddr := flag.String("port", ":5000", "API Server listen on port ?")
	flag.Parse()

	// Auth handlers
	auth.Post("/auth", authHandler.HandleAuth)

	// User Routes
	apiV1.Post("/user/", userHandler.HandlePostUser)
	apiV1.Get("/user/", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id", userHandler.HandleUpdateUser)

	// Hotel Routes
	apiV1.Get("/hotel", hotelHandler.HanldeGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms) // --> get rooms for hotel

	// room Routes
	apiV1.Post("/room/:id/book", roomHandler.HandleRoomBooking)
	apiV1.Get("/room", roomHandler.HandleGetRooms) // --> get all rooms (not for specific hotel)

	// booking routes
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiV1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// admin routes
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	// Listen on default port 5000
	fmt.Println("----------------------")
	err = app.Listen(*listenAddr)
	fmt.Println(err)
}
