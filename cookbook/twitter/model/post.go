package model

import "go.mongodb.org/mongo-driver/v2/bson"

type (
	Post struct {
		ID      bson.ObjectID `json:"id" bson:"_id,omitempty"`
		To      string        `json:"to" bson:"to"`
		From    string        `json:"from" bson:"from"`
		Message string        `json:"message" bson:"message"`
	}
)
