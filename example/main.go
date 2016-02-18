package main

import (
	"bytes"
	"crypto/subtle"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/gigablah/dashing-go.v1"
	_ "gopkg.in/gigablah/dashing-go.v1/example/jobs"
)

func tokenAuthMiddleware(h http.Handler) http.Handler {
	auth := []byte(os.Getenv("TOKEN"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if len(auth) == 0 {
      h.ServeHTTP(w, r)
      return
    }
		if r.Method == "POST" {
			body, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			r.Body = ioutil.NopCloser(bytes.NewReader(body))

			var data map[string]interface{}
			json.Unmarshal(body, &data)
			token, ok := data["auth_token"]
			if !ok {
				log.Printf("Auth token missing")
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			if result := subtle.ConstantTimeCompare(auth, []byte(token.(string))); result != 1 {
				log.Printf("Invalid auth token: expected %s, got %s", auth, token)
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dash := dashing.NewDashing().Start()
	log.Fatal(http.ListenAndServe(":"+port, tokenAuthMiddleware(dash)))
}
