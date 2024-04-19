package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yahialm/GoReserve/db"
	"github.com/yahialm/GoReserve/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	Store db.Store
}

func NewHotelHandler(store db.Store) *HotelHandler {
	return &HotelHandler{
		Store: store,
	}
}

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func (h *HotelHandler) HandleGetRooms(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"hotelID": oid}
	var rooms []*types.Room
	rooms, err = h.Store.Room.GetRooms(ctx.Context(), filter)
	if err != nil {
		return err
	}
	return ctx.JSON(rooms)
}

func (h *HotelHandler) HanldeGetHotels(ctx *fiber.Ctx) error {
	// var queryparams HotelQueryParams
	// err := ctx.QueryParser(&queryparams)
	// if err != nil {
	// 	return err
	// }
	hotels, err := h.Store.Hotel.GetHotels(ctx.Context(), bson.M{})
	if err != nil {
		return err
	}
	ctx.JSON(ctx.Status(fiber.StatusOK))
	return ctx.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(ctx *fiber.Ctx) error {
	param := ctx.Params("id")
	hotel, err := h.Store.Hotel.GetHotelByID(ctx.Context(), param)
	if err != nil {
		return err
	}
	return ctx.JSON(hotel)
}
