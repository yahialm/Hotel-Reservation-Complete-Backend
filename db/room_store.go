package db

import (
	"context"

	"github.com/yahialm/GoReserve/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
)

const RoomsColl = "rooms"

type RoomStore interface {
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
	// GetHotels(context.Context) (*[]types.User, error)
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	// UpdateHotel(context.Context, bson.M, types.UpdateUserParams) (error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll *mongo.Collection
	hotelStore HotelStore
}

//										dependency injection
// 											............
func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore{
	return &MongoRoomStore{
		client: client,
		coll: client.Database(DBNAME).Collection(RoomsColl),
		hotelStore: hotelStore,
	}
}

func (mrs *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	curs, err  := mrs.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := curs.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil 
}

func (mrs *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error){
	res, err := mrs.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	
	err = mrs.hotelStore.UpdateHotel(ctx, filter, update)
	if err != nil {
		return room, err
	}
	return room, nil
}