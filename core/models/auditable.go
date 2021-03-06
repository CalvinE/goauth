package models

import (
	"time"

	"github.com/calvine/goauth/core/nullable"
)

type auditable struct {
	CreatedByID    string                  `bson:"createdById"`
	CreatedOnDate  time.Time               `bson:"createdOnDate"`
	ModifiedByID   nullable.NullableString `bson:"modifiedById"`
	ModifiedOnDate nullable.NullableTime   `bson:"modifiedOnDate"`
}
