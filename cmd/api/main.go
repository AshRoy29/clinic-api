package main

import (
	. "clinic-api/config"
	"fmt"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"time"
)

//type config struct {
//	port int
//	env string
//	db struct{
//		dsn string
//	}
//	jwt struct{
//		secret string
//	}
//}
//
//type AppStatus struct {
//	Status string `json:"status"`
//	Environment string `json:"environment"`
//	Version string `json:"version"`
//}
//
//type application struct {
//	config config
//	logger *log.Logger
//	models models.Models
//}

var config Config

func init() {
	config.Read()
}

//func initFlags() {
//	config.Server.Port = flag.Lookup("Port").Value.(flag.Getter).Get().(string)
//}

func main() {

	//flag.Parse()
	//initFlags()

	//logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//srv := &http.Server{
	//	Addr:         fmt.Sprintf(config.Server.Port),
	//	Handler:      Routes(),
	//	IdleTimeout:  time.Minute,
	//	ReadTimeout:  10 * time.Second,
	//	WriteTimeout: 30 * time.Second,
	//}

	abc := "Hello"
	xyz := "12345"

	fmt.Println(abc + xyz)

	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"})
	methodsOK := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"})
	originOK := handlers.AllowedOrigins([]string{"*"})

	log.Println("Starting server on port", config.Server.Port)
	x := time.Now().Format("02/01/2006 Monday")
	fmt.Println(x)

	log.Fatal(http.ListenAndServe(config.Server.Port, handlers.CORS(headersOK, methodsOK, originOK)(Routes())))

}
