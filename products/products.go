//package products contains all the requested methods for working with products table
package products

import (
	"encoding/json"
	"fmt"
	"github.com/KseniiaL/AdcashTestAssignment/categories"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
	"io/ioutil"
	"log"
	"net/http"
)

// product stores information about product fields.
type product struct {
	ProductID          string `json:"ProductID"`
	ProductName        string `json:"ProductName"`
	ProductDescription string `json:"ProductDescription"`
	Price			   int	  `json:"Price"`
	CategoryID 		   string `json:"CategoryID"`
}

// allProducts is the slice of product structs
type allProducts []product

// products is the simple imitation of the DB (products table)
var products = allProducts {
	{
		ProductID: 			"bq4foj37jhfipc5nqri0",
		ProductName: 		"Nike SuperRep Go",
		ProductDescription: "Women's Training Shoe",
		Price: 				100,
		CategoryID: 		"bq4fasj7jhfi127rimlg",
	},
	{
		ProductID: 			"bq5457j7jhfi2s58o030",
		ProductName: 		"Nike Icon Clash",
		ProductDescription: "Women's Seamless Light-Support Sports Bra",
		Price: 				50,
		CategoryID: 		"bq4fasj7jhfi127rimlg",
	},
}

// GetAllProducts returns products slice in JSON format as a response
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	//define new encoder that writes to the w
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf(err.Error())
		w.WriteHeader(500)
	}
}

// GetProductById gets a product id from the request link and looks for the corresponding item in []products
func GetProductById(w http.ResponseWriter, r *http.Request) {
	//get product id from the link
	productID := mux.Vars(r)["id"]

	//find the product with the given id in the slice
	var prod product
	for _, singleProduct := range products {
		if singleProduct.ProductID == productID {
			prod = singleProduct
			break
		}
	}

	//return the product information to ResponseWriter
	//or log the encoding error
	if prod.ProductID == productID {
		if err := json.NewEncoder(w).Encode(prod); err != nil {
			log.Printf(err.Error())
			w.WriteHeader(500)
			return
		}
	} else {
		fmt.Fprintf(w, "Product with ID %s not found", productID)
	}
}

// GetProductsOfCategory gets a category id from the request link and returns all products of the given category in response
func GetProductsOfCategory(w http.ResponseWriter, r *http.Request) {
	//get category id from the link
	categoryID := mux.Vars(r)["id"]

	//define new slice and append it with products which have the same categoryID
	productsOfCategory := make([]product, 0)
	for _, singleProduct := range products {
		if singleProduct.CategoryID == categoryID {
			productsOfCategory = append(productsOfCategory, singleProduct)
		}
	}

	//return the products list to ResponseWriter
	//or log the encoding error
	if err := json.NewEncoder(w).Encode(productsOfCategory); err != nil {
		log.Printf(err.Error())
		w.WriteHeader(500)
	}
}

// DeleteProduct gets a product id from the request link and removes corresponding item from the slice
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	//get product id from the link
	productID := mux.Vars(r)["id"]
	productsLength := len(products)

	//find the product with the given id and remove from the slice
	for i, singleProduct := range products {
		if singleProduct.ProductID == productID {
			products = append(products[:i], products[i+1:]...)
			fmt.Fprintf(w, "The category with ID %v has been deleted successfully", productID)
		}
	}

	//report product with the given id not exists
	if productsLength == len(products) {
		fmt.Fprintf(w, "Product with ID %s not found", productID)
	}
}

// CreateProduct creates a new sample of product, fills it with the information from the request body,
// and appends to the []products slice
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct product

	//get the information containing in request's body
	//or report an error
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the product name and description only in order to create")
		log.Fatal(err.Error())
	}

	//unmarshal the information from JSON into the product instance
	//or report an error
	if err = json.Unmarshal(reqBody, &newProduct); err != nil {
		log.Printf("Body parse error, %v", err.Error())
		w.WriteHeader(400)
		return
	}
	//generate unique productID
	newProduct.ProductID = xid.New().String()

	//check if the category given exists in categories
	var usedCategory categories.Category
	for _, singleCategory := range categories.Categories {
		if singleCategory.CategoryID == newProduct.CategoryID {
			usedCategory = singleCategory
			break
		}
	}

	//if exists append the new product to the slice
	if len(newProduct.CategoryID) != 0 && usedCategory.CategoryID == newProduct.CategoryID {
		products = append(products, newProduct)
		w.WriteHeader(http.StatusCreated)

		//return the product in response
		//or report an error
		if err = json.NewEncoder(w).Encode(newProduct); err != nil {
			log.Printf(err.Error())
			w.WriteHeader(500)
			return
		}
	} else {
		fmt.Fprintf(w, "Category with ID \"%s\" not found. Kindly enter data with the category ID", newProduct.CategoryID)
	}
}

// UpdateProduct gets a product id from the request link and replaces the fields in the corresponding product
// with the given ones in the request body
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	//get product id from the link
	productID := mux.Vars(r)["id"]
	var updateProduct product

	//get the information containing in request's body
	//or report an error
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the product name and description only in order to update")
	}

	//unmarshal the information from JSON into the product instance
	//or report an error
	if err = json.Unmarshal(reqBody, &updateProduct); err != nil {
		log.Printf("Body parse error, %v", err.Error())
		w.WriteHeader(400)
		return
	}

	//find the given product in the slice by id
	for i, singleProduct := range products {
		if singleProduct.ProductID == productID {
			//change the fields
			singleProduct.ProductName = updateProduct.ProductName
			singleProduct.ProductDescription = updateProduct.ProductDescription
			singleProduct.Price = updateProduct.Price

			//check if a categoryID exists in the categories
			var usedCategory categories.Category
			for _, singleCategory := range categories.Categories {
				if singleCategory.CategoryID == updateProduct.CategoryID {
					usedCategory = singleCategory
					break
				}
			}
			//then replace the categoryID in the product
			if len(usedCategory.CategoryID) != 0 {
				singleProduct.CategoryID = updateProduct.CategoryID
			}

			products = append(products[:i], singleProduct)
			//return the product in response
			//or report an error
			if err = json.NewEncoder(w).Encode(singleProduct); err != nil {
				log.Printf(err.Error())
				w.WriteHeader(500)
				return
			}
		}
	}
}