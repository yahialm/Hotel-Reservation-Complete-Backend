package db

import (
	"context"

	"github.com/yahialm/GoReserve/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelColl = "hotels"

type HotelStore interface {
	GetHotelByID(context.Context, string) (*types.Hotel, error)
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(context.Context, bson.M, bson.M) error
	GetHotels(context.Context, bson.M) ([]*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(hotelColl),
	}
}

func (mhs *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	res, err := mhs.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := res.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, err
}

func (mhs *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	var (
		fetchedHotel = types.Hotel{}
		oid, err     = primitive.ObjectIDFromHex(id)
	)
	if err != nil {
		return nil, err
	}
	res := mhs.coll.FindOne(ctx, bson.M{"_id": oid})
	if err := res.Decode(&fetchedHotel); err != nil {
		return nil, err
	}
	return &fetchedHotel, nil
}

func (mhs *MongoHotelStore) InsertHotel(ctx context.Context, h *types.Hotel) (*types.Hotel, error) {
	res, err := mhs.coll.InsertOne(ctx, h)
	if err != nil {
		return nil, err
	}
	h.ID = res.InsertedID.(primitive.ObjectID)
	return h, nil
}

func (mhs *MongoHotelStore) UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := mhs.coll.UpdateOne(ctx, filter, update)
	return err
}

// func (mhs *MongoHotelStore) DeleteHotel(ctx context.Context, id string) (error) {
// 	oid, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = mhs.coll.DeleteOne(ctx, bson.M{"_id": oid})
// 	if err != nil {
// 		return err
// 	}
// 	// TODO: Check if i > 0 to confirm delete ops.
// 	// i := res.DeletedCount
// 	return nil
// }
