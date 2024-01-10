package utils

import (
	"fmt"
	"time"
)

func ParseDateInLocation(dateString, layout string) (time.Time, error) {
	parsedDate, err := time.ParseInLocation(layout, dateString, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return parsedDate, nil
}

func ParseBool(str string) (bool, error) {
	switch str {
	case "y", "Y":
		return true, nil
	case "n", "N":
		return false, nil
	}
	return false, fmt.Errorf("must be y/n or Y/N")
}
