package setup

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// EnvironmentConfig loads the .env file for the whole program.
// Use os.Getenv("LABEL_NAME") to access.
func EnvironmentConfig() {
	// get .env filepath
	absPath, err := filepath.Abs("./configs/.env")
	if err != nil {
		LogCommon(err).Fatal("Config filepath")
	}

	// get .env variables
	err = godotenv.Load(absPath)
	if err != nil {
		LogCommon(err).Fatal("Loading .env file")
	}
}

// EnvToInt converts the given .env variable into an integer.
func EnvToInt(s string) int {
	// get apps port
	port, err := strconv.Atoi(os.Getenv(s))
	if err != nil {
		LogCommon(err).
			WithField("port", port).
			Fatal("Failed env int conversion")
	}

	return port
}
