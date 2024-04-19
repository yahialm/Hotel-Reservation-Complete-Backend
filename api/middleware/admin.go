package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yahialm/GoReserve/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok:= c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.JSON(fiber.Map{"err": "not authorized"})
	}
	if !user.IsAdmin {
		return c.JSON(fiber.Map{"err": "not authorized"})
	}
	return c.Next()
}
