//package main initializes router and describes all the routes and funcs to handle
package main

import (
	"fmt"
	"github.com/KseniiaL/AdcashTestAssignment/categories"
	"github.com/KseniiaL/AdcashTestAssignment/products"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/categories", categories.GetAllCategories).Methods("GET")
	router.HandleFunc("/categories/{id}", categories.GetCategoryById).Methods("GET")
	router.HandleFunc("/categories/new", categories.CreateCategory).Methods("POST")
	router.HandleFunc("/categories/{id}", categories.DeleteCategory).Methods("DELETE")
	router.HandleFunc("/categories/{id}", categories.UpdateCategory).Methods("PATCH")
	router.HandleFunc("/products", products.GetAllProducts).Methods("GET")
	router.HandleFunc("/products/{id}", products.GetProductById).Methods("GET")
	router.HandleFunc("/products/new", products.CreateProduct).Methods("POST")
	router.HandleFunc("/products/{id}", products.UpdateProduct).Methods("PATCH")
	router.HandleFunc("/products/{id}", products.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/products/category/{id}", products.GetProductsOfCategory).Methods("GET")
	fmt.Println("Server running on: 8080")
	//run the server
	log.Fatal(http.ListenAndServe(":8080", router))
}