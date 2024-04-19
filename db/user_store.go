package db

import (
	"context"

	"github.com/yahialm/GoReserve/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const DBName = "hotel-reservation"
const userColl = "users"

type Dropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	Dropper // Embedded interface

	GetUserByID(context.Context, string) (*types.User, error)
	GetUserByEmail(context.Context, bson.M) (*types.User, error)
	GetUsers(context.Context) (*[]types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	UpdateUser(context.Context, bson.M, types.UpdateUserParams) error
	DeleteUser(context.Context, string) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userColl),
	}
}

func (ms *MongoUserStore) Drop(ctx context.Context) error {
	return ms.coll.Drop(ctx)
}

func (ms *MongoUserStore) GetUserByEmail(ctx context.Context, filter bson.M) (*types.User, error){
	res := ms.coll.FindOne(ctx, filter)
	var u types.User
	if err := res.Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (ms *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error {
	update := bson.D{
		{
			Key: "$set", Value: params.ToBSON(),
		},
	}
	_, err := ms.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// TODO: Maybe it would be good to handle if we did not delete any user
	// Maybe we can log the error
	_, err = ms.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (ms *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	insertedUser, err := ms.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = insertedUser.InsertedID.(primitive.ObjectID) //cast
	return user, nil
}

func (ms *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User
	err = ms.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ms *MongoUserStore) GetUsers(ctx context.Context) (*[]types.User, error) {
	var users = []types.User{}
	curs, err := ms.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = curs.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}
