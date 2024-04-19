package api

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yahialm/GoReserve/db"
	"github.com/yahialm/GoReserve/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	UserStore db.UserStore
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RespAuthParams struct {
	User *types.User `json:"user"`
	Token string `json:"token"`
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		UserStore: userStore,
	}
}

func (h *AuthHandler) HandleAuth(c *fiber.Ctx) error {
	var loginParams AuthParams
	err := c.BodyParser(&loginParams) 
	if err != nil {
		return err
	}
	filter := bson.M{"email": loginParams.Email}
	u, err :=  h.UserStore.GetUserByEmail(c.Context(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("invalid credentiels")
		}
		return err
	}
	if !types.IsValidPassword(u.EncryptedPassword, loginParams.Password) {
		return fmt.Errorf("invalid credentiels")
	}
	token := createToken(u)
	// TODO: Remove the next line (redirects credentiels into aplication logs)
	fmt.Println("Authenticated -> ", u, "\ntoken -> ", token)
	resAuthParams := RespAuthParams{
		u,
		token,
	}
	return c.JSON(resAuthParams)
}

func createToken(user *types.User) string {
	expires := time.Now().Add(time.Hour * 4)
	expiresUnix := expires.UnixMilli()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, 
		jwt.MapClaims{ 
			"id": user.ID, 
			"email": user.Email,
			"expires": expiresUnix,
		},
	)
	key := os.Getenv("JWT_SECRET")
	token, err := t.SignedString([]byte(key))
	if err != nil {
		fmt.Println(err.Error())
	}
	return token
}