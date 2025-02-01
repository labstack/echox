package handler

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echox/cookbook/twitter/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handler) Signup(c echo.Context) (err error) {
	// Bind
	u := &model.User{ID: bson.NewObjectID()}
	if err = c.Bind(u); err != nil {
		return
	}

	// Validate
	if u.Email == "" || u.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid email or password"}
	}

	// Save user
	if _, err = h.Client.Database("twitter").Collection("users").InsertOne(c.Request().Context(), u); err != nil {
		return
	}

	return c.JSON(http.StatusCreated, u)
}

func (h *Handler) Login(c echo.Context) (err error) {
	// Bind
	u := new(model.User)
	if err = c.Bind(u); err != nil {
		return
	}

	// Find user
	cursor, err := h.Client.Database("twitter").Collection("users").
		Find(c.Request().Context(), bson.M{"email": u.Email, "password": u.Password})
	if err != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid email or password"}
	}
	defer cursor.Close(c.Request().Context())
	if cursor.Next(c.Request().Context()) {
		if err = cursor.Decode(u); err != nil {
			return err
		}
	} else {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid email or password"}
	}

	//-----
	// JWT
	//-----

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response
	u.Token, err = token.SignedString([]byte(Key))
	if err != nil {
		return err
	}

	u.Password = "" // Don't send password
	return c.JSON(http.StatusOK, u)
}

func (h *Handler) Follow(c echo.Context) (err error) {
	userID := userIDFromToken(c)
	id := c.Param("id")

	// Add a follower to user
	targetUserID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	if _, err = h.Client.Database("twitter").
		Collection("users").
		UpdateOne(c.Request().Context(), bson.M{"_id": targetUserID}, bson.M{"$addToSet": bson.M{"followers": userID}}); err != nil {
		return err
	}

	return
}

func userIDFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["id"].(string)
}
