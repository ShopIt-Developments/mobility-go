package id

import (
    "time"
    "math/rand"
)

const alphanumericChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const alphanumericLength = 64

type Id struct {
    Id string `json:"id"`
}

func Alphanumeric() string {
    rand.Seed(time.Now().UTC().UnixNano())

    result := make([]byte, alphanumericLength)

    for i := 0; i < alphanumericLength; i++ {
        result[i] = alphanumericChars[rand.Intn(len(alphanumericChars))]
    }

    return string(result)
}