//package categories contains test for categories.go
package categories

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//categoryByIdTest is a structure for testing GET /category/{id} requests
var categoryByIdTest = []struct {
	categoryId   string // category id
	expected 	 string // expected result
	expectedCode int //expected http response status code
}{
	{
		"bq4fasj7jhfi127rimlg",
		`{"CategoryID":"bq4fasj7jhfi127rimlg","CategoryName":"Shopping Products","CategoryDescription":"Products consumers purchase and consume on a less frequent schedule compared to convenience products."}`,
		200,
	},
	{
		"randomID",
		`Category with ID randomID not found`,
		412,
	},
}

//TestGetAllCategories tests whether GetAllCategories func returns the right response body and status
func TestGetAllCategories(t *testing.T) {
	//Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/categories", nil)
	if err != nil {
		t.Fatal(err)
	}
	//Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllCategories)

	handler.ServeHTTP(rr, req)

	// Check the response status code is what we expect
	assert.Equal(t, 200, rr.Code, "OK response is expected")
	// Check the response body is what we expect
	expected := `[{"CategoryID":"bq4fasj7jhfi127rimlg","CategoryName":"Shopping Products","CategoryDescription":"Products consumers purchase and consume on a less frequent schedule compared to convenience products."},{"CategoryID":"bq4fb3b7jhfi7v7uo39g","CategoryName":"Specialty Products","CategoryDescription":"Products that are more expensive relative to convenience and shopping products."}]`
	assert.JSONEq(t, expected, rr.Body.String(), "Expected response body to be the same")
}

//TestGetCategoryById tests whether GetCategoryById func returns the right response bodies and statuses
//while iterating over categoryByIdTest
func TestGetCategoryById(t *testing.T) {
	for _, p := range categoryByIdTest {
		//request url
		path := "/categories/" + p.categoryId
		req, err := http.NewRequest("GET", path, nil)
		//passing the {id} var as it was not recognizable for some reason
		req = mux.SetURLVars(req, map[string]string{"id": p.categoryId})
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetCategoryById)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, p.expectedCode, rr.Code, "OK response is expected")
		assert.Equal(t, p.expected, strings.TrimSuffix(rr.Body.String(),"\n"), "Response body is expected to be equal to expected value")
	}
}

//TestCreateCategory tests whether CreateCategory func returns the right status and actually appends a Category to []Categories
func TestCreateCategory(t *testing.T) {
	//initial length of []Categories
	initialLen := len(Categories)
	//parameters passed to request body
	requestBody := &Category{
		CategoryName: 		"Super Cool Category",
		CategoryDescription: "Brand new cool Category",
	}
	jsonCategory, _ := json.Marshal(requestBody)
	//Create a request to pass to the handler with request body as a third parameter
	req, err := http.NewRequest("POST", "/categories/new", bytes.NewBuffer(jsonCategory))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateCategory)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 201, rr.Code, "Created response is expected")
	//the length of []Categories should increase after creating new category
	assert.NotEqual(t, initialLen, len(Categories), "Expected length to increase after creating new Category")
}

//TestCreateCategoryEmptyBody tests whether CreateCategory func returns the right status and does not append a Category to []Categories
//because of the empty request body
func TestCreateCategoryEmptyBody (t *testing.T) {
	//initial length of []Categories
	initialLen := len(Categories)
	//empty body
	requestBody := &Category{}
	jsonCategory, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/categories/new", bytes.NewBuffer(jsonCategory))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateCategory)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 422, rr.Code, "Unprocessable Entity response is expected")
	//the length of []Categories should not change after trying to create new empty category
	assert.Equal(t, initialLen, len(Categories), "Expected length to stay the same after adding empty category name")
}

//TestCreateCategoryWrongJSONSyntax tests whether CreateCategory func returns the right status and does not append a Category to []Categories
//because of the wrong syntax in JSON request body
func TestCreateCategoryWrongJSONSyntax(t *testing.T) {
	//initial length of []Categories
	initialLen := len(Categories)
	//parameters passed to request body
	requestBody := `{{"CategoryID":"bq4fasj7jhfi127rimlg","CategoryName":"Name",,,}}`
	req, err := http.NewRequest("POST", "/categories/new", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateCategory)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 400, rr.Code, "Bad request response is expected")
	assert.Equal(t, initialLen, len(Categories), "Expected length to stay the same after wrong syntax json")

}

//TestDeleteCategory tests whether DeleteCategory func returns the right status and actually deletes a Category from []Categories
func TestDeleteCategory(t *testing.T) {
	//initial length of []products
	initialLen := len(Categories)

	req, err := http.NewRequest("DELETE", "/categories/bq4fasj7jhfi127rimlg", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "bq4fasj7jhfi127rimlg"})

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteCategory)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code, "OK response is expected")
	//the length of []Categories should decrease after deleting category
	assert.NotEqual(t, initialLen, len(Categories), "Expected length to decrease after deleting new Category")
}

//TestDeleteCategoryWrongID tests whether DeleteCategory func returns the right status and does not delete a Category from []Categories
//because the ID is wrong
func TestDeleteCategoryWrongID(t *testing.T) {
	//initial length of []products
	initialLen := len(Categories)

	req, err := http.NewRequest("DELETE", "/categories/randomID", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "randomID"})

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteCategory)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 412, rr.Code, "Precondition Failed response is expected")
	//the length of []Categories should not change after trying to delete non-existing category
	assert.Equal(t, initialLen, len(Categories), "Expected length to stay the same after creating new Category")
}

//TestUpdateCategory tests whether UpdateCategory func returns the right status and does not change []Categories, but updates fields
func TestUpdateCategory(t *testing.T) {
	//initial length of []products
	initialLen := len(Categories)
	//parameters passed to request body
	requestBody := &Category{
		CategoryName: 		"Super Cool Category",
		CategoryDescription: "Brand new cool Category",
	}
	jsonProduct, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("PATCH", "/categories/bq4fasj7jhfi127rimlg", bytes.NewBuffer(jsonProduct))
	req = mux.SetURLVars(req, map[string]string{"id": "bq4fasj7jhfi127rimlg"})

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateCategory)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code, "OK response is expected")
	assert.Equal(t, initialLen, len(Categories), "Expected length to stay the same after creating new product")
}

//TestUpdateCategoryWrongJSONSyntax tests whether UpdateCategory func returns the right status and does update fields
//because of the wrong syntax in JSON request body
func TestUpdateCategoryWrongJSONSyntax(t *testing.T) {
	//initial length of []products
	initialLen := len(Categories)
	//parameters passed to request body
	requestBody := `{{"CategoryID":"bq4fasj7jhfi127rimlg","CategoryName":"Name",,,}}`
	req, err := http.NewRequest("PATCH", "/categories/bq4fasj7jhfi127rimlg", bytes.NewBufferString(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": "bq4fasj7jhfi127rimlg"})

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateCategory)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 400, rr.Code, "Bad request response is expected")
	assert.Equal(t, initialLen, len(Categories), "Expected length to stay the same after updating product")
}