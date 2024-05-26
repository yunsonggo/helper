package read

import "github.com/joho/godotenv"

func InEnv(files ...string) (envMap map[string]string, err error) {
	return godotenv.Read(files...)
}
