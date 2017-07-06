package ping

import (
	"context"
	"fmt"
	"github.com/b3ntly/twelvefactor_ping/ping"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	// BACKGROUND_CTX makex important variables noticeable immediately
	BACKGROUND_CTX = context.Background()
	// DEFAULT_LOGGER replaces this with something like Logrus or Zap in production
	DEFAULT_LOGGER = log.New(os.Stdout, "", log.Lshortfile)
	// DEFAULT_MUX, replace with gorilla etc. if desired
	DEFAULT_MUX = http.NewServeMux()

	// read configuration from our environment with default values
	// ...expose as much configuration as possible to your users/ops team
	// PORT to serve the application on
	PORT = getEnv("PORT", "9090")
	// the path from which to reply
	// PATH to return a response from
	ENDPOINT = getEnv("ENDPOINT", "/ping")
	// DEFAULT_RESPONSE to send
	DEFAULT_RESPONSE = getEnv("DEFAULT_RESPONSE", "PONG")
	// REQ_TIMEOUT timeout for http handler
	REQ_TIMEOUT = time.Duration(getEnvInt("REQ_TIMEOUT", 500)) * time.Millisecond
	// SERVER_READ_TIMEOUT timeout for the server, REQ_TIMEOUT should occur before a server timeout in most circumstances
	SERVER_READ_TIMEOUT = time.Duration(getEnvInt("SERVER_READ_TIMEOUT", 1000)) * time.Millisecond
	// SERVER_WRITE_TIMEOUT timeout for the server, REQ_TIMEOUT should occur before a server timeout
	SERVER_WRITE_TIMEOUT = time.Duration(getEnvInt("SERVER_WRITE_TIMEOUT", 2000)) * time.Millisecond
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
			DEFAULT_LOGGER.Fatal(err)
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
		ctxWithTimeout, cancel := context.WithTimeout(ctx, REQ_TIMEOUT)

		// defer can to prevent context leaking
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctxWithTimeout))
	})
}

// break individual initialization steps into small sequential pieces
func buildServer(ctx context.Context, mux http.Handler) *http.Server {
	return &http.Server{
		// Defaults to no timeouts, which is really bad
		// See: https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		// And: https://blog.cloudflare.com/exposing-go-on-the-internet/
		ReadTimeout:  SERVER_READ_TIMEOUT,
		WriteTimeout: SERVER_WRITE_TIMEOUT,
		Addr:         fmt.Sprintf(":%s", PORT),
		Handler:      mux,
	}
}

// fail hard on any initialization errors, include context per convention
func listen(ctx context.Context, s *http.Server) {
	DEFAULT_LOGGER.Fatal(s.ListenAndServe())
}

func main() {
	// instantiate the service
	service := ping.New(DEFAULT_RESPONSE, DEFAULT_LOGGER)

	// build the router with desired middleware
	DEFAULT_MUX.Handle(ENDPOINT, injectContextWithTimeout(BACKGROUND_CTX, http.HandlerFunc(service.Endpoint)))

	// instantiate the http.Server with our router
	server := buildServer(BACKGROUND_CTX, DEFAULT_MUX)

	// start the server
	listen(BACKGROUND_CTX, server)
}
