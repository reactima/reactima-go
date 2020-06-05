package utils

// ACHTUNG ! /util & codes should be consistent across all boilerplates

import (
	"net/http"
)

// Dont remove dummy function for referrences
func statusOK() int {
	return http.StatusOK
}

// Custom status codes
const (
	StatusNoErrorDescription    = 1001
	StatusErrorFromDB           = 1002
	StatusMissingToken          = 1003
	StatusTokenGenerationFailed = 1004
	StatusMissingRequiedFields  = 1005
	StatusInvalidJSON           = 1006
	StatusDeletedFromDB         = 1007
	StatusErrorFromGithubAPI    = 1008
	StatusErrorLogic            = 1009
	StatusErrorBindValidate     = 1010
	StatusErrorDeleteIDs        = 1011
	StatusErrorUpdateIDs        = 1012
	StatusErrorUpload           = 1013
	StatusErrorUpdatePassword   = 1014
	StatusWrondField         = 1015

	// 0... Golang net/http standard
	// 1... database related
	// 2... API related request & response
	// 3... batch related
	// 4... logic related

)

var desc = map[int]string{
	StatusNoErrorDescription:    "No error description, recommend to add more error codes",
	StatusErrorLogic:            "Failed business logic",
	StatusMissingToken:          "JWT must be provided inside of Authorization Bearer header",
	StatusTokenGenerationFailed: "Failed attempt to generate JWT",
	StatusMissingRequiedFields:  "All fields must be provided",
	StatusInvalidJSON:           "Parsing error",
	StatusErrorFromDB:           "Error returned by database",
	StatusDeletedFromDB:         "Successfully removed from database",
	StatusErrorFromGithubAPI:    "Error returned by Github API",
	StatusErrorBindValidate:     "Failed BindValidate",
	StatusErrorDeleteIDs:        "IDs failed to be deleted ",
	StatusErrorUpdateIDs:        "IDs failed to be updated ",
	StatusErrorUpload:           "Upload error",
	StatusErrorUpdatePassword:   "Password failed to be updated",
	StatusWrondField:   "Wrond Field Requested",
}

// StatusTitle returns a text for the HTTP status code. It returns StatusNoErrorDescription if the code is unknown.
func StatusTitle(code int) string {
	s, ok := desc[code]
	if ok {
		return s
	}
	return desc[StatusNoErrorDescription]
}
