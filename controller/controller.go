package controller

import (
	. "clinic-api/config"
	"clinic-api/helpers"
	"clinic-api/models"
	. "clinic-api/repository"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
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

//GetAllAppointments displays all the appointments
func GetAllAppointments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	allAppointments := Repo.GetAllAppointments()
	json.NewEncoder(w).Encode(allAppointments)
}

//CreateAppointment creates an appointment
func CreateAppointment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var appointment models.Appointment
	//var doctor models.Doctors
	//var appt models.Appt
	_ = json.NewDecoder(r.Body).Decode(&appointment)

	//for x := range doctor.Appt {
	//	if x ==
	//	for i, v := range appt.Slots {
	//		if v == appointment.Time {
	//			appt.Slots = append(appt.Slots[:i], appt.Slots[i+1:]...)
	//			break
	//		}
	//	}
	//}

	Repo.InsertAppointment(appointment)
	json.NewEncoder(w).Encode(appointment)

}

//UpdateAppointment updates a particular appointment
func UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	Repo.UpdateAppointment(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

//DeleteAppointment deletes a particular appointment
func DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	Repo.DeleteAppointment(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

//DeleteAllAppointment deletes all appointments
func DeleteAllAppointment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	count := Repo.DeleteAllAppointment()
	json.NewEncoder(w).Encode(count)
}

//Specialty

//CreateSpecialty inserts a specialty
func CreateSpecialty(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var specialty models.Specialties
	_ = json.NewDecoder(r.Body).Decode(&specialty)
	Repo.InsertSpecialty(specialty)
	json.NewEncoder(w).Encode(specialty)

}

//GetAllSpecialties displays all specialties
func GetAllSpecialties(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	allSpecialties := Repo.GetAllSpecialties()
	json.NewEncoder(w).Encode(allSpecialties)
}

//Doctors

//CreateDoctor creates a doctor
func CreateDoctor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var appts models.Appt

	var dbdoctor models.Doctors
	var doctor models.Doctors
	_ = json.NewDecoder(r.Body).Decode(&doctor)
	doctor.Image = imageDir

	doctor.Appt = make([]models.Appt, 7)
	for i := 0; i < 7; i++ {
		apptTaken := 0
		slots, appNo := helpers.Time(doctor.StartTime, doctor.EndTime, doctor.Duration)
		appts.Slots = make([]string, appNo+1)
		today := time.Now().AddDate(0, 0, i)
		appts.Slots = slots
		appts.Date = today.Format("02/01/2006")
		fmt.Println(appts.Slots)
		appts.ApptNo = appNo
		appts.ApptTaken = apptTaken
		doctor.Appt[i] = appts
	}

	var err error

	dbdoctor = Repo.CheckEmailDoc(doctor.Email)
	if dbdoctor.Email != "" {
		helpers.ErrorJSON(w, errors.New("email already exists"))
		return
	} else {
		doctor.Password, err = GenerateHashPassword(doctor.FirstName + doctor.Phone)
		if err != nil {
			helpers.ErrorJSON(w, errors.New("error hashing password"))
		}
	}
	//doctor.StartTime = time.Kitchen
	//doctor.EndTime = time.Kitchen
	Repo.InsertDoctor(doctor)
	json.NewEncoder(w).Encode(doctor)
	//json.NewEncoder(w).Encode(slots)
	//json.NewEncoder(w).Encode(appNo)
}

//GetDoctorsBySpecialties displays doctors filtered by a particular specialty
func GetDoctorsBySpecialties(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	doctorBySpecialty := Repo.GetDoctorsBySpecialties(params["id"])
	json.NewEncoder(w).Encode(doctorBySpecialty)
}

//DoctorsByID displays doctors filtered by a particular id
func DoctorsByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	doctorByID := Repo.GetDoctorsByID(params["id"])

	json.NewEncoder(w).Encode(doctorByID)

}

//GetAllDoctors displays all doctors
func GetAllDoctors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	allDoctors := Repo.GetAllDoctors()
	json.NewEncoder(w).Encode(allDoctors)
}

//UpdateDoctor updates a particular doctor by id
func UpdateDoctor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var doctor models.Doctors

	_ = json.NewDecoder(r.Body).Decode(&doctor)

	Repo.UpdateDoctor(doctor)
	json.NewEncoder(w).Encode(doctor)
}

//DeleteDoctor deletes a doctor by id
func DeleteDoctor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	Repo.DeleteDoctor(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

//InsertProfileImage inserts an image for profile picture
func InsertProfileImage(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("content-type", "application/json")
	//w.Header().Set("Allow-Control-Allow-Methods", "POST")

	r.ParseMultipartForm(10 * 1024 * 1024)
	file, handler, err := r.FormFile("image")

	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}
	defer file.Close()

	fmt.Println("file info")
	fmt.Println("file name:", handler.Filename)
	fmt.Println("file size:", handler.Size)
	fmt.Println("file type:", handler.Header.Get("Content-Type"))

	tempFile, err := ioutil.TempFile("img", "img-*.jpg")
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		helpers.ErrorJSON(w, err)
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

//SignUp registers new users
func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var err error

	var user models.User
	var dbuser models.User

	_ = json.NewDecoder(r.Body).Decode(&user)
	dbuser = Repo.CheckEmail(user.Email)
	if dbuser.Email != "" {
		helpers.ErrorJSON(w, errors.New("email already exists"))
		return
	} else {
		user.Password, err = GenerateHashPassword(user.Password)
		if err != nil {
			helpers.ErrorJSON(w, errors.New("error hashing password"))
		}

		user.Prescriptions = make([]string, 0, 1)
		Repo.InsertUser(user)
		json.NewEncoder(w).Encode(user)
	}

}

//UsersByID finds users by their ID
func UsersByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	userByID := Repo.GetUsersByID(params["id"])
	json.NewEncoder(w).Encode(userByID)
}

//GetPrescriptionsByUser filters prescription by user
func GetPrescriptionsByUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	userPrescription := Repo.GetUserPrescriptionByID(params["id"])
	json.NewEncoder(w).Encode(userPrescription)
}

//InsertPrescription creates prescription
func InsertPrescription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	params := mux.Vars(r)
	//var userInfo models.User

	r.ParseMultipartForm(10 * 1024 * 1024)
	file, handler, err := r.FormFile("pres")

	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}
	defer file.Close()

	fmt.Println("file info")
	fmt.Println("file name:", handler.Filename)
	fmt.Println("file size:", handler.Size)
	fmt.Println("file type:", handler.Header.Get("Content-Type"))

	tempFile, err := ioutil.TempFile("pres", "pres-*.pdf")
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		helpers.ErrorJSON(w, err)
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

//GenerateJWT generates JWT
func GenerateJWT(email, role string) (string, error) {
	mySigningKey := []byte("secret")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	log.Println(token)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	log.Println(tokenString)

	if err != nil {
		fmt.Errorf("something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

//CheckPasswordHash checks hashed password with alphanumeric password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//SignIn signs in user
func SignIn(w http.ResponseWriter, r *http.Request) {
	var authDetails models.Authentication
	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	var authuser models.User
	authuser = Repo.AuthEmail(authDetails.Email)
	if authuser.Email == "" {
		helpers.ErrorJSON(w, errors.New("incorrect email"))
		return
	}

	check := CheckPasswordHash(authDetails.Password, authuser.Password)

	if !check {
		helpers.ErrorJSON(w, errors.New("incorrect password"))
		return
	}

	validToken, err := GenerateJWT(authuser.Email, authuser.Role)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	var token models.Token
	token.ID = authuser.ID
	token.Email = authuser.Email
	token.Role = authuser.Role
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

//DoctorSignIn sign in doctor
func DoctorSignIn(w http.ResponseWriter, r *http.Request) {
	var authDetails models.Authentication
	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	var authdoc models.Doctors
	authdoc = Repo.AuthEmailDoc(authDetails.Email)
	if authdoc.Email == "" {
		helpers.ErrorJSON(w, errors.New("incorrect email"))
		return
	}

	check := CheckPasswordHash(authDetails.Password, authdoc.Password)

	if !check {
		helpers.ErrorJSON(w, errors.New("incorrect password"))
		return
	}

	validToken, err := GenerateJWT(authdoc.Email, authdoc.Role)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	var token models.Token
	token.ID = authdoc.ID
	token.Email = authdoc.Email
	token.Role = authdoc.Role
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

//ChangePassword resets password to a new password
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var err error

	var doctor models.Doctors

	_ = json.NewDecoder(r.Body).Decode(&doctor)

	doctor.Password, err = GenerateHashPassword(doctor.Password)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("error hashing password"))
	}

	Repo.DocNewPassword(doctor)
}

//func GetPatientInfo(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("content-type", "application/json")
//	w.Header().Set("Allow-Control-Allow-Methods", "GET")
//
//	params := mux.Vars(r)
//	patientInfo := Repo.GetPatientInfo(params["id"])
//	json.NewEncoder(w).Encode(patientInfo)
//}

//AdminIndex admin welcome page
func AdminIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized."))
		return
	}
	w.Write([]byte("Welcome, Admin."))
}

//UserIndex user welcome page
func UserIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "user" {
		w.Write([]byte("Not Authorized."))
		return
	}

	w.Write([]byte("Welcome, User."))
}

//DoctorIndex doctor welcome page
func DoctorIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "doctor" {
		w.Write([]byte("Not Authorized."))
		return
	}

	w.Write([]byte("Welcome, Doctor."))
}

//func BookNow(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("content-type", "application/json")
//	w.Header().Set("Allow-Control-Allow-Methods", "GET")
//
//	var doctor models.Doctors
//	_ = json.NewDecoder(r.Body).Decode(&doctor)
//
//	slots := Time(doctor.StartTime, doctor.EndTime, 10)
//
//	json.NewEncoder(w).Encode(slots)
//}

//GetAppointmentsByUserID displays appointment filtered by user ID
func GetAppointmentsByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	userAppt := Repo.UserAppointments(params["id"])
	json.NewEncoder(w).Encode(userAppt)
}

//GetAppointmentsByDoctorID displays appointment filtered by doctor ID
func GetAppointmentsByDoctorID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "GET")

	params := mux.Vars(r)
	doctorAppt := Repo.DoctorAppointments(params["id"])
	json.NewEncoder(w).Encode(doctorAppt)
}
