package mongo

import (
	"context"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repoModels "github.com/calvine/goauth/dataaccess/mongo/internal/models"
	"github.com/calvine/richerror/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ProjUserOnly = bson.M{
		"_id":                            1,
		"passwordHash":                   1,
		"consecutiveFailedLoginAttempts": 1,
		"lockedOutUntil":                 1,
		"lastLoginDate":                  1,
	}
	ProjUserWithSpecificContact = bson.M{
		"_id":                            1,
		"passwordHash":                   1,
		"consecutiveFailedLoginAttempts": 1,
		"lockedOutUntil":                 1,
		"lastLoginDate":                  1,
		"contacts.$":                     1,
	}
)

// TODO: need to update these to use new rich errors

// userRepo is the repository struct for the user side of mongo db access. since other models related to users are embedded it makes sense (at least right now) to use a single struct for the related repository interfaces.
type userRepo struct {
	mongoClient    *mongo.Client
	dbName         string
	collectionName string
}

func NewUserRepo(client *mongo.Client) userRepo {
	return userRepo{client, DB_NAME, USER_COLLECTION}
}

func NewUserRepoWithNames(client *mongo.Client, dbName, collectionName string) userRepo {
	return userRepo{client, dbName, collectionName}
}

func (ur userRepo) GetUserByID(ctx context.Context, id string) (models.User, errors.RichError) {
	var repoUser repoModels.RepoUser
	options := options.FindOneOptions{
		Projection: ProjUserOnly,
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return repoUser.ToCoreUser(), coreerrors.NewFailedToParseObjectIDError(id, err, true)
	}
	filter := bson.M{"_id": oid}
	err = ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&repoUser)
	user := repoUser.ToCoreUser()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fields := map[string]interface{}{
				"_id": id,
			}
			return user, coreerrors.NewNoUserFoundError(fields, true)
		}
		return user, coreerrors.NewRepoQueryFailedError(err, true)
	}
	return user, nil
}

func (ur userRepo) GetUserAndContactByPrimaryContact(ctx context.Context, contactType, contactPrincipal string) (models.User, models.Contact, errors.RichError) {
	var receiver struct {
		User    repoModels.RepoUser      `bson:",inline"`
		Contact []repoModels.RepoContact `bson:"contacts"`
	}
	var user models.User
	var contact models.Contact

	options := options.FindOneOptions{
		Projection: ProjUserWithSpecificContact,
	}
	filter := bson.M{
		"contacts": bson.D{
			{
				Key: "$elemMatch", Value: bson.D{
					{Key: "isPrimary", Value: true},
					{Key: "type", Value: contactType},
					{Key: "principal", Value: contactPrincipal},
				},
			},
		},
	}
	err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&receiver)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fields := map[string]interface{}{
				"contacts.isPrimary": true,
				"contacts.type":      contactType,
				"contacts.principal": contactPrincipal,
			}
			return user, contact, coreerrors.NewNoUserFoundError(fields, true)
		}
		return user, contact, coreerrors.NewRepoQueryFailedError(err, true)
	}
	user = receiver.User.ToCoreUser()
	contact = receiver.Contact[0].ToCoreContact()
	contact.UserID = user.ID
	return user, contact, nil
}

func (ur userRepo) GetUserByPrimaryContact(ctx context.Context, contactType, contactPrincipal string) (models.User, errors.RichError) {
	var repoUser repoModels.RepoUser
	options := options.FindOneOptions{
		Projection: ProjUserOnly,
	}
	filter := bson.M{
		"contacts": bson.D{
			{
				Key: "$elemMatch", Value: bson.D{
					{Key: "isPrimary", Value: true},
					{Key: "type", Value: contactType},
					{Key: "principal", Value: contactPrincipal},
				},
			},
		},
	}
	err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).FindOne(ctx, filter, &options).Decode(&repoUser)
	user := repoUser.ToCoreUser()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fields := map[string]interface{}{
				"contacts.isPrimary": true,
				"contacts.type":      contactType,
				"contacts.principal": contactPrincipal,
			}
			return user, coreerrors.NewNoUserFoundError(fields, true)
		}
		return user, coreerrors.NewRepoQueryFailedError(err, true)
	}
	return user, nil
}

func (ur userRepo) AddUser(ctx context.Context, user *models.User, createdByID string) errors.RichError {
	user.AuditData.CreatedByID = createdByID
	user.AuditData.CreatedOnDate = time.Now().UTC()
	result, err := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).InsertOne(ctx, user, nil)
	if err != nil {
		return coreerrors.NewRepoQueryFailedError(err, true)
	}
	oid := result.InsertedID.(primitive.ObjectID)
	// oid, ok := result.InsertedID.(primitive.ObjectID)
	// if !ok {
	// 	return mongoerrors.NewMongoFailedToParseObjectID(result.InsertedID, true)
	// }
	user.ID = oid.Hex()
	return nil
}

func (ur userRepo) UpdateUser(ctx context.Context, user *models.User, modifiedByID string) errors.RichError {
	user.AuditData.ModifiedByID.Set(modifiedByID)
	user.AuditData.ModifiedOnDate.Set(time.Now().UTC())
	repoUser, err := repoModels.CoreUser(*user).ToRepoUser()
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id": bson.M{
			"$eq": repoUser.ObjectID,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"passwordHash":                   repoUser.PasswordHash,
			"consecutiveFailedLoginAttempts": repoUser.ConsecutiveFailedLoginAttempts,
			"lockedOutUntil":                 repoUser.LockedOutUntil.GetPointerCopy(),
			"lastLoginDate":                  repoUser.LastLoginDate.GetPointerCopy(),
			"modifiedById":                   repoUser.AuditData.ModifiedByID.GetPointerCopy(),
			"modifiedOnDate":                 repoUser.AuditData.ModifiedOnDate.GetPointerCopy(),
		},
	}
	result, updateErr := ur.mongoClient.Database(ur.dbName).Collection(ur.collectionName).UpdateOne(ctx, filter, update)
	if updateErr != nil {
		return coreerrors.NewRepoQueryFailedError(updateErr, true)
	}
	if result.ModifiedCount == 0 {
		fields := map[string]interface{}{
			"_id": user.ID,
		}
		return coreerrors.NewNoUserFoundError(fields, true)
	}
	return nil
}
