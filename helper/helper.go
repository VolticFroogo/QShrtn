package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// GenerateRandomString generates a random string of a specified length.
func GenerateRandomString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// ThrowErr is the function used to throw an error when on an HTML request.
func ThrowErr(w http.ResponseWriter, r *http.Request, err error) {
	// Create the new URL to redirect the user to.
	// This will be the error page along with some debug information.
	// Example: forum.com/error/sql%3A+no+rows+in+result+set
	newURL := fmt.Sprintf("/error/%v", url.QueryEscape(err.Error()))

	// Redirect the user to the new URL.
	http.Redirect(w, r, newURL, http.StatusTemporaryRedirect)
	log.Printf("Error: %v", err)
}

// JSONResponse sends a client a JSON response.
func JSONResponse(data interface{}, w http.ResponseWriter) (err error) {
	dataJSON, err := json.Marshal(data) // Encode response into JSON.
	if err != nil {
		return
	}
	w.Write(dataJSON) // Write JSON data to response writer.
	return
}

// Seed seeds the random number generator.
func Seed() {
	rand.Seed(time.Now().UnixNano())
}
