package repotest

import (
	"context"
	"testing"
	"time"

	"github.com/calvine/goauth/core"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
)

var (
	testUser1 = models.User{
		PasswordHash: "passwordhash1",
	}
	nonExistantUserID string
)

func setupUserRepoTestData(_ *testing.T, testingHarness RepoTestHarnessInput) {
	nonExistantUserID = testingHarness.IDGenerator(false)
}

func testUserRepo(t *testing.T, testHarness RepoTestHarnessInput) {
	setupUserRepoTestData(t, testHarness)
	// functionality tests
	t.Run("AddUser", func(t *testing.T) {
		_testAddUser(t, *testHarness.UserRepo)
	})
	t.Run("UpdateUser", func(t *testing.T) {
		_testUpdateUser(t, *testHarness.UserRepo)
	})
	t.Run("GetUserByID", func(t *testing.T) {
		_testGetUserByID(t, *testHarness.UserRepo)
	})
	t.Run("GetUserByPrimaryContact", func(t *testing.T) {
		_testGetUserByPrimaryContact(t, *testHarness.UserRepo)
	})
	t.Run("GetUserAndContactByContact", func(t *testing.T) {
		_testGetUserAndContactByContact(t, *testHarness.UserRepo)
	})
}

func _testAddUser(t *testing.T, userRepo repo.UserRepo) {
	createdByID := "user repos tests"

	err := userRepo.AddUser(context.TODO(), &testUser1, createdByID)
	if err != nil {
		t.Log(err.Error())
		t.Error("failed to add user to database", err.GetErrorCode())
	}
	if testUser1.AuditData.CreatedByID != createdByID {
		t.Error("failed to set the test user 1 CreatedByID to the right value", testUser1.AuditData.CreatedByID, createdByID)
	}
}

func _testUpdateUser(t *testing.T, userRepo repo.UserRepo) {
	preUpdateDate := time.Now().UTC()
	newPasswordHash := "another secure password hash"
	testUser1.PasswordHash = newPasswordHash
	err := userRepo.UpdateUser(context.TODO(), &testUser1, testUser1.ID)
	if err != nil {
		t.Log(err.Error())
		t.Error("failed to update user", err.GetErrorCode())
	}
	if !testUser1.AuditData.ModifiedOnDate.Value.After(preUpdateDate) {
		t.Error("ModifiedOnDate should be after the preUpdate for test", preUpdateDate, testUser1.AuditData.ModifiedOnDate)
	}
	retreivedUser, err := userRepo.GetUserByID(context.TODO(), testUser1.ID)
	if err != nil {
		t.Log(err.Error())
		t.Error("failed to retreive updated user to check that fields were updated.", err.GetErrorCode())
	}
	if retreivedUser.PasswordHash != newPasswordHash {
		t.Error("password hash did not update.", retreivedUser.PasswordHash, newPasswordHash)
	}
}

func _testGetUserByID(t *testing.T, userRepo repo.UserRepo) {
	userID := initialTestUser.ID
	retreivedUser, err := userRepo.GetUserByID(context.TODO(), userID)
	if err != nil {
		t.Log(err.Error())
		t.Error("error getting user with id", userID, err.GetErrorCode())
	}
	if retreivedUser.PasswordHash != initialTestUser.PasswordHash {
		t.Error("retreivedUser should have same data as user with id tested", retreivedUser, initialTestUser)
	}
}

func _testGetUserByPrimaryContact(t *testing.T, userRepo repo.UserRepo) {
	contactType, principal := core.CONTACT_TYPE_EMAIL, initialTestContact.Principal
	retreivedUser, err := userRepo.GetUserByPrimaryContact(context.TODO(), contactType, principal)
	if err != nil {
		t.Log(err.Error())
		t.Error("failed to retreive user via primary contact info", contactType, principal, err.GetErrorCode())
	}
	if retreivedUser.ID != initialTestUser.ID {
		t.Error("expected retreivedUser and initialTestUser ID to match", retreivedUser.ID, initialTestUser.ID)
	}
}

func _testGetUserAndContactByContact(t *testing.T, userRepo repo.UserRepo) {
	contactType, principal := core.CONTACT_TYPE_EMAIL, initialTestContact.Principal
	retreivedUser, retreivedContact, err := userRepo.GetUserAndContactByContact(context.TODO(), contactType, principal)
	if err != nil {
		t.Log(err.Error())
		t.Error("failed to retreive user via primary contact info", contactType, principal, err.GetErrorCode())
	}
	if retreivedUser.ID != initialTestUser.ID {
		t.Error("expected retreivedUser and initialTestUser ID to match", retreivedUser.ID, initialTestUser.ID)
	}
	if retreivedContact.Principal != principal {
		t.Error("expected retreivedContact.Principal and the test principal to be the same", retreivedContact.Principal, principal)
	}
	if retreivedContact.UserID != retreivedUser.ID {
		t.Error("expected user.ID and contact.userID to match", retreivedContact.UserID, retreivedUser.ID)
	}
}
