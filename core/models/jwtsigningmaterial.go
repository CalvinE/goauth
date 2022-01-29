package models

import (
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/jwt"
	"github.com/calvine/goauth/core/nullable"
	"github.com/calvine/richerror/errors"
	"github.com/google/uuid"
)

type JWTSigningMaterial struct {
	ID            string                        `bson:"-"`
	KeyID         string                        `bson:"keyId"`
	AlgorithmType jwt.JWTSingingAlgorithmFamily `bson:"algorithmType"`
	HMACSecret    nullable.NullableString       `bson:"hmacSecret"`
	Expiration    nullable.NullableTime         `bson:"expiration"`
	Disabled      bool                          `bson:"disabled"`
	AuditData     auditable                     `bson:",inline"`
	// PublicKey nullable.NullableString `bson:"publicKey"`
	// PrivateKey nullable.NullableString `bson:"privateKey"`
}

func NewHMACJWTSigningMaterial(secret string, expiration nullable.NullableTime) JWTSigningMaterial {
	return JWTSigningMaterial{
		KeyID:         uuid.Must(uuid.NewRandom()).String(), // TODO: make a function to create random unique key ids?
		AlgorithmType: jwt.HMAC,
		HMACSecret: nullable.NullableString{
			HasValue: true,
			Value:    secret,
		},
		Expiration: expiration,
		Disabled:   false,
	}
}

func (jsm *JWTSigningMaterial) IsExpired() bool {
	now := time.Now().UTC()
	if jsm.Expiration.HasValue && jsm.Expiration.Value.Before(now) {
		return true
	}
	return false
}

func (jsm *JWTSigningMaterial) ToSigner() (jwt.Signer, errors.RichError) {
	switch jsm.AlgorithmType {
	case jwt.HMAC:
		hmacOptions, err := jwt.NewHMACSigningOptions(jsm.HMACSecret.Value)
		if err != nil {
			return nil, err
		}
		return hmacOptions, nil
	default:
		err := coreerrors.NewJWTSigningMaterialAlgorithmTypeNotSupportedError(string(jsm.AlgorithmType), true)
		return nil, err
	}
}
