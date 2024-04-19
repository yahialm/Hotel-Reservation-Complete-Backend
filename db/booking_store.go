package db

import (
	"context"
	"github.com/yahialm/GoReserve/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookColl = "bookings"

type BookingState struct {
	Canceled bool `json:"canceled"`
}

type BookingStore interface {
	InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error)
	GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error)
	GetBookingById(ctx context.Context, filter bson.M) (*types.Booking, error)
	CancelBooking(ctx context.Context, id bson.M, update bson.M) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(bookColl),
	}
}

func (bs *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	res, err := bs.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = res.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (bs *MongoBookingStore) CancelBooking(ctx context.Context, id bson.M, update bson.M) error {
	// Cancel booking by editing the "canceled" field of the booking
	bookingState := BookingState{
		Canceled: true,
	}
	up := bson.D{
		{
			Key: "$set", Value: bookingState,
		},
	}
	_, err := bs.coll.UpdateOne(ctx, id, up)
	return err
}

func (bs *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	cur, err := bs.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err = cur.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (bs *MongoBookingStore) GetBookingById(ctx context.Context, filter bson.M) (*types.Booking, error) {
	booking := types.Booking{}
	res := bs.coll.FindOne(ctx, filter)
	if err := res.Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}
