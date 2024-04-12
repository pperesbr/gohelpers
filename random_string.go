package gohelpers

import (
	"math/rand"
	"time"
)

const (
	charsetLowerCase     = "abcdefghijklmnopqrstuvwxyz0123456789"
	charsetUpperCase     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetCaseSensitive = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

const (
	LowerCase = iota + 1
	UpperCase
)

func randomString(length int, charset string) string {

	b := make([]byte, length)

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func GenerateRandomString(length int, charsetType int) string {

	switch charsetType {
	case LowerCase:
		return randomString(length, charsetLowerCase)
	case UpperCase:
		return randomString(length, charsetUpperCase)
	default:
		return randomString(length, charsetCaseSensitive)
	}
}
