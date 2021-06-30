package models

import (
	"github.com/calvine/goauth/core/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CoreAddress models.Address

type RepoAddress struct {
	ObjectId    primitive.ObjectID `bson:"id"`
	CoreAddress `bson:",inline"`
}

func (ra RepoAddress) ToCoreAddress() models.Address {
	oidString := ra.ObjectId.Hex()
	ra.CoreAddress.ID = oidString

	return models.Address(ra.CoreAddress)
}

func (ca CoreAddress) ToRepoAddress() (RepoAddress, error) {
	oid, err := primitive.ObjectIDFromHex(ca.ID)
	if err != nil {
		return RepoAddress{}, err
	}
	return RepoAddress{
		ObjectId:    oid,
		CoreAddress: ca,
	}, nil
}

func (ca CoreAddress) ToRepoAddressWithoutId() RepoAddress {
	return RepoAddress{
		CoreAddress: ca,
	}
}
