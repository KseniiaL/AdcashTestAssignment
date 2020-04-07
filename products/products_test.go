//package products contains test for products.go
package products

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

//prodByIdTest is a structure for testing GET /product/{id} requests
var prodByIdTest = []struct {
	id       	 string // id
	expected 	 string // expected result
	expectedCode int //expected http response status code
}{
	{
		"bq4foj37jhfipc5nqri0",
		`{"ProductID":"bq4foj37jhfipc5nqri0","ProductName":"Nike SuperRep Go","ProductDescription":"Women's Training Shoe","Price":100,"CategoryID":"bq4fasj7jhfi127rimlg"}`,
		200,
	},
	{
		"bq5457j7jhfi2s58o030",
		`{"ProductID":"bq5457j7jhfi2s58o030","ProductName":"Nike Icon Clash","ProductDescription":"Women's Seamless Light-Support Sports Bra","Price":50,"CategoryID":"bq4fasj7jhfi127rimlg"}`,
		200,
	},
	{
		"randomID",
		`Product with ID randomID not found`,
		412,
	},
}

//prodByCategoryTest is a structure for testing GET /product/category/{id} requests
var prodByCategoryTest = []struct {
	categoryId   string // category id
	expected 	 string // expected result
	expectedCode int //expected http response status code
}{
	{
		"bq4fasj7jhfi127rimlg",
		`[{"ProductID":"bq4foj37jhfipc5nqri0","ProductName":"Nike SuperRep Go","ProductDescription":"Women's Training Shoe","Price":100,"CategoryID":"bq4fasj7jhfi127rimlg"},{"ProductID":"bq5457j7jhfi2s58o030","ProductName":"Nike Icon Clash","ProductDescription":"Women's Seamless Light-Support Sports Bra","Price":50,"CategoryID":"bq4fasj7jhfi127rimlg"}]`,
		200,
	},
	{
		"randomID",
		`[]`,
		200,
	},
}

//TestGetAllProducts tests whether GetAllProducts func returns the right response body and status
func TestGetAllProducts(t *testing.T) {
	//Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/products", nil)
	if err != nil {
		t.Fatal(err)
	}
	//Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllProducts)

	handler.ServeHTTP(rr, req)
	// Check the response status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"ProductID":"bq4foj37jhfipc5nqri0","ProductName":"Nike SuperRep Go","ProductDescription":"Women's Training Shoe","Price":100,"CategoryID":"bq4fasj7jhfi127rimlg"},{"ProductID":"bq5457j7jhfi2s58o030","ProductName":"Nike Icon Clash","ProductDescription":"Women's Seamless Light-Support Sports Bra","Price":50,"CategoryID":"bq4fasj7jhfi127rimlg"}]`
	assert.JSONEq(t, expected, rr.Body.String(), "Expected response body to be the same")
}

//TestGetProductById tests whether GetProductById func returns the right response bodies and statuses
//while iterating over prodByIdTest
func TestGetProductById(t *testing.T) {
	for _, p := range prodByIdTest {
		//request url
		path := "/products/" + p.id
		req, err := http.NewRequest("GET", path, nil)
		//passing the {id} var as it was not recognizable for some reason
		req = mux.SetURLVars(req, map[string]string{"id": p.id})
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetProductById)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, p.expectedCode, rr.Code, "OK response is expected")
		assert.Equal(t, p.expected, strings.TrimSuffix(rr.Body.String(),"\n"), "Response body is expected to be equal to expected value")
	}
}

//TestGetProductsOfCategory tests whether GetProductsOfCategory func returns the right response bodies and statuses
//while iterating over prodByCategoryTest
func TestGetProductsOfCategory(t *testing.T) {
	for _, p := range prodByCategoryTest {
		path := "/products/category/" + p.categoryId
		req, err := http.NewRequest("GET", path, nil)
		req = mux.SetURLVars(req, map[string]string{"id": p.categoryId})

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetProductsOfCategory)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, p.expectedCode, rr.Code, "Response code is expected to be different")
		assert.Equal(t, p.expected, strings.TrimSuffix(rr.Body.String(),"\n"), "Response body is expected to be equal to expected value")
	}
}

//TestDeleteProduct tests whether DeleteProduct func returns the right status and actually deletes a product from []products
func TestDeleteProduct(t *testing.T) {
	//initial length of []products
	initialLen := len(products)

	req, err := http.NewRequest("DELETE", "/products/bq4foj37jhfipc5nqri0", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "bq4foj37jhfipc5nqri0"})

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteProduct)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code, "OK response is expected")
	//the length of []products should decrease after deleting product
	assert.NotEqual(t, initialLen, len(products), "Expected length to decrease after creating new product")
}

//TestDeleteProductWrongID tests whether DeleteProduct func returns the right status and does not delete a product from []products
//because the ID is wrong
func TestDeleteProductWrongID(t *testing.T) {
	//initial length of []products
	initialLen := len(products)

	req, err := http.NewRequest("DELETE", "/products/randomID", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "randomID"})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteProduct)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 412, rr.Code, "Precondition Failed response is expected")
	//the length of []products should not change after trying to delete non-existing product
	assert.Equal(t, initialLen, len(products), "Expected length to stay same after creating new product")
}

//TestCreateProduct tests whether CreateProduct func returns the right status and actually appends a product to []products
func TestCreateProduct(t *testing.T) {
	//initial length of []products
	initialLen := len(products)
	//parameters passed to request body
	requestBody := &product{
		ProductName: 		"Super Cool Product",
		ProductDescription: "Brand new cool product",
		Price: 				1000,
		CategoryID: 		"bq4fasj7jhfi127rimlg",
	}
	jsonProduct, _ := json.Marshal(requestBody)
	//Create a request to pass to the handler with request body as a third parameter
	req, err := http.NewRequest("POST", "/products/new", bytes.NewBuffer(jsonProduct))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateProduct)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 201, rr.Code, "Created response is expected")
	//the length of []products should increase after creating new product
	assert.NotEqual(t, initialLen, len(products), "Expected length to increase after creating new product")
}

//TestCreateProductNonExistingCategory tests whether CreateCategory func returns the right status and does not append
//a product to []products because of the non-existing category to refer to
func TestCreateProductNonExistingCategory(t *testing.T) {
	//initial length of []products
	initialLen := len(products)
	//parameters passed to request body with wrong CategoryID
	requestBody := &product{
		ProductName: 		"Super Cool Product",
		ProductDescription: "Brand new cool product",
		Price: 				1000,
		CategoryID: 		"randomCategoryID",
	}
	jsonProduct, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/products/new", bytes.NewBuffer(jsonProduct))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateProduct)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 422, rr.Code, "Unprocessable Entity response is expected")
	//the length of []products should not change after trying to create new product but passing wrong CategoryID
	//product should be connected to the existing category
	assert.Equal(t, initialLen, len(products), "Expected length to stay same after creating new product")
}

//TestCreateProductEmptyBody tests whether CreateCategory func returns the right status and does not append
//a product to []products because of the empty request body
func TestCreateProductEmptyBody(t *testing.T) {
	//initial length of []products
	initialLen := len(products)
	//empty body
	requestBody := &product{}
	jsonProduct, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/products/new", bytes.NewBuffer(jsonProduct))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateProduct)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 422, rr.Code, "Created response is expected")
	//the length of []products should not change after trying to create new empty product
	assert.Equal(t, initialLen, len(products), "Expected length to increase after creating new product")
}

//TestCreateProductEmptyBody tests whether CreateCategory func returns the right status and does not append
//a product to []products because of the wrong syntax in JSON request body
func TestCreateProductWrongJSONSyntax(t *testing.T) {
	//initial length of []Categories
	initialLen := len(products)
	//parameters passed to request body
	requestBody := `{{"ProductID":"bq4foj37jhfipc5nqri0","ProductName":"Name",,,}}`
	req, err := http.NewRequest("POST", "/products/new", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateProduct)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 400, rr.Code, "Bad request response is expected")
	assert.Equal(t, initialLen, len(products), "Expected length to stay the same after wrong syntax json")

}

//TestUpdateProduct tests whether UpdateProduct func returns the right status and does not change []products, but updates fields
func TestUpdateProduct(t *testing.T) {
	//initial length of []products
	initialLen := len(products)
	//parameters passed to request body
	requestBody := &product{
		ProductName: 		"Super Cool Product",
		ProductDescription: "Brand new cool product",
		Price: 				1000,
		CategoryID: 		"bq4fasj7jhfi127rimlg",
	}
	jsonProduct, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("PATCH", "/products/bq4foj37jhfipc5nqri0", bytes.NewBuffer(jsonProduct))
	req = mux.SetURLVars(req, map[string]string{"id": "bq4foj37jhfipc5nqri0"})

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateProduct)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code, "OK response is expected")
	assert.Equal(t, initialLen, len(products), "Expected length to stay the same after creating new product")
}

//TestUpdateProductWrongJSONSyntax tests whether UpdateProduct func returns the right status and does update fields
//because of the wrong syntax in JSON request body
func TestUpdateProductWrongJSONSyntax(t *testing.T) {
	//initial length of []Categories
	initialLen := len(products)
	//parameters passed to request body
	requestBody := `{{"ProductID":"bq4foj37jhfipc5nqri0","ProductName":"Name",,,}}`
	req, err := http.NewRequest("POST", "/products/bq4foj37jhfipc5nqri0", bytes.NewBufferString(requestBody))
	req = mux.SetURLVars(req, map[string]string{"id": "bq4foj37jhfipc5nqri0"})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateProduct)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 400, rr.Code, "Bad request response is expected")
	assert.Equal(t, initialLen, len(products), "Expected length to stay the same after wrong syntax json")

}