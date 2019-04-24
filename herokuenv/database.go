package herokuenv

import "os"

func init() {
	DatabaseURI = os.Getenv("DATABASE_URL")
}

var DatabaseURI string
