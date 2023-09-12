package route

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func parseBody(w http.ResponseWriter, r *http.Request, bodyHandler interface{}) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil && err != io.EOF {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	err = json.Unmarshal(body, bodyHandler)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return true
}
