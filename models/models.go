package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

//type Models struct {
//	DB DBmodel
//	Col Colmodel
//}
//
//func NewModels(db *mongo.Client) Models {
//	return Models{
//		DB: DBmodel{DB: db},
//	}
//}
//
//func Collection(col *mongo.Collection) Models {
//	return Models{
//		Col: Colmodel{Col: col},
//	}
//}

type Appointment struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
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

type Specialties struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name"`
}

type Doctors struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName   string             `json:"first_name"`
	LastName    string             `json:"last_name"`
	Phone       string             `json:"phone"`
	Email       string             `json:"email"`
	Specialties Specialties        `json:"specialties"`
	Degrees     []string           `json:"degrees"`
	Description string             `json:"description"`
	Awards      []string           `json:"awards"`
	Image       string             `json:"image"`
	Role        string             `json:"role"`
	//CreatedAt time.Time `json:"created_at"`
	//UpdatedAt time.Time `json:"updated_at"`
}

//type Patient struct {
//	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
//	FirstName     string             `json:"first_name"`
//	LastName      string             `json:"last_name"`
//	Phone         string             `json:"phone"`
//	Email         string             `json:"email"`
//	Prescriptions []string           `json:"prescriptions"`
//	CreatedAt     time.Time          `json:"created_at"`
//	UpdatedAt     time.Time          `json:"updated_at"`
//}

type User struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName     string             `json:"first_name"`
	LastName      string             `json:"last_name"`
	Phone         string             `json:"phone"`
	Email         string             `json:"email"`
	Password      string             `json:"password"`
	Prescriptions []string           `json:"prescriptions"`
	Role          string             `json:"role"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Role        string             `json:"role"`
	Email       string             `json:"email"`
	TokenString string             `json:"token"`
}

type Role struct {
	Admin  string `json:"admin"`
	Doctor string `json:"doctor"`
	User   string `json:"user"`
}
