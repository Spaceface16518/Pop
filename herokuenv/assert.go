package herokuenv

// DatabaseURIExists checks if the DATABASE_URL variable exists.
func DatabaseURIExists() bool {
	return DatabaseURI != ""
}
