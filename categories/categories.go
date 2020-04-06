//package categories contains all the requested methods for working with categories table
package categories

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"io/ioutil"
	"log"
	"net/http"
)

// Category stores information about category fields
type Category struct {
	CategoryID  		string `json:"CategoryID"`
	CategoryName       	string `json:"CategoryName"`
	CategoryDescription string `json:"CategoryDescription"`
}

// allCategories is the slice of Category structs
type allCategories []Category

// Categories is the simple imitation of the DB (categories table)
var Categories = allCategories {
	{
		CategoryID:  		 "bq4fasj7jhfi127rimlg",
		CategoryName:        "Shopping Products",
		CategoryDescription: "Products consumers purchase and consume on a less frequent schedule compared to convenience products.",
	},
	{
		CategoryID:  		 "bq4fb3b7jhfi7v7uo39g",
		CategoryName:        "Specialty Products",
		CategoryDescription: "Products that are more expensive relative to convenience and shopping products.",
	},
}

// GetAllCategories returns Categories slice in JSON format as a response
func GetAllCategories (w http.ResponseWriter, r *http.Request) {
	//define new encoder that writes to the w
	//or report an error
	if err := json.NewEncoder(w).Encode(Categories); err != nil {
		log.Printf(err.Error())
		w.WriteHeader(500)
	}
}

// GetCategoryById gets a category id from the request link and looks for the corresponding item in []Categories
func GetCategoryById (w http.ResponseWriter, r *http.Request) {
	//get category id from the link
	categoryID := mux.Vars(r)["id"]

	//find the category with the given id in the slice
	var givenCategory Category
	for _, singleCategory := range Categories {
		if singleCategory.CategoryID == categoryID {
			givenCategory = singleCategory
		}
	}

	//return the Category information to ResponseWriter
	//or log the encoding error
	if givenCategory.CategoryID == categoryID {
		if err := json.NewEncoder(w).Encode(givenCategory); err != nil {
			log.Printf(err.Error())
			w.WriteHeader(500)
			return
		}
	} else {
		fmt.Fprintf(w, "Category with ID %s not found", categoryID)
	}
}

// CreateCategory creates a new sample of Category, fills it with the information from the request body,
// and appends to the []Categories slice
func CreateCategory (w http.ResponseWriter, r *http.Request) {
	var newCategory Category

	//get the information containing in request's body
	//or report an error
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the category name and description only in order to update")
		log.Fatal(err.Error())
	}

	//generate unique categoryID
	newCategory.CategoryID = xid.New().String()
	//unmarshal the information from JSON into the Category instance
	//or report an error
	if err = json.Unmarshal(reqBody, &newCategory); err != nil {
		log.Printf("Body parse error, %v", err.Error())
		w.WriteHeader(400)
		return
	}

	//append the new category to the slice
	Categories = append(Categories, newCategory)
	w.WriteHeader(http.StatusCreated)

	//return the category in response
	//or report an error
	if err = json.NewEncoder(w).Encode(newCategory); err != nil {
		log.Printf(err.Error())
		w.WriteHeader(500)
		return
	}
}

// DeleteCategory gets a category id from the request link and removes corresponding item from the slice
func DeleteCategory (w http.ResponseWriter, r *http.Request) {
	//get category id from the link
	categoryID := mux.Vars(r)["id"]
	//to define if some category was found in slice and deleted
	categoriesLength := len(Categories)

	//find the category with the given id and remove from the slice
	for i, singleCategory := range Categories {
		if singleCategory.CategoryID == categoryID {
			Categories = append(Categories[:i], Categories[i+1:]...)
			fmt.Fprintf(w, "The category with ID %v has been deleted successfully", categoryID)
		}
	}

	//report category with the given id not exists
	if categoriesLength == len(Categories) {
		fmt.Fprintf(w, "Category with ID %s not found", categoryID)
	}
}

// UpdateCategory gets a Category id from the request link and replaces the fields in the corresponding Category
// with the given ones in the request body
func UpdateCategory (w http.ResponseWriter, r *http.Request) {
	//get category id from the link
	categoryID := mux.Vars(r)["id"]
	var updateCategory Category

	//get the information containing in request's body
	//or report an error
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the category name and description only in order to update")
	}

	//unmarshal the information from JSON into the Category instance
	//or report an error
	if err = json.Unmarshal(reqBody, &updateCategory); err != nil {
		log.Printf("Body parse error, %v", err.Error())
		w.WriteHeader(400)
		return
	}

	//find the given Category in the slice by id
	for i, singleCategory := range Categories {
		if singleCategory.CategoryID == categoryID {
			//change the fields
			singleCategory.CategoryName = updateCategory.CategoryName
			singleCategory.CategoryDescription = updateCategory.CategoryDescription

			Categories = append(Categories[:i], singleCategory)

			//return the Category in response
			//or report an error
			if err = json.NewEncoder(w).Encode(singleCategory); err != nil {
				log.Printf(err.Error())
				w.WriteHeader(500)
				return
			}
		}
	}
}
