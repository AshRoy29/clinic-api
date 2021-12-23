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

	route.HandleFunc("/admin", IsAuthorized(AdminIndex)).Methods("GET")
	route.HandleFunc("/user", IsAuthorized(UserIndex)).Methods("GET")

	route.HandleFunc("/api/doctor", controller.CreateDoctor).Methods("POST")
	route.HandleFunc("/api/doctorsp/{id}", controller.GetDoctorsBySpecialties).Methods("GET")

	route.HandleFunc("/api/signup", controller.SignUp).Methods("POST")
	route.HandleFunc("/api/signin", controller.SignIn).Methods("POST")

	route.HandleFunc("/api/pres/{id}", controller.GetPrescriptionsByUser).Methods("GET")

	route.HandleFunc("/api/image", controller.InsertPrescription).Methods("POST")

	route.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})

	return route
}
