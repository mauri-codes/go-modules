package utils

import (
	"math/rand"
	"time"
)

const (
	LOWER_CASE_LETTERS = "abcdefghijklmnopqrstuvwxyz"
	UPPER_CASE_LETTERS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NUMBERS            = "0123456789"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset)-1)]
	}
	return string(b)
}

func RandomUpperCaseString(length int) string {
	return StringWithCharset(length, UPPER_CASE_LETTERS)
}

func RandomNumberString(length int) string {
	return StringWithCharset(length, NUMBERS)
}
