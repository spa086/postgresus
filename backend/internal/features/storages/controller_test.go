package storages

import (
	"net/http"
	local_storage "postgresus-backend/internal/features/storages/models/local"
	"postgresus-backend/internal/features/users"
	test_utils "postgresus-backend/internal/util/testing"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_SaveNewStorage_StorageReturnedViaGet(t *testing.T) {
	user := users.GetTestUser()
	router := createRouter()
	storage := createTestStorage(user.UserID)

	var savedStorage Storage
	test_utils.MakePostRequestAndUnmarshal(
		t, router, "/api/v1/storages", user.Token, storage, http.StatusOK, &savedStorage,
	)

	verifyStorageData(t, storage, &savedStorage)
	assert.NotEmpty(t, savedStorage.ID)

	// Verify storage is returned via GET
	var retrievedStorage Storage
	test_utils.MakeGetRequestAndUnmarshal(
		t,
		router,
		"/api/v1/storages/"+savedStorage.ID.String(),
		user.Token,
		http.StatusOK,
		&retrievedStorage,
	)

	verifyStorageData(t, &savedStorage, &retrievedStorage)

	// Verify storage is returned via GET all storages
	var storages []Storage
	test_utils.MakeGetRequestAndUnmarshal(
		t, router, "/api/v1/storages", user.Token, http.StatusOK, &storages,
	)

	assert.Contains(t, storages, savedStorage)
}

func Test_UpdateExistingStorage_UpdatedStorageReturnedViaGet(t *testing.T) {
	user := users.GetTestUser()
	router := createRouter()
	storage := createTestStorage(user.UserID)

	// Save initial storage
	var savedStorage Storage
	test_utils.MakePostRequestAndUnmarshal(
		t, router, "/api/v1/storages", user.Token, storage, http.StatusOK, &savedStorage,
	)

	// Modify storage name
	updatedName := "Updated Storage " + uuid.New().String()
	savedStorage.Name = updatedName

	// Update storage
	var updatedStorage Storage
	test_utils.MakePostRequestAndUnmarshal(
		t, router, "/api/v1/storages", user.Token, savedStorage, http.StatusOK, &updatedStorage,
	)

	// Verify updated data
	assert.Equal(t, updatedName, updatedStorage.Name)
	assert.Equal(t, savedStorage.ID, updatedStorage.ID)

	// Verify through GET
	var retrievedStorage Storage
	test_utils.MakeGetRequestAndUnmarshal(
		t,
		router,
		"/api/v1/storages/"+updatedStorage.ID.String(),
		user.Token,
		http.StatusOK,
		&retrievedStorage,
	)

	verifyStorageData(t, &updatedStorage, &retrievedStorage)

	// Verify storage is returned via GET all storages
	var storages []Storage
	test_utils.MakeGetRequestAndUnmarshal(
		t, router, "/api/v1/storages", user.Token, http.StatusOK, &storages,
	)

	assert.Contains(t, storages, updatedStorage)
}

func Test_DeleteStorage_StorageNotReturnedViaGet(t *testing.T) {
	user := users.GetTestUser()
	router := createRouter()
	storage := createTestStorage(user.UserID)

	// Save initial storage
	var savedStorage Storage
	test_utils.MakePostRequestAndUnmarshal(
		t, router, "/api/v1/storages", user.Token, storage, http.StatusOK, &savedStorage,
	)

	// Delete storage
	test_utils.MakeDeleteRequest(
		t, router, "/api/v1/storages/"+savedStorage.ID.String(), user.Token, http.StatusOK,
	)

	// Try to get deleted storage, should return error
	response := test_utils.MakeGetRequest(
		t, router, "/api/v1/storages/"+savedStorage.ID.String(), user.Token, http.StatusBadRequest,
	)

	assert.Contains(t, string(response.Body), "error")

	// Verify storage is not returned via GET all storages
	var storages []Storage
	test_utils.MakeGetRequestAndUnmarshal(
		t, router, "/api/v1/storages", user.Token, http.StatusOK, &storages,
	)

	assert.NotContains(t, storages, savedStorage)
}

func Test_TestDirectStorageConnection_ConnectionEstablished(t *testing.T) {
	user := users.GetTestUser()
	router := createRouter()
	storage := createTestStorage(user.UserID)

	response := test_utils.MakePostRequest(
		t, router, "/api/v1/storages/direct-test", user.Token, storage, http.StatusOK,
	)

	assert.Contains(t, string(response.Body), "successful")
}

func Test_TestExistingStorageConnection_ConnectionEstablished(t *testing.T) {
	user := users.GetTestUser()
	router := createRouter()
	storage := createTestStorage(user.UserID)

	var savedStorage Storage
	test_utils.MakePostRequestAndUnmarshal(
		t, router, "/api/v1/storages", user.Token, storage, http.StatusOK, &savedStorage,
	)

	// Test connection to existing storage
	response := test_utils.MakePostRequest(
		t,
		router,
		"/api/v1/storages/"+savedStorage.ID.String()+"/test",
		user.Token,
		nil,
		http.StatusOK,
	)

	assert.Contains(t, string(response.Body), "successful")
}

func Test_CallAllMethodsWithoutAuth_UnauthorizedErrorReturned(t *testing.T) {
	router := createRouter()
	storage := createTestStorage(uuid.New())

	// Test endpoints without auth
	endpoints := []struct {
		method string
		url    string
		body   interface{}
	}{
		{"GET", "/api/v1/storages", nil},
		{"GET", "/api/v1/storages/" + uuid.New().String(), nil},
		{"POST", "/api/v1/storages", storage},
		{"DELETE", "/api/v1/storages/" + uuid.New().String(), nil},
		{"POST", "/api/v1/storages/" + uuid.New().String() + "/test", nil},
		{"POST", "/api/v1/storages/direct-test", storage},
	}

	for _, endpoint := range endpoints {
		testUnauthorizedEndpoint(t, router, endpoint.method, endpoint.url, endpoint.body)
	}
}

func testUnauthorizedEndpoint(
	t *testing.T,
	router *gin.Engine,
	method, url string,
	body interface{},
) {
	test_utils.MakeRequest(t, router, test_utils.RequestOptions{
		Method:         method,
		URL:            url,
		Body:           body,
		ExpectedStatus: http.StatusUnauthorized,
	})
}

func createRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	controller := GetStorageController()
	v1 := router.Group("/api/v1")
	controller.RegisterRoutes(v1)
	return router
}

func createTestStorage(userID uuid.UUID) *Storage {
	return &Storage{
		UserID:       userID,
		Type:         StorageTypeLocal,
		Name:         "Test Storage " + uuid.New().String(),
		LocalStorage: &local_storage.LocalStorage{},
	}
}

func verifyStorageData(t *testing.T, expected *Storage, actual *Storage) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.UserID, actual.UserID)
}
