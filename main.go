package twelvefactor_ping

import (
	"github.com/b3ntly/twelvefactor_ping/ping"
	"net/http"
	"os"
	"log"
	"fmt"
)

var (
	PORT = getenv("PORT", "9090")
	DEFAULT_RESPONSE = getenv("DEFAULT_RESPONSE", "PONG")
	DEFAULT_LOGGER = log.New(os.Stdout, "", log.Lshortfile)
)

func getenv(key string, _default string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return _default
}

func buildServer(handler func(w http.ResponseWriter, r *http.Request)) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: http.HandlerFunc(handler),
	}
}

func listen(s *http.Server){
	DEFAULT_LOGGER.Fatal(s.ListenAndServe())
}

func main(){
	service := ping.New(DEFAULT_RESPONSE, DEFAULT_LOGGER)
	server := buildServer(service.Endpoint)
	listen(server)
}