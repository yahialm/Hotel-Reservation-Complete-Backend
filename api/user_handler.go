package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/yahialm/GoReserve/db"
	"github.com/yahialm/GoReserve/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	UserStore db.UserStore //every struct that implement UserStore interface can
	// be used for DB op in the user handler
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		UserStore: userStore,
	}
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var (
		// values = bson.M{}
		params = types.UpdateUserParams{}
		id     = c.Params("id")
	)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	if err := c.BodyParser(&params); err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"error": "Parsing error"})
	}
	if err := h.UserStore.UpdateUser(c.Context(), bson.M{"_id": oid}, params); err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"error": "Can't update user"})
	}
	c.JSON(c.Status(fiber.StatusNoContent))
	return c.JSON(fiber.Map{"updated": id})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.UserStore.DeleteUser(c.Context(), id)
	if err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"error": "Can't delete the given user"})
	}
	return c.JSON(fiber.Map{"msg": "Deleted successfully"})
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"error": "Parsing error"})
	}
	if len(params.ValidateUser()) > 0 {
		c.JSON(c.Status(fiber.StatusBadRequest))
		return c.JSON(params.ValidateUser())
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	insertedUser, err := h.UserStore.InsertUser(c.Context(), user)
	if err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	u, err := h.UserStore.GetUserByID(c.Context(), id)
	if err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(fiber.Map{"error": "Not found"})
		}
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(u)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.UserStore.GetUsers(c.Context())
	if err != nil {
		c.JSON(c.Status(fiber.StatusInternalServerError))
		return c.JSON(fiber.Map{"error": err.Error()})
	}
	c.JSON(c.Status(fiber.StatusAccepted))
	return c.JSON(users)
}
