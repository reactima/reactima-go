package bindvalidate

const (
	// authorization
	missingToken          = "3001"
	tokenGenerationFailed = "3002"
	// validation
	missingRequiedFields = "4001"
	// data corruption
	invalidJSON = "4101"

	errorFromDB = "1001"

	deletedFromDB = "1204"
)

var statusTitle = map[string]string{
	missingToken:          "JWT must be provided inside of Authorization Bearer header",
	tokenGenerationFailed: "Failed attempt to generate JWT",
	missingRequiedFields:  "All fields must be provided",
	invalidJSON:           "Parsing error",
	errorFromDB:           "Error returned by database",
	deletedFromDB:         "Successfully removed from database",
}

// StatusTitle returns a text for the HTTP status code. It returns the empty string if the code is unknown.
func StatusTitle(code string) string {
	return statusTitle[code]
}
