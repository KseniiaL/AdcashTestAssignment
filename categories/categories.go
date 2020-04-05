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

type Category struct {
	CategoryID  		string `json:"CategoryID"`
	CategoryName       	string `json:"CategoryName"`
	CategoryDescription string `json:"CategoryDescription"`
}

type allCategories []Category

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

func GetAllCategories (w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Categories)
	log.Println("GET: Categories")
}

func GetCategoryById (w http.ResponseWriter, r *http.Request) {
	categoryID := mux.Vars(r)["id"]

	for _, singleCategory := range Categories {
		if singleCategory.CategoryID == categoryID {
			json.NewEncoder(w).Encode(singleCategory)
		}
	}
	//TODO handle category not found
	log.Println("GET: Categories/", categoryID)
}

func CreateCategory (w http.ResponseWriter, r *http.Request) {
	var newCategory Category

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err.Error())
		fmt.Fprintf(w, "Kindly enter data with the category name and description only in order to update")
	}

	newCategory.CategoryID = xid.New().String()
	json.Unmarshal(reqBody, &newCategory)

	Categories = append(Categories, newCategory)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newCategory)
	log.Println("POST: Categories")
}

func DeleteCategory (w http.ResponseWriter, r *http.Request) {
	categoryID := mux.Vars(r)["id"]
	categoriesLength := len(Categories)

	for i, singleCategory := range Categories {
		if singleCategory.CategoryID == categoryID {
			Categories = append(Categories[:i], Categories[i+1:]...)
			fmt.Fprintf(w, "The category with ID %v has been deleted successfully", categoryID)
		}
	}

	if categoriesLength == len(Categories) {
		fmt.Fprintf(w, "Product with ID %s not found", categoryID)
	}
	log.Println("DELETE: Categories/", categoryID)
}

func UpdateCategory (w http.ResponseWriter, r *http.Request) {
	categoryID := mux.Vars(r)["id"]
	var updateCategory Category
	//TODO handle category not found
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the category name and description only in order to update")
	}
	json.Unmarshal(reqBody, &updateCategory)

	for i, singleCategory := range Categories {
		if singleCategory.CategoryID == categoryID {
			singleCategory.CategoryName = updateCategory.CategoryName
			singleCategory.CategoryDescription = updateCategory.CategoryDescription
			Categories = append(Categories[:i], singleCategory)
			json.NewEncoder(w).Encode(singleCategory)
		}
	}
	log.Println("PATCH: Categories/", categoryID)
}
