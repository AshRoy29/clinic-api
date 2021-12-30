package repository

import (
	. "clinic-api/config"
	"clinic-api/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type DBRepo struct{}

var config Config

var appointmentCol = new(mongo.Collection)
var doctorsCol = new(mongo.Collection)
var specialtiesCol = new(mongo.Collection)
var usersCol = new(mongo.Collection)

const AppointmentCol = "appointments"
const DoctorsCol = "doctors"
const SpecialtiesCol = "specialties"
const UserCol = "users"

//IMAGE variables
var dir string
var imageDir string

var userID string
var prescription string
var prescriptions []string

func init() {
	config.Read()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Database.Uri))
	if err != nil {
		log.Fatal(err)
	}
	appointmentCol = client.Database(config.Database.DatabaseName).Collection(AppointmentCol)
	doctorsCol = client.Database(config.Database.DatabaseName).Collection(DoctorsCol)
	specialtiesCol = client.Database(config.Database.DatabaseName).Collection(SpecialtiesCol)
	usersCol = client.Database(config.Database.DatabaseName).Collection(UserCol)

}

//MongoDB helpers

//APPOINTMENT
//insert 1 record
func (p *DBRepo) InsertAppointment(appointment models.Appointment) {
	inserted, err := appointmentCol.InsertOne(context.Background(), appointment)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted one movie in db with id: ", inserted.InsertedID)

	//id, _ := primitive.ObjectIDFromHex(appointment.DoctorID)
	filter := bson.M{"_id": appointment.DoctorID, "appt.date": appointment.Date}

	_, err = doctorsCol.UpdateOne(context.Background(), filter, bson.M{"$pull": bson.M{"appt.$.slots": appointment.Time}})
	if err != nil {
		log.Fatal(err)
	}
}

//update 1 record
func (p *DBRepo) UpdateAppointment(appointmentID string) {
	id, _ := primitive.ObjectIDFromHex(appointmentID)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"message": "HELLO"}}

	result, err := appointmentCol.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
}

//delete 1 record
func (p *DBRepo) DeleteAppointment(appointmentID string) {
	id, _ := primitive.ObjectIDFromHex(appointmentID)
	filter := bson.M{"_id": id}

	deleteCount, err := appointmentCol.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Appointment deleted with delete count: ", deleteCount)
}

//delete all records from MongoDB
func (p *DBRepo) DeleteAllAppointment() int64 {
	deleteResult, err := appointmentCol.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of appointments deleted: ", deleteResult.DeletedCount)

	return deleteResult.DeletedCount
}

//get all appointments from database
func (p *DBRepo) GetAllAppointments() []primitive.M {
	cur, err := appointmentCol.Find(context.Background(), bson.D{{}})
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

//SPECIALTY

func (p *DBRepo) InsertSpecialty(specialty models.Specialties) {
	inserted, err := specialtiesCol.InsertOne(context.Background(), specialty)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted one movie in db with id: ", inserted.InsertedID)
}

func (p *DBRepo) GetAllSpecialties() []primitive.M {
	cur, err := specialtiesCol.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var specialties []primitive.M

	for cur.Next(context.Background()) {
		var specialty bson.M
		err := cur.Decode(&specialty)
		if err != nil {
			log.Fatal(err)
		}
		specialties = append(specialties, specialty)
	}

	defer cur.Close(context.Background())

	return specialties
}

//DOCTORS
func (p *DBRepo) InsertDoctor(doctor models.Doctors) {
	inserted, err := doctorsCol.InsertOne(context.Background(), doctor)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted one movie in db with id: ", inserted.InsertedID)
}

func (p *DBRepo) GetDoctorsBySpecialties(specialtyID string) []primitive.M {
	id, _ := primitive.ObjectIDFromHex(specialtyID)
	log.Println(id)
	filter := bson.M{"specialties._id": id}

	cur, err := doctorsCol.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	var doctors []primitive.M

	for cur.Next(context.Background()) {
		var doctor bson.M
		err := cur.Decode(&doctor)
		if err != nil {
			log.Fatal(err)
		}
		doctors = append(doctors, doctor)
	}

	defer cur.Close(context.Background())

	return doctors
}

func (p *DBRepo) GetDoctorsByID(doctorID string) models.Doctors {
	id, _ := primitive.ObjectIDFromHex(doctorID)
	log.Println(id)
	filter := bson.M{"_id": id}
	var doctor models.Doctors

	err := doctorsCol.FindOne(context.Background(), filter).Decode(&doctor)
	if err != nil {
		log.Fatal(err)
	}

	return doctor
}

func (p *DBRepo) UpdateSlots(doctorID string, appts models.Appt) {
	id, _ := primitive.ObjectIDFromHex(doctorID)
	filter := bson.M{"_id": id, "appt.date": appts.Date}

	_, err := doctorsCol.UpdateOne(context.Background(), filter, bson.M{"$set": bson.M{"appt.slots": appts.Slots}})
	if err != nil {
		log.Fatal(err)
	}
}

func (p *DBRepo) InsertUser(user models.User) {
	//user.Prescriptions = make([]string, 1)
	user.CreatedAt = time.Now()
	inserted, err := usersCol.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted one movie in db with id: ", inserted.InsertedID)
}

func (p *DBRepo) GetUsersByID(userID string) models.User {
	id, _ := primitive.ObjectIDFromHex(userID)
	log.Println(id)
	filter := bson.M{"_id": id}
	var user models.User

	err := usersCol.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	return user
}

func (p *DBRepo) CheckEmail(email string) models.User {
	filter := bson.M{"email": email}
	var user models.User
	err := usersCol.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return user
	}

	return user
}

func (p *DBRepo) AuthEmail(email string) models.User {
	filter := bson.M{"email": email}
	var authUser models.User
	err := usersCol.FindOne(context.Background(), filter).Decode(&authUser)
	if err != nil {
		return authUser
	}

	return authUser
}

//func (p *DBRepo) GetPatientInfo(patientID string) models.Patient {
//	id, _ := primitive.ObjectIDFromHex(patientID)
//	filter := bson.M{"_id": id}
//
//	var patient models.Patient
//	err := usersCol.FindOne(context.Background(), filter).Decode(&patient)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return patient
//}

func (p *DBRepo) GetUserPrescriptionByID(patientID string) []string {
	id, _ := primitive.ObjectIDFromHex(patientID)
	filter := bson.M{"_id": id}
	//projection := bson.D{
	//	{"prescriptions", 1},
	//	{"_id", 0},
	//}
	//err := userCol.FindOne(context.Background(), filter)
	var pres models.User
	err := usersCol.FindOne(context.Background(), filter).Decode(&pres)
	if err != nil {
		log.Fatal(err)
	}

	return pres.Prescriptions

}

func (p *DBRepo) InsertPrescription(patientID string, pres string) {
	id, _ := primitive.ObjectIDFromHex(patientID)
	filter := bson.M{"_id": id}

	//svar pres models.User
	//err := usersCol.FindOne(context.Background(), filter).Decode(&pres)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//i := cap(pres.Prescriptions)
	//
	//log.Println(i)
	//log.Println(pres.Prescriptions)
	//log.Println(prescription)
	//pres.Prescriptions = make([]string, i+1)

	//if i==1 {
	//	_, err = userCol.UpdateOne(context.Background(), filter, bson.M{"$set": bson.M{"prescriptions.0": prescription}})
	//} else {
	_, err := usersCol.UpdateOne(context.Background(), filter, bson.M{"$push": bson.M{"prescriptions": pres}})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(pres)
}

func (p *DBRepo) SubmitPrescription(patientID string) {
	id, _ := primitive.ObjectIDFromHex(patientID)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"prescriptions": prescriptions}}

	result, err := appointmentCol.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modified count: ", result.ModifiedCount)

}
