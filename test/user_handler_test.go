package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/yahialm/GoReserve/api"
	"github.com/yahialm/GoReserve/db"
	"github.com/yahialm/GoReserve/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const tdbname = "hotel-reservation-test"

type testdb struct{
	db.UserStore
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		t.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func ( tdb *testdb) Teardown(t *testing.T) {
	ctx := context.TODO()
	fmt.Println("--- DROPPING DATABASE")
	tdb.UserStore.Drop(ctx)
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.Teardown(t)

	app := fiber.New()
	userHandler := api.NewUserHandler(tdb.UserStore)

	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName : "james",
		LastName : "foobar",
		Email : "email.a@h.fa",
		Password: "ikzejbFLZYEBFUVK",
	}

	b, err := json.Marshal(params)
	if err != nil {
		t.Errorf(err.Error())
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Errorf(err.Error())
	}
	var user types.User
	err = json.NewDecoder(res.Request.Body).Decode(&user)
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println(user)
	// TODO: Should add some tests for userID
	if user.FirstName != params.FirstName {
		t.Errorf("Expected a different firstname")
	}
	if user.LastName != params.LastName {
		t.Errorf("Expected a different lastname")
	}
	if user.Email != params.Email {
		t.Errorf("Expected a different email")
	} 
}

// TODO: Must test the other handlers (GET, PUT and DELETE)