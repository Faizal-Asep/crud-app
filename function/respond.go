package function

import (
	"encoding/json"
	"net/http"
)

func JSONRespond(w http.ResponseWriter, code int, data interface{}) {
	response, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func ErrorRespond(w http.ResponseWriter, code int, message string) {
	JSONRespond(w, code, map[string]string{"error": message})
}

func ByteRespond(w http.ResponseWriter, code int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
