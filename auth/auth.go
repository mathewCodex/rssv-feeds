package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetApiKey extracts an api key from the header of an http request
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization") //using http std lib to get Auth value
	if val == ""{
		return "", errors.New("no authentication info found")
	}

	vals := strings.Split(val, " ")//splitting the value with a white space
	if (len(vals) != 2){
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "ApiKey"{
		return "", errors.New("malformed first part of auth header")
	}
	return vals[1], nil //else return the api key 
}