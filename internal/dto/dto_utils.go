package dto

import (
	"log"
	"time"
	"unicode"
)

const dateTimeFormat = time.RFC3339

func GetCurTime() string {
	return time.Now().Format(dateTimeFormat)
}

func IsDateTimeBefore(before string, after string) bool {
	beforeTime, _ := time.Parse(dateTimeFormat, before)
	afterTime, _ := time.Parse(dateTimeFormat, after)
	return beforeTime.Before(afterTime)

}

func IsDateTimeValid(dateTime string) bool {
	_, err := time.Parse(dateTimeFormat, dateTime)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func IsOppCodeValid(opCode string) bool {
	if len(opCode) != 3 {
		return false
	}
	for _, digit := range opCode {
		if !unicode.IsDigit(digit) {
			return false
		}
	}
	return true
}

func IsPhoneNumberValid(phoneNumber string) bool {
	if len(phoneNumber) != 11 || phoneNumber[0] != '7' {
		return false
	}
	for _, digit := range phoneNumber {
		if !unicode.IsDigit(digit) {
			return false
		}
	}
	return true
}

func IsUtcValid(utc string) bool {
	if len(utc) != 6 || (utc[0] != '+' && utc[0] != '-') || utc[3] != ':' {
		return false
	}
	if !unicode.IsDigit(rune(utc[1])) || !unicode.IsDigit(rune(utc[2])) || !unicode.IsDigit(rune(utc[4])) || !unicode.IsDigit(rune(utc[5])) {
		return false
	}
	if utc[1] > '2' || (utc[1] == '2' && utc[2] > '3') {
		return false
	}
	return true
}
