package handler

import "go.mongodb.org/mongo-driver/v2/mongo"

type (
	Handler struct {
		Client *mongo.Client
	}
)

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)
