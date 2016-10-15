package network

import (
    "net/http"
    "issue"
    "errors"
)

func Authorization(w http.ResponseWriter, r *http.Request) string {
    header := r.Header.Get("Authorization")

    if header == "" {
        issue.Handle(w, errors.New("unauthorized"), http.StatusUnauthorized)
    }

    return header
}