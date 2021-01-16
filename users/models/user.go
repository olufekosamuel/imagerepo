package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	CreatedAt string             `json:"createdat,omitempty"`
	UpdatedAt string             `json:"updatedat,omitempty"`
}

type Image struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID          primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Image           string             `json:"image" bson:"image"`
	Text            string             `json:"text" bson:"text"`
	Type            string             `json:"type" bson:"type"`
	Characteristics []string           `json:"characteristics" bson:"characteristics"`
	CreatedAt       string             `json:"createdat,omitempty"`
	UpdatedAt       string             `json:"updatedat,omitempty"`
}

type DataResponse map[string]interface{}
type Response struct {
	Status string       `json:"status"`
	Error  bool         `json:"error"`
	Msg    string       `json:"msg,omitempty"`
	Data   DataResponse `json:"data,omitempty"`
}
