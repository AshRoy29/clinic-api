package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Appointment struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Phone     string             `json:"phone"`
	Email     string             `json:"email"`
	Date      time.Time          `json:"date"`
	Time      time.Time          `json:"time"`
	Message   string             `json:"message"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
