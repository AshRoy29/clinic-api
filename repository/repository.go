package repository

import (
	"clinic-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataBaseRepo interface {
	InsertAppointment(appointment models.Appointment)
	UpdateAppointment(appointmentID string)
	DeleteAppointment(appointmentID string)
	DeleteAllAppointment() int64
	GetAllAppointments() []primitive.M
	InsertSpecialty(specialty models.Specialties)
	GetAllSpecialties() []primitive.M
	InsertDoctor(doctor models.Doctors)
	GetDoctorsBySpecialties(specialtyID string) []primitive.M
	InsertUser(user models.User)
	CheckEmail(email string) models.User
	AuthEmail(email string) models.User
	GetPatientInfo(patientID string) models.Patient
	GetUserPrescription(patientID string)
	InsertPrescription(patientID string)
	SubmitPrescription(patientID string)
}
