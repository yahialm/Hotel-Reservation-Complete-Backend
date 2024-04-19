package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yahialm/GoReserve/db"
	"github.com/yahialm/GoReserve/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingHandler struct {
	Store db.Store
}

func NewBookingHandler(store db.Store) *BookingHandler {
	return &BookingHandler{
		Store: store,
	}
}


// TODO: To be admin authorized
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	filter := bson.M{}
	bookings, err := h.Store.Booking.GetBookings(c.Context(), filter)
	if err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"err": "can't get bookings due to internal error"})
	}
	c.JSON(c.Status(fiber.StatusAccepted))
	return c.JSON(bookings)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	param := c.Params("id")
	bookingOID, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"err": "conversion error of id to object ID"})
	}
	filter := bson.M{"_id": bookingOID}
	b, err := h.Store.Booking.GetBookingById(c.Context(), filter)
	if err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"err": "Booking not found"})
	}
	user, ok := c.Context().UserValue("user").(types.Booking)
	if !ok {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"err": "unauthorized"})
	}
	if user.UserID != b.UserID {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"err": "inapropriate action"})
	}
	// Cancelling the fetched booking
	update := bson.M{"$set": bson.M{"canceled": "true"}}
	if err := h.Store.Booking.CancelBooking(c.Context(), filter, update); err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"err": "internal problems encountered"})
	}

	c.JSON(c.Status(fiber.StatusAccepted))
	return c.JSON(fiber.Map{"msg": "booking state has been updated"})
}

// TODO: To be user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	bookID := c.Params("id")
	bookOID, err := primitive.ObjectIDFromHex(bookID)
	if err != nil {
		return c.JSON(fiber.Map{"err": "something goes wrong!"})
	}
	filter := bson.M{"_id": bookOID}
	booking, err := h.Store.Booking.GetBookingById(c.Context(), filter)
	if err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"err": "something goes wrong!"})
	}
	user := c.Context().UserValue("user").(*types.User)
	ok := user.ID == booking.UserID 
	if !ok {
		c.JSON(c.Status(fiber.StatusUnauthorized))
		return c.JSON(fiber.Map{"err": "Not allowed"})
	}
	c.JSON(c.Status(fiber.StatusOK))
	return c.JSON(booking)
}