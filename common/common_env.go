package common

import "os"

func GetPort(varName string) string {
	port := os.Getenv(varName)
	if port == "" {
		port = ":9883"
	}
	return port
}
