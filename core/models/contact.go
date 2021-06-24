package models

import (
	"github.com/calvine/goauth/core/nullable"
)

// Contact is a model that represents a contact method for a user like phone or email.
type Contact struct {
	ID               string                  `bson:"id"`
	UserID           string                  `bson:"userId"`
	Name             nullable.NullableString `bson:"name"`
	Principal        string                  `bson:"principal"`
	Type             string                  `bson:"type"`
	IsPrimary        bool                    `bson:"isPrimary"`
	ConfirmationCode nullable.NullableString `bson:"confirmCode"`
	ConfirmedDate    nullable.NullableTime   `bson:"confirmedDate"`
	auditable
}
