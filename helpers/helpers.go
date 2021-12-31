package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Time(startTime, endTime string, interval time.Duration) ([]string, int) {
	tStart, err := time.Parse("15:04", startTime)
	tEnd, err := time.Parse("15:04", endTime)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tStart.Format(time.Kitchen))
	fmt.Println(tEnd.Format(time.Kitchen))
	//json.NewEncoder(w).Encode(tEnd.Format(time.Kitchen))
	slots := []string{tStart.Format(time.Kitchen)}
	count := 0
	startP := tStart

	for startP != tEnd {
		d := startP.Add(time.Minute * interval)
		startP = d
		x := d.Format(time.Kitchen)

		slots = append(slots, x)

		count++

	}
	fmt.Println(slots)
	fmt.Println(count)
	//json.NewEncoder(w).Encode(slots)

	return slots, count

	//app := "4:20PM"
	//
	//for i, v := range slots {
	//	if v == app {
	//		slots = append(slots[:i], slots[i+1:]...)
	//		break
	//	}
	//}
	//fmt.Println(slots)

}

func WriteJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	wrapper := make(map[string]interface{})

	wrapper[wrap] = data

	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func ErrorJSON(w http.ResponseWriter, err error) {
	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError{
		Message: err.Error(),
	}

	WriteJSON(w, http.StatusBadRequest, theError, "error")
}
