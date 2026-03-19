//go:build ignore

package main

var userProfile = map[string]interface{}{
	"firstName":      "John",
	"lastName":       "Doe",
	"birthDate":      "1990-05-15",
	"emailAddress":   "john_doe@example.com",
	"phoneNumber":    "555-123-4567",
	"mailingAddress": map[string]interface{}{
		"streetName":     "Main Street",
		"houseNumber":    123,
		"apartmentUnit":  "4B",
		"zipCode":        "10001",
		"cityName":       "New York",
	},
}
