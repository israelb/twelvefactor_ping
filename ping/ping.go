package ping

import (
	"encoding/json"
	"log"
	"net/http"
)

type Service struct {
	PingResponse string
	Logger       *log.Logger
}

// inject dependencies via an explicit constructor. Though sometimes people will read environmental variables or
// initialize defaults here I prefer to do so explicitly within the program entry-point.
func New(pingResponse string, logger *log.Logger) *Service {
	return &Service{PingResponse: pingResponse, Logger: logger}
}

// Alternatively you can return a mux or sub-router from this subpackage
func (s *Service) Endpoint(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(s.PingResponse)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		return
	}

	s.Logger.Println(err)

	// here you can decide what type of error to return to the user, you should usually refrain from the actual error
	// in case it contains sensitive information. That said you should try to tell the user something helpful.
	http.Error(w, "", http.StatusInternalServerError)
}
