package helpers

import (
	"fmt"
	"os"
	"strconv"
)

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}

	return url
}

func InitMachineID() (uint16, error) {
	idStr := os.Getenv("MACHINE_ID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid MACHINE_ID: %v", err)
	}
	if id < 0 || id > 65535 {
		return 0, fmt.Errorf("MACHINE_ID must be between 0 and 65535")
	}
	return uint16(id), nil
}
