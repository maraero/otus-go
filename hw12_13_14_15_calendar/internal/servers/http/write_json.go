package serverhttp

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	resp, _ := json.Marshal(v)
	w.Write(resp)
}
