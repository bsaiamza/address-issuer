package utils

import (
	"fmt"
)

// func IDValidator(id, gender, countryOfBirth string) (string, error) {
func IDValidator(id string) (string, error) {
	err := validChecksum(id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func validChecksum(id string) error {
	valid := validateChecksum(id)

	if !valid {
		return fmt.Errorf("invalid ID number")
	}

	return nil
}

func validateChecksum(pan string) bool {
	/* Validate string with Luhn (mod-10) */
	var alter bool
	var checksum int

	for position := len(pan) - 1; position > -1; position-- {
		digit := int(pan[position] - 48)
		if alter {
			digit = digit * 2
			if digit > 9 {
				digit = (digit % 10) + 1
			}
		}
		alter = !alter
		checksum += digit
	}
	return checksum%10 == 0
}
