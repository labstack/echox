package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echox/cookbook/twitter/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (h *Handler) CreatePost(c echo.Context) (err error) {
	userID, err := bson.ObjectIDFromHex(userIDFromToken(c))
	if err != nil {
		return err
	}
	u := &model.User{
		ID: userID,
	}
	p := &model.Post{
		ID:   bson.NewObjectID(),
		From: u.ID.Hex(),
	}
	if err = c.Bind(p); err != nil {
		return
	}

	// Validation
	if p.To == "" || p.Message == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid to or message fields"}
	}

	// Find user from database by id
	if err = h.Client.Database("twitter").Collection("users").FindOne(context.TODO(), bson.M{"_id": u.ID}).Decode(u); err != nil {
		return err
	}

	// Save post in database
	if _, err = h.Client.Database("twitter").Collection("posts").InsertOne(context.TODO(), p); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, p)
}

func (h *Handler) FetchPost(c echo.Context) (err error) {
	userID := userIDFromToken(c)
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	// Defaults
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 100
	}

	// Retrieve posts from database
	posts := []*model.Post{}
	cur, err := h.Client.Database("twitter").Collection("posts").
		Find(context.TODO(), bson.M{"to": userID}, options.Find().SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)))
	if err != nil {
		return
	}
	defer cur.Close(context.TODO())

	if err = cur.All(context.TODO(), &posts); err != nil {
		return
	}

	return c.JSON(http.StatusOK, posts)
}
