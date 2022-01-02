package main

import (
	"clinic-api/controller"
	"github.com/gorilla/mux"
	"net/http"
)

func Routes() *mux.Router {
	route := mux.NewRouter()

	route.HandleFunc("/api/appointments", controller.GetAllAppointments).Methods("GET")
	route.HandleFunc("/api/appointment", controller.CreateAppointment).Methods("POST")
	route.HandleFunc("/api/appointment/{id}", controller.UpdateAppointment).Methods("PUT")
	route.HandleFunc("/api/appointment/{id}", controller.DeleteAppointment).Methods("DELETE")
	route.HandleFunc("/api/deleteallapps", controller.DeleteAllAppointment).Methods("DELETE")

	route.HandleFunc("/api/specialties", controller.GetAllSpecialties).Methods("GET")
	route.HandleFunc("/api/specialty", controller.CreateSpecialty).Methods("POST")

	route.HandleFunc("/admin", IsAuthorized(controller.AdminIndex)).Methods("GET")
	route.HandleFunc("/user", IsAuthorized(controller.UserIndex)).Methods("GET")
	route.HandleFunc("/doctor", IsAuthorized(controller.DoctorIndex)).Methods("GET")

	route.HandleFunc("/api/doctor", controller.CreateDoctor).Methods("POST")
	route.HandleFunc("/api/doctor/{id}", controller.DoctorsByID).Methods("GET")
	route.HandleFunc("/api/doctorsp/{id}", controller.GetDoctorsBySpecialties).Methods("GET")
	route.HandleFunc("/api/doctor/update/{id}", controller.UpdateDoctor).Methods("POST")
	route.HandleFunc("/api/doctor/appointments/{id}", controller.GetAppointmentsByDoctorID).Methods("GET")
	route.HandleFunc("/api/doctors", controller.GetAllDoctors).Methods("GET")
	route.HandleFunc("/api/ddoctor/{id}", controller.DeleteDoctor).Methods("DELETE")
	route.HandleFunc("/api/doctor/signin", controller.DoctorSignIn).Methods("POST")
	route.HandleFunc("/api/doctor/change", controller.ChangePassword).Methods("POST")

	route.HandleFunc("/api/signup", controller.SignUp).Methods("POST")
	route.HandleFunc("/api/signin", controller.SignIn).Methods("POST")

	route.HandleFunc("/api/patient/{id}", controller.UsersByID).Methods("GET")
	route.HandleFunc("/api/prescriptions/{id}", controller.GetPrescriptionsByUser).Methods("GET")
	route.HandleFunc("/api/patient/appointments/{id}", controller.GetAppointmentsByUserID).Methods("GET")

	route.HandleFunc("/api/image", controller.InsertProfileImage).Methods("POST")
	route.HandleFunc("/api/prescription/{id}", controller.InsertPrescription).Methods("PUT")

	//route.HandleFunc("/time", controller.Time).Methods("GET")

	route.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})

	return route
}
