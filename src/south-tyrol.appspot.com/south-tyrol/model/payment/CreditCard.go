package payment

import "time"

type CreditCard struct {
	CardNumber     string `json:"card_number"`
	ExpirationDate time.Time `json:"expiration_date"`
	SecurityCode   string `json:"security_code"`
}
