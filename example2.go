package main

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
)

var (
	ErrValueTooLong = errors.New("cookie value too long ")
	ErrInvalidValue = errors.New("invalid cookie value")
)

// Write() function encodes a cookie value to base64 and ensures
// the length of the cookie is no more than 4096 bytes before writing it
func write(w http.ResponseWriter, cookie http.Cookie) error {
	// Encode the cookie value using base64
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	//check the total length of the cookie contents. Return the ErrValueTooLong
	// error if it's more than b4096 bytes.
	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}

	// write the cookie as normal
	http.SetCookie(w, &cookie)
	return nil
}

// Read() function which read a cookie from the current request and
//
//	decodes the cookie value from base64
func Read(r *http.Request, name string) (string, error) {
	// Read the cookie as normal
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	// Decode the base64-encoded cookie value. If the cookie didn't contain a
	// valid base64 value, this operation will fail and we return an
	// ErrInvalidValue error.
	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}

	// Return the decoded cookie value
	return string(value), nil

}

func main() {
	// start a web server with the two endpoints
	mux := http.NewServeMux()
	mux.HandleFunc("/set", setCookieHandler)
	mux.HandleFunc("/get", getCookieHandler)

	log.Print("Listening...")
	err := http.ListenAndServe(":5000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

// Modified cookie handler:
func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize a cookie as normal
	cookie := http.Cookie{
		Name:     "exampleCookie",
		Value:    "Hello Zoë!!",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	// Write the cookie, If ther is an erro due to and encoding failure or
	// it being too long) then log the error and send a 500 Internal server error
	// response

	err := write(w, cookie) // Use our helper function.

	if err != nil {
		log.Println(err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	// Write a HTTP resonse as normal
	w.Write([]byte("cookie set!"))
}

// Modified getCookieHandler
func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	// Use the Read() function to retrieve the cookie value, additionally
	// checking for the ErrInvalidValue error and handling it as necessary.

	value, err := Read(r, "exampleCookie")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found ", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}
	// Echo out the cookie value in the response body

	w.Write([]byte(value))

}