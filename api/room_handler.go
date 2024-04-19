package api

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yahialm/GoReserve/db"
	"github.com/yahialm/GoReserve/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	Store db.Store
}

func NewRoomHandler(store db.Store) *RoomHandler {
	return &RoomHandler{
		Store: store,
	}
}

func (bookingParams *BookingParams) validate() error {
	now := time.Now()
	if bookingParams.FromDate.Before(now) || bookingParams.TillDate.Before(now) {
		return fmt.Errorf("choose a correct date")
	}
	return nil
}

type BookingParams struct{
	FromDate time.Time `json:"fromDate"`
	TillDate time.Time `json:"tillDate"`
	NumbPersons int `json:"numbPers"`
}

// Get All rooms, use bson.M{}
func (rh *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := rh.Store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"msg": "something wrong! try again"})
	}
	c.JSON(c.Status(fiber.StatusAccepted))
	return c.JSON(rooms)
}

func (rh *RoomHandler) HandleRoomBooking(c *fiber.Ctx) error {
	var bookingParams BookingParams
	err := c.BodyParser(&bookingParams)
	if err != nil {
		return c.JSON(err.Error())
	}
	if bookingParams.validate() != nil {
		return c.JSON(fiber.Map{"msg": "invalid date"})
	}
	roomID := c.Params("id")

	user := c.Context().Value("user").(*types.User)
	roomOID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}

	//...
	ok, err := rh.isRoomOpenForBooking(c, roomOID, bookingParams)
	if err != nil {
		return c.JSON(c.Status(fiber.StatusInternalServerError))
	}
	if !ok {
		c.JSON(c.Status(fiber.StatusOK))
		return c.JSON(fiber.Map{"msg": "This room is already booked"})
	}

	booking := &types.Booking{
		FromDate: bookingParams.FromDate,
		TillDate: bookingParams.TillDate,
		RoomID: roomOID,
		UserID: user.ID,
		NumbPersons: int64(bookingParams.NumbPersons),
	}
	// Insert booking if no book for the same room exist in this period
	b, err := rh.Store.Booking.InsertBooking(c.Context(), booking)
	if err != nil {
		return c.JSON(fmt.Errorf("can't book the room, try again"))
	}
	return c.JSON(b)
}




// ------------------------------------------------------------------------------

// Helper function
func (rh *RoomHandler) isRoomOpenForBooking(c *fiber.Ctx, roomOID primitive.ObjectID, bookingParams BookingParams) (bool, error) {
	//Check for any booking for this room in same period
	filter := bson.M{
		"roomID": roomOID,
		"fromDate": bson.M{
			"$gte": bookingParams.FromDate,
		},
		"tillDate": bson.M{
			"$lte": bookingParams.TillDate,
		},
	}
	bookings, err := rh.Store.Booking.GetBookings(c.Context(), filter)
	if err != nil {
		return false, err
	}
	if len(bookings) > 0 {
		return false, nil
	}
	return true, nil
}