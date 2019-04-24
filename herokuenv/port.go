package herokuenv

import "os"

func init() {
	if portEnvVar := os.Getenv("PORT"); portEnvVar != "" {
		Port = portEnvVar
	} else {
		Port = "8080"
	}
}

var Port string
