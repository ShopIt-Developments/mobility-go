package util

import (
    "time"
    "math/rand"
)

const alphanumericChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const alphanumericLength = 64

func Alphanumeric() string {
    rand.Seed(time.Now().UTC().UnixNano())

    result := make([]byte, alphanumericLength)

    for i := 0; i < alphanumericLength; i++ {
        result[i] = alphanumericChars[rand.Intn(len(alphanumericChars))]
    }

    return string(result)
}