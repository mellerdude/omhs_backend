package tests

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"omhs-backend/models"

	"github.com/stretchr/testify/assert"
)

// TestCreateDocument tests the creation of a new document in the database.
func TestCreateDocument(t *testing.T) {
	router, _ := initializeRouterAndControllers(client)

	adminToken := AdminLogin(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))

	doc := models.Document{
		Data: map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		},
	}

	body, code := createDocument(router, "testdb", "testcollection", adminToken, doc)
	assert.Equal(t, http.StatusOK, code)

	var createdDoc models.Document
	json.Unmarshal([]byte(body), &createdDoc)
	assert.Equal(t, doc.Data["field1"], createdDoc.Data["field1"])
	assert.Equal(t, doc.Data["field2"], createdDoc.Data["field2"])

	// Cleanup
	_, code = deleteDocument(router, "testdb", "testcollection", createdDoc.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)

	requestsTestManager.RegisterTest(t, "TestCreateDocument")
}

// TestGetDocument tests retrieving a document by its ID from the database.
func TestGetDocument(t *testing.T) {
	router, _ := initializeRouterAndControllers(client)

	adminToken := AdminLogin(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))

	field1 := generateRandomString(10)
	field2 := generateRandomString(10)
	value1 := generateRandomString(10)
	value2 := generateRandomString(10)
	doc := models.Document{
		Data: map[string]interface{}{
			field1: value1,
			field2: value2,
		},
	}

	body, code := createDocument(router, "testdb", "testcollection", adminToken, doc)
	assert.Equal(t, http.StatusOK, code)

	var createdDoc models.Document
	json.Unmarshal([]byte(body), &createdDoc)

	body, code = getDocument(router, "testdb", "testcollection", createdDoc.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)

	var retrievedDoc models.Document
	json.Unmarshal([]byte(body), &retrievedDoc)
	assert.Equal(t, createdDoc.Data[field1], retrievedDoc.Data[field1])
	assert.Equal(t, createdDoc.Data[field2], retrievedDoc.Data[field2])

	// Cleanup
	_, code = deleteDocument(router, "testdb", "testcollection", createdDoc.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)

	requestsTestManager.RegisterTest(t, "TestGetDocument")
}

// TestUpdateDocument tests updating an existing document in the database.
func TestUpdateDocument(t *testing.T) {
	router, _ := initializeRouterAndControllers(client)

	adminToken := AdminLogin(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))

	doc := models.Document{
		Data: map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		},
	}

	body, code := createDocument(router, "testdb", "testcollection", adminToken, doc)
	assert.Equal(t, http.StatusOK, code)

	var createdDoc models.Document
	json.Unmarshal([]byte(body), &createdDoc)

	updatedDoc := models.Document{
		Data: map[string]interface{}{
			"field1": "new_value1",
			"field2": "new_value2",
		},
	}

	body, code = updateDocument(router, "testdb", "testcollection", createdDoc.ID.Hex(), adminToken, updatedDoc)
	assert.Equal(t, http.StatusOK, code)

	var updatedDocResponse models.Document
	json.Unmarshal([]byte(body), &updatedDocResponse)
	assert.Equal(t, updatedDoc.Data["field1"], updatedDocResponse.Data["field1"])
	assert.Equal(t, updatedDoc.Data["field2"], updatedDocResponse.Data["field2"])

	// Cleanup
	_, code = deleteDocument(router, "testdb", "testcollection", createdDoc.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)

	requestsTestManager.RegisterTest(t, "TestUpdateDocument")
}

// TestDeleteDocument tests deleting a document from the database.
func TestDeleteDocument(t *testing.T) {
	router, _ := initializeRouterAndControllers(client)

	adminToken := AdminLogin(router, os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASS"))

	doc := models.Document{
		Data: map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		},
	}

	body, code := createDocument(router, "testdb", "testcollection", adminToken, doc)
	assert.Equal(t, http.StatusOK, code)

	var createdDoc models.Document
	json.Unmarshal([]byte(body), &createdDoc)

	_, code = deleteDocument(router, "testdb", "testcollection", createdDoc.ID.Hex(), adminToken)
	assert.Equal(t, http.StatusOK, code)

	requestsTestManager.RegisterTest(t, "TestDeleteDocument")
}
