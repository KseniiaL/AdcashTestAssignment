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

type product struct {
	ProductID          string `json:"ProductID"`
	ProductName        string `json:"ProductName"`
	ProductDescription string `json:"ProductDescription"`
	Price			   int	  `json:"Price"`
	CategoryID 		   string `json:"CategoryID"`
}

type allProducts []product

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

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(products)
	log.Println("GET: Products")
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["id"]

	var Prod product
	for _, singleProduct := range products {
		if singleProduct.ProductID == productID {
			Prod = singleProduct
			break
		}
	}

	if Prod.ProductID == productID {
		json.NewEncoder(w).Encode(Prod)
	} else {
		fmt.Fprintf(w, "Product with ID %s not found", productID)
	}

	log.Println("GET: Products/", productID)
}

func GetProductsOfCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := mux.Vars(r)["id"]

	productsOfCategory := make([]product, 0)
	for _, singleProduct := range products {
		if singleProduct.CategoryID == categoryID {
			productsOfCategory = append(productsOfCategory, singleProduct)
		}
	}

	json.NewEncoder(w).Encode(productsOfCategory)
	
	log.Println("GET: Products of category ", categoryID)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["id"]
	productsLength := len(products)

	for i, singleProduct := range products {
		if singleProduct.ProductID == productID {
			products = append(products[:i], products[i+1:]...)
			fmt.Fprintf(w, "The category with ID %v has been deleted successfully", productID)
		}
	}

	if productsLength == len(products) {
		fmt.Fprintf(w, "Product with ID %s not found", productID)
	}
	log.Println("DELETE: Products/", productID)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct product

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err.Error())
		fmt.Fprintf(w, "Kindly enter data with the product name and description only in order to create")
	}

	json.Unmarshal(reqBody, &newProduct)
	newProduct.ProductID = xid.New().String()

	var usedCategory categories.Category
	for _, singleCategory := range categories.Categories {
		if singleCategory.CategoryID == newProduct.CategoryID {
			usedCategory = singleCategory
			break
		}
	}

	if len(newProduct.CategoryID) != 0 && usedCategory.CategoryID == newProduct.CategoryID {
		products = append(products, newProduct)
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(newProduct)
	} else {
		fmt.Fprintf(w, "Category with ID \"%s\" not found. Kindly enter data with the category ID", newProduct.CategoryID)
	}

	log.Println("POST: Products")
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["id"]
	var updateProduct product
	//TODO handle category not found
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the product name and description only in order to update")
	}
	json.Unmarshal(reqBody, &updateProduct)

	for i, singleProduct := range products {
		if singleProduct.ProductID == productID {
			singleProduct.ProductName = updateProduct.ProductName
			singleProduct.ProductDescription = updateProduct.ProductDescription
			singleProduct.Price = updateProduct.Price

			var usedCategory categories.Category
			for _, singleCategory := range categories.Categories {
				if singleCategory.CategoryID == updateProduct.CategoryID {
					usedCategory = singleCategory
					break
				}
			}
			if len(usedCategory.CategoryID) != 0 {
				singleProduct.CategoryID = updateProduct.CategoryID
			}

			products = append(products[:i], singleProduct)
			json.NewEncoder(w).Encode(singleProduct)
		}
	}
	log.Println("PATCH: Products/", productID)
}