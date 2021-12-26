package controller

import (
	. "clinic-api/config"
	"clinic-api/models"
	. "clinic-api/repository"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

//IMAGE variables
var dir string
var imageDir string

var userID string
var prescription string
var prescriptions []string

var Repo DBRepo
var config Config

//actual CONTROLLERS

//Appointment

func GetAllAppointments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	allAppointments := Repo.GetAllSpecialties()
	json.NewEncoder(w).Encode(allAppointments)
}

func CreateAppointment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var appointment models.Appointment
	_ = json.NewDecoder(r.Body).Decode(&appointment)
	Repo.InsertAppointment(appointment)
	json.NewEncoder(w).Encode(appointment)

}

func UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	Repo.UpdateAppointment(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	Repo.DeleteAppointment(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllAppointment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	count := Repo.DeleteAllAppointment()
	json.NewEncoder(w).Encode(count)
}

//Specialty

func CreateSpecialty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var specialty models.Specialties
	_ = json.NewDecoder(r.Body).Decode(&specialty)
	Repo.InsertSpecialty(specialty)
	json.NewEncoder(w).Encode(specialty)

}

func GetAllSpecialties(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	allSpecialties := Repo.GetAllSpecialties()
	json.NewEncoder(w).Encode(allSpecialties)
}

//Doctors

func CreateDoctor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var doctor models.Doctors
	_ = json.NewDecoder(r.Body).Decode(&doctor)
	doctor.Image = imageDir
	Repo.InsertDoctor(doctor)
	json.NewEncoder(w).Encode(doctor)

}

func GetDoctorsBySpecialties(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	doctorBySpecialty := Repo.GetDoctorsBySpecialties(params["id"])
	json.NewEncoder(w).Encode(doctorBySpecialty)
}

func DoctorsByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	doctorByID := Repo.GetDoctorsByID(params["id"])
	json.NewEncoder(w).Encode(doctorByID)
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
	dbuser = Repo.CheckEmail(user.Email)
	if dbuser.Email != "" {
		fmt.Println("email already in use")
		return
	} else {
		user.Password, err = GenerateHashPassword(user.Password)
		if err != nil {
			log.Fatalln("error in password hash")
		}

		user.Prescriptions = make([]string, 0, 1)
		Repo.InsertUser(user)
		json.NewEncoder(w).Encode(user)
	}

}

func GetPrescriptionsByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	Repo.GetUserPrescriptionByID(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func InsertPrescription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	//var userInfo models.User

	r.ParseMultipartForm(10 * 1024 * 1024)
	file, handler, err := r.FormFile("pres")

	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	fmt.Println("file info")
	fmt.Println("file name:", handler.Filename)
	fmt.Println("file size:", handler.Size)
	fmt.Println("file type:", handler.Header.Get("Content-Type"))

	tempFile, err := ioutil.TempFile("pres", "pres-*.pdf")
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

	prescriptions := Repo.GetUserPrescriptionByID(params["id"])

	i := cap(prescriptions)

	prescriptions = make([]string, i+1)

	//userInfo.Prescriptions[i] = prescription

	Repo.InsertPrescription(params["id"], prescription)
	json.NewEncoder(w).Encode(params["id"])
}

func GenerateJWT(email, role string) (string, error) {
	var mySigningKey = []byte(config.Jwt.Secret)
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
	authuser = Repo.AuthEmail(authDetails.Email)
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

func GetPatientInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	patientInfo := Repo.GetPatientInfo(params["id"])
	json.NewEncoder(w).Encode(patientInfo)
}
