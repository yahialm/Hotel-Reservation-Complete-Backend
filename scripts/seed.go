package main

import (
	"context"
	"fmt"
	"log"

	"os"

	"github.com/yahialm/GoReserve/db"
	"github.com/yahialm/GoReserve/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	hotelStore *db.MongoHotelStore
	roomStore *db.MongoRoomStore
	userStore db.UserStore
	ctx = context.Background()
)

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		panic(err)
	}

	if err = client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}

func seedUser(isAdmin bool, fname, lname, email , password string) {
	user, _ := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName: lname,
		Email: email,
		Password: password,
	})
	user.IsAdmin = isAdmin
	u, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal("seeding prob with user")
	}
	fmt.Println(u)
}

func seed(name, location string, rating int) {
	hotel := types.Hotel{
		Name:    name,
		Location: location,
		Rooms: []primitive.ObjectID{},
		Rating: rating,
	}

	room1 := types.Room{
		Type:      types.DeluxRoomType,
		Price: 291,
		Size: "small",
	}

	room2 := types.Room{
		Type:      types.DoubleRoomType,
		Price: 230,
		Size: "large",
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		fmt.Print(err.Error())
	}
	room1.HotelID = insertedHotel.ID
	room2.HotelID = insertedHotel.ID
	_, err = roomStore.InsertRoom(ctx, &room1)
	if err != nil {
		fmt.Print(err.Error())
	}

	_, err = roomStore.InsertRoom(context.TODO(), &room2)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Fprintln(os.Stdout, []any{"seeding the database  with -> ", hotel}...)
	fmt.Fprintln(os.Stdout, []any{"seeding the database  with -> ", room1}...)
	fmt.Fprintln(os.Stdout, []any{"seeding the database  with -> ", room2}...)
}


func main() {
	seed("Ibisos", "Berkane", 3)
	seed("Jorga", "Casa", 5)
	seed("Tayta", "Kenitra", 4)
	seedUser(false, "yahia", "lamhafad", "yahia@gm.com", "normalUserpwd")
	seedUser(true, "admin", "admin", "admin@admin.ad", "adminPassword")
}