package controller

import (
	"clinic-api/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

const connectionString = "mongodb://localhost:27017"
const dbName = "clinic"

//IMPORTANT
var appointmentCol *mongo.Collection
var doctorsCol *mongo.Collection
var specialtiesCol *mongo.Collection
var userCol *mongo.Collection

//IMAGE variables
var dir string
var imageDir string

var userID string
var prescription string
var prescriptions []string

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

	createDB := client.Database(dbName)
	appointmentCol = createDB.Collection("appointment")
	doctorsCol = createDB.Collection("doctors")
	specialtiesCol = createDB.Collection("specialties")
	userCol = createDB.Collection("users")

	//collection instance
	fmt.Println("Collection instance is ready")
}

//MongoDB helpers

//APPOINTMENT
//insert 1 record
func insertAppointment(appointment models.Appointment) {
	inserted, err := appointmentCol.InsertOne(context.Background(), appointment)
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

	result, err := appointmentCol.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
}

//delete 1 record
func deleteAppointment(appointmentID string) {
	id, _ := primitive.ObjectIDFromHex(appointmentID)
	filter := bson.M{"_id": id}

	deleteCount, err := appointmentCol.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Appointment deleted with delete count: ", deleteCount)
}

//delete all records from MongoDB
func deleteAllAppointment() int64 {
	deleteResult, err := appointmentCol.DeleteMany(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of appointments deleted: ", deleteResult.DeletedCount)

	return deleteResult.DeletedCount
}

//get all appointments from database
func getAllAppointments() []primitive.M {
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

func insertSpecialty(specialty models.Specialties) {
	inserted, err := specialtiesCol.InsertOne(context.Background(), specialty)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted one movie in db with id: ", inserted.InsertedID)
}

func getAllSpecialties() []primitive.M {
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
func insertDoctor(doctor models.Doctors) {
	doctor.Image = imageDir
	inserted, err := doctorsCol.InsertOne(context.Background(), doctor)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted one movie in db with id: ", inserted.InsertedID)
}

func getDoctorsBySpecialties(specialtyID string) []primitive.M {
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

func insertUser(user models.User) {
	inserted, err := userCol.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted one movie in db with id: ", inserted.InsertedID)
}

func checkEmail(email string) models.User {
	filter := bson.M{"email": email}
	var user models.User
	err := userCol.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return user
	}

	return user
}

func authEmail(email string) models.User {
	filter := bson.M{"email": email}
	var authUser models.User
	err := userCol.FindOne(context.Background(), filter).Decode(&authUser)
	if err != nil {
		return authUser
	}

	return authUser
}

func getUserPrescription(patientID string) {
	id, _ := primitive.ObjectIDFromHex(patientID)
	filter := bson.M{"_id": id}
	//projection := bson.D{
	//	{"prescriptions", 1},
	//	{"_id", 0},
	//}
	//err := userCol.FindOne(context.Background(), filter)
	var pres models.User
	err := userCol.FindOne(context.Background(), filter).Decode(&pres)
	if err != nil {
		log.Fatal(err)
	}

	//log.Println(pres.Prescriptions)

}

//func insertPrescription() {
//	//id, _ := primitive.ObjectIDFromHex(patientID)
//	//filter := bson.M{"prescriptions": id}
//	i := len()
//
//	if patientID==userID {
//		prescriptions[i] = prescription
//	}
//}

func submitPrescription(patientID string) {
	id, _ := primitive.ObjectIDFromHex(patientID)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"prescriptions": prescriptions}}

	result, err := appointmentCol.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modified count: ", result.ModifiedCount)

}

//actual CONTROLLERS

//Appointment

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

//Specialty

func CreateSpecialty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var specialty models.Specialties
	_ = json.NewDecoder(r.Body).Decode(&specialty)
	insertSpecialty(specialty)
	json.NewEncoder(w).Encode(specialty)

}

func GetAllSpecialties(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	allSpecialties := getAllSpecialties()
	json.NewEncoder(w).Encode(allSpecialties)
}

//Doctors

func CreateDoctor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var doctor models.Doctors
	_ = json.NewDecoder(r.Body).Decode(&doctor)
	insertDoctor(doctor)
	json.NewEncoder(w).Encode(doctor)

}

func GetDoctorsBySpecialties(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	doctorBySpecialty := getDoctorsBySpecialties(params["id"])
	json.NewEncoder(w).Encode(doctorBySpecialty)
}

func InsertProfileImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	r.ParseMultipartForm(10 * 1024 * 1024)
	file, handler, err := r.FormFile("image")

	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	fmt.Println("file info")
	fmt.Println("file name:", handler.Filename)
	fmt.Println("file size:", handler.Size)
	fmt.Println("file type:", handler.Header.Get("Content-Type"))

	tempFile, err := ioutil.TempFile("img", "img-*.jpg")
	if err != nil {
		log.Println(err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}
	tempFile.Write(fileBytes)

	fmt.Println("SUCCESS")

	fmt.Println(tempFile.Name())

	//files, _ := os.ReadDir(dir)
	path, _ := filepath.Abs(dir)
	//for _, file := range files {
	fmt.Println("path:", path)
	imageDir = filepath.Join(path, tempFile.Name())
	log.Println(imageDir)

	//params := mux.Vars(r)
	//doctorBySpecialty := getDoctorsBySpecialties(params["id"])
	//json.NewEncoder(w).Encode(doctorBySpecialty)
}

//GenerateHashPassword generates hashed password
func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var err error
	var user models.User
	var dbuser models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	dbuser = checkEmail(user.Email)
	if dbuser.Email != "" {
		fmt.Println("email already in use")
		return
	} else {
		user.Password, err = GenerateHashPassword(user.Password)
		if err != nil {
			log.Fatalln("error in password hash")
		}

		insertUser(user)
		json.NewEncoder(w).Encode(user)
	}

}

func GetPrescriptionsByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	getUserPrescription(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func InsertPrescription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	//params := mux.Vars(r)
	//userID = params["id"]

	r.ParseMultipartForm(10 * 1024 * 1024)
	file, handler, err := r.FormFile("image")

	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	fmt.Println("file info")
	fmt.Println("file name:", handler.Filename)
	fmt.Println("file size:", handler.Size)
	fmt.Println("file type:", handler.Header.Get("Content-Type"))

	tempFile, err := ioutil.TempFile("img", "img-*.pdf")
	if err != nil {
		log.Println(err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}
	tempFile.Write(fileBytes)

	fmt.Println("SUCCESS")

	fmt.Println(tempFile.Name())

	//files, _ := os.ReadDir(dir)
	path, _ := filepath.Abs(dir)
	//for _, file := range files {
	fmt.Println("path:", path)
	prescription = filepath.Join(path, tempFile.Name())
	log.Println(prescription)

	//insertPrescription()
	//json.NewEncoder(w).Encode(userID)
}

//func SignUp(w http.ResponseWriter, r *http.Request) {
//	var user models.User
//	err := json.NewDecoder(r.Body).Decode(&user)
//	if err != nil {
//
//		fmt.Println("Error in reading body")
//		w.Header().Set("Content-Type", "application/json")
//		json.NewEncoder(w).Encode(err)
//		return
//	}
//
//
//}

func GenerateJWT(email, role string) (string, error) {
	secretkey := "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160"
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var authDetails models.Authentication
	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		fmt.Println("Error in reading body")
		return
	}

	var authuser models.User
	authuser = authEmail(authDetails.Email)
	if authuser.Email == "" {
		fmt.Println("username or password incorrect")
		return
	}

	check := CheckPasswordHash(authDetails.Password, authuser.Password)

	if !check {
		fmt.Println("username or password incorrect")
		return
	}

	validToken, err := GenerateJWT(authuser.Email, authuser.Password)
	if err != nil {
		fmt.Println("failed to generate token")
		return
	}

	var token models.Token
	token.Email = authuser.Email
	token.Role = authuser.Role
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
