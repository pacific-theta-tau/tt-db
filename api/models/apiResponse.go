// apiResponse.go: Model for API handler responses
package models

import (
	"encoding/json"
	"log"
	"net/http"
)


// API response format in JSend notation. See: https://github.com/omniti-labs/jsend
// @Description JSON response format for all API calls
type APIResponse struct {
    Status      string `json:"status"`
    Message     string `json:"message,omitnil"` //Omit if nil
    Data        interface{} `json:"data,omitnil"` //Omit if nil
}

// helper function to set headers and encode json
func sendResponse(w http.ResponseWriter, statusCode int, response APIResponse) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Error encoding JSON: %v", err)
    }
}

// 200~ status codes
func RespondWithSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
    response := APIResponse{
        Status: "success",
        Data: data,
    }
    sendResponse(w, statusCode, response)
}

// 500~ status codes (Server error)
func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
    response := APIResponse{
        Status:  "error",
        Message: message,
    }
    sendResponse(w, statusCode, response)
}

// 400~ status codes (Client error)
func RespondWithFail(w http.ResponseWriter, statusCode int, message string) {
    response := APIResponse{
        Status: "fail",
        Message: message,
    }
    sendResponse(w, statusCode, response)
}
