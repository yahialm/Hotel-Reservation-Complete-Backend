package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type BookingParams struct{
// 	RoomID primitive.ObjectID `json:"roomID,omitempty"`
// 	UserID primitive.ObjectID `json:"userID,omitempty"`
// 	NumbPersons int `json:"numbPers,omitempty"`
// 	FromDate time.Time `json:"fromDate,omitempty"`
// 	TillDate time.Time `json:"tillDate,omitempty"`
// }

type Booking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RoomID      primitive.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	UserID      primitive.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	NumbPersons int64              `bson:"numbPers,omitempty" json:"numbPers,omitempty"`
	FromDate    time.Time          `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	TillDate    time.Time          `bson:"tillDate,omitempty" json:"tillDate,omitempty"`
	Cancelled   bool               `bson:"canceled" json:"canceled"`
}
