package converter

import (
	"fmt"
	"strconv"
	"strings"
)

// ConvertStringNumberToNewBase converts the number provided as a string from the provided base to a new base
func ConvertStringNumberToNewBase(originalNumber string, originalBase string, newBase string) string {
	originalBaseInt := stringToInt(originalBase, 10)
	originalNumberInt64 := stringToInt64(originalNumber, originalBaseInt)
	newBaseInt := stringToInt(newBase, 10)
	finalValue := strconv.FormatInt(originalNumberInt64, newBaseInt)
	return finalValue
}

// StringToInt converts the given string to an int
func StringToInt(number string, base int) int {
	return stringToInt(number, base)
}

func stringToInt(number string, base int) int {
	value := stringToInt64(number, base)
	finalValue := int(value)
	return finalValue
}

// StringToInt64 converts the given string to an int64
func StringToInt64(number string, base int) int64 {
	return stringToInt64(number, base)
}

func stringToInt64(number string, base int) int64 {
	value, err := strconv.ParseInt(strings.TrimSpace(number), base, 64)
	if err != nil {
		fmt.Println(err)
	}
	return value
}
