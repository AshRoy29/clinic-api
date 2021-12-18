package controller

import (
	"clinic-api/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

const connectionString = "mongodb://localhost:27017"
const dbName = "clinic"
const colName = "appointment"

//IMPORTANT
var collection *mongo.Collection

//connect with MongoDB

func init() {
	//client options
	clientOptions := options.Client().ApplyURI(connectionString)

	//connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mongo connection success")

	collection = client.Database(dbName).Collection(colName)

	//collection instance
	fmt.Println("Collection instance is ready")
}

//MongoDB helpers

//insert 1 record
func insertAppointment(appointment models.Appointment) {
	inserted, err := collection.InsertOne(context.Background(), appointment)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted one movie in db with id: ", inserted.InsertedID)
}

//update 1 record
func updateAppointment(appointmentID string) {
	id, _ := primitive.ObjectIDFromHex(appointmentID)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"message": "HELLO"}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
}

//delete 1 record
func deleteAppointment(appointmentID string) {
	id, _ := primitive.ObjectIDFromHex(appointmentID)
	filter := bson.M{"_id": id}

	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Appointment deleted with delete count: ", deleteCount)
}

//delete all records from MongoDB
func deleteAllAppointment() int64 {
	deleteResult, err := collection.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of appointments deleted: ", deleteResult.DeletedCount)

	return deleteResult.DeletedCount
}

//get all appointments from database
func getAllAppointments() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var appointments []primitive.M

	for cur.Next(context.Background()) {
		var appointment bson.M
		err := cur.Decode(&appointment)
		if err != nil {
			log.Fatal(err)
		}
		appointments = append(appointments, appointment)
	}

	defer cur.Close(context.Background())

	return appointments
}

//actual controllers

func GetAllAppointments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	allAppointments := getAllAppointments()
	json.NewEncoder(w).Encode(allAppointments)
}

func CreateAppointment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var appointment models.Appointment
	_ = json.NewDecoder(r.Body).Decode(&appointment)
	insertAppointment(appointment)
	json.NewEncoder(w).Encode(appointment)

}

func UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	updateAppointment(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteAppointment(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllAppointment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	count := deleteAllAppointment()
	json.NewEncoder(w).Encode(count)
}
