package repository

import (
	. "clinic-api/config"
	"clinic-api/helpers"
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

//MongoDB Repository

//InsertAppointment inserts an appointment in MongoDB database
func (p *DBRepo) InsertAppointment(appointment models.Appointment) {
	inserted, err := appointmentCol.InsertOne(context.Background(), appointment)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted one appointment in db with id: ", inserted.InsertedID)

	filter := bson.M{"_id": appointment.DoctorID, "appt.date": appointment.Date}

	_, err = doctorsCol.UpdateOne(context.Background(), filter, bson.M{"$pull": bson.M{"appt.$.slots": appointment.Time}})
	if err != nil {
		log.Fatal(err)
	}
}

//UpdateAppointment updates a particular appointment
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

//DeleteAppointment deletes a particular appointment
func (p *DBRepo) DeleteAppointment(appointmentID string) {
	id, _ := primitive.ObjectIDFromHex(appointmentID)
	filter := bson.M{"_id": id}

	deleteCount, err := appointmentCol.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Appointment deleted with delete count: ", deleteCount)
}

//DeleteAllAppointment deletes all appointments from database
func (p *DBRepo) DeleteAllAppointment() int64 {
	deleteResult, err := appointmentCol.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of appointments deleted: ", deleteResult.DeletedCount)

	return deleteResult.DeletedCount
}

//GetAllAppointments retrieves all appointments from database
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

//InsertSpecialty inserts a specialty in database
func (p *DBRepo) InsertSpecialty(specialty models.Specialties) {
	inserted, err := specialtiesCol.InsertOne(context.Background(), specialty)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted one movie in db with id: ", inserted.InsertedID)
}

//GetAllSpecialties displays all specialties
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

//DeleteSpecialty deletes a particular specialty
func (p *DBRepo) DeleteSpecialty(specialtyID string) {
	id, _ := primitive.ObjectIDFromHex(specialtyID)
	filter := bson.M{"_id": id}

	deleteCount, err := specialtiesCol.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Appointment deleted with delete count: ", deleteCount)
}

//DOCTORS

//InsertDoctor inserts a doctor in database
func (p *DBRepo) InsertDoctor(doctor models.Doctors) {
	inserted, err := doctorsCol.InsertOne(context.Background(), doctor)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted one movie in db with id: ", inserted.InsertedID)
}

//GetDoctorsBySpecialties displays doctors based on a particular specialty
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

//GetDoctorsByID retrieves a doctor from database by id
func (p *DBRepo) GetDoctorsByID(doctorID string) models.Doctors {
	id, _ := primitive.ObjectIDFromHex(doctorID)
	log.Println(id)
	filter := bson.M{"_id": id}
	var doctor models.Doctors
	var appts models.Appt

	err := doctorsCol.FindOne(context.Background(), filter).Decode(&doctor)
	if err != nil {
		log.Fatal(err)
	}
	var index int

	if doctor.Appt[0].Date == time.Now().Format("02/01/2006") {
		return doctor
	} else {
		for j, v := range doctor.Appt {
			//fmt.Println(j, v)
			if v.Date == time.Now().Format("02/01/2006") {
				index = j
				break
			} else {
				index = 7
				break
			}
		}
		fmt.Println(index)
		fmt.Println(doctor.Appt[0].Date)
		for i := 0; i != index; i++ {
			_, err = doctorsCol.UpdateOne(context.Background(), filter, bson.M{"$pull": bson.M{"appt": doctor.Appt[i]}})
			if err != nil {
				log.Fatal(err)
			}

			slots, appNo := helpers.Time(doctor.StartTime, doctor.EndTime, doctor.Duration)
			today := time.Now().AddDate(0, 0, 7-(index-i))
			appts.Slots = slots
			appts.Date = today.Format("02/01/2006")
			fmt.Println(appts.Slots)
			appts.ApptNo = appNo

			_, err = doctorsCol.UpdateOne(context.Background(), filter, bson.M{"$push": bson.M{"appt": appts}})
		}

		return doctor
	}
}

//GetAllDoctors displays all doctors
func (p *DBRepo) GetAllDoctors() []primitive.M {
	cur, err := doctorsCol.Find(context.Background(), bson.D{{}})
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

//UpdateDoctor updates doctor information
func (p *DBRepo) UpdateDoctor(doctor models.Doctors) {
	filter := bson.M{"_id": doctor.ID}
	update := bson.M{"$set": bson.M{"phone": doctor.Phone, "start_time": doctor.StartTime, "end_time": doctor.EndTime, "duration": doctor.Duration}}

	_, err := doctorsCol.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

}

//DeleteDoctor deletes doctor
func (p *DBRepo) DeleteDoctor(doctorID string) {
	id, _ := primitive.ObjectIDFromHex(doctorID)
	filter := bson.M{"_id": id}

	deleteCount, err := doctorsCol.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Appointment deleted with delete count: ", deleteCount)
}

//InsertUser creates user
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

func (p *DBRepo) CheckEmailDoc(email string) models.Doctors {
	filter := bson.M{"email": email}
	var doctor models.Doctors
	err := doctorsCol.FindOne(context.Background(), filter).Decode(&doctor)
	if err != nil {
		return doctor
	}

	return doctor
}

func (p *DBRepo) AuthEmailDoc(email string) models.Doctors {
	filter := bson.M{"email": email}
	var authDoc models.Doctors
	err := doctorsCol.FindOne(context.Background(), filter).Decode(&authDoc)
	if err != nil {
		return authDoc
	}

	return authDoc
}

func (p *DBRepo) DocNewPassword(authdoc models.Doctors) {
	filter := bson.M{"email": authdoc.Email}
	update := bson.M{"$set": bson.M{"password": authdoc.Password}}

	_, err := doctorsCol.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *DBRepo) DeleteUser(userID string) {
	id, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"_id": id}

	deleteCount, err := usersCol.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Appointment deleted with delete count: ", deleteCount)
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

func (p *DBRepo) DeletePrescription(patientID string, presName string) {
	id, _ := primitive.ObjectIDFromHex(patientID)
	filter := bson.M{"_id": id, "prescriptions": presName}

	deleteCount, err := usersCol.UpdateOne(context.Background(), filter, bson.M{"$pull": bson.M{"prescriptions": presName}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Appointment deleted with delete count: ", deleteCount)
}

func (p *DBRepo) UserAppointments(patientID string) []primitive.M {
	id, _ := primitive.ObjectIDFromHex(patientID)
	filter := bson.M{"userid": id}

	cur, err := appointmentCol.Find(context.Background(), filter)
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

func (p *DBRepo) DoctorAppointments(doctorID string) []primitive.M {
	id, _ := primitive.ObjectIDFromHex(doctorID)
	filter := bson.M{"doctorid": id}

	cur, err := appointmentCol.Find(context.Background(), filter)
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
