package twelvefactor_ping

import (
	"os"
	"log"
	"fmt"
	"time"
	"context"
	"strconv"
	"net/http"
	"github.com/b3ntly/twelvefactor_ping/ping"
)

var (
	// make important variables noticeable immediately
	BACKGROUND_CTX = context.Background()
	DEFAULT_LOGGER = log.New(os.Stdout, "", log.Lshortfile)

	// read configuration from our environment with default values
	// ...expose as much configuration as possible to your users/ops team
	PORT = getEnv("PORT", "9090")
	DEFAULT_RESPONSE = getEnv("DEFAULT_RESPONSE", "PONG")
	DEFAULT_TIMEOUT = time.Duration(getEnvInt("DEFAULT_TIMEOUT", 500)) * time.Millisecond
)

func getEnv(key string, _default string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return _default
}

// notice that we fail hard with log.Fatal if any part of configuration goes wrong
func getEnvInt(key string, _default int) int {
	if val := os.Getenv(key); val != "" {
		ival, err := strconv.Atoi(val)

		if err != nil {
			log.Fatal(err)
		}

		return ival
	}
	return _default
}

// Using middleware to compose route definitions can be extremely powerful. Here we define a middleware
// that injects a context with a timeout. If that timeout is exceeded the request returns a error code
// indicating timeout.
func injectContextWithTimeout(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxWithTimeout, _ := context.WithTimeout(ctx, DEFAULT_TIMEOUT)
		next.ServeHTTP(w, r.WithContext(ctxWithTimeout))
	})
}

// break individual initialization steps into small sequential pieces
func buildServer(ctx context.Context, handler func(w http.ResponseWriter, r *http.Request)) *http.Server {
	return &http.Server{
		// Defaults to no timeouts, which is really bad
		// See: https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		// And: https://blog.cloudflare.com/exposing-go-on-the-internet/
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 2 * time.Second,

		Addr:    fmt.Sprintf(":%s", PORT),
		// Wrap the http handler in a middleware
		Handler: injectContextWithTimeout(ctx, http.HandlerFunc(handler)),
	}
}

// fail hard on any initialization errors
func listen(ctx context.Context, s *http.Server){
	log.Fatal(s.ListenAndServe())
}

func main(){
	service := ping.New(DEFAULT_RESPONSE, DEFAULT_LOGGER)
	server := buildServer(BACKGROUND_CTX, service.Endpoint)
	listen(BACKGROUND_CTX, server)
}