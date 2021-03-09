package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type Product struct {
	ID    int             `json:"id"`
	Code  string          `json:"code"`
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price" gorm:"type:decimal(16,2);"`
}
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	dbSourcePath := "root:dikihenas@tcp(localhost:3306)/go_rest_api_crud?charset=utf8&&parseTime=True"
	db, err = gorm.Open(mysql.Open(dbSourcePath), &gorm.Config{})
	if err != nil {
		log.Println("failed to connect database ", err)
	} else {
		log.Println("Connection established")
	}
	db.AutoMigrate(&Product{})

	handleRequest()
}

func handleRequest() {
	log.Println("Start Development on localhost:8081")
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/api/products", getProductList).Methods("GET")
	router.HandleFunc("/api/products/{id}", detailProduct).Methods("GET")
	router.HandleFunc("/api/products", createProduct).Methods("POST")
	router.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	router.HandleFunc("/api/products/{id}", deletProduct).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8081", router))

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to HomePage")
}
func getProductList(w http.ResponseWriter, r *http.Request) {
	product := []Product{}

	db.Find(&product)

	res := Result{
		Code:    200,
		Data:    product,
		Message: "Succes get List of Product",
	}

	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

func detailProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	vars := mux.Vars(r)
	productID := vars["id"]

	db.First(&product, productID)

	res := Result{
		Code:    200,
		Data:    product,
		Message: "Succes Product Detail",
	}

	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	json.NewDecoder(r.Body).Decode(&product)
	db.Create(&product)

	res := Result{
		Code:    200,
		Data:    product,
		Message: "Succes Add to DB",
	}

	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
func updateProduct(w http.ResponseWriter, r *http.Request) {
	var updateData Product
	vars := mux.Vars(r)
	productID := vars["id"]
	json.NewDecoder(r.Body).Decode(&updateData)

	var product Product
	db.First(&product, productID)
	db.Model(&product).Updates(updateData)

	res := Result{
		Code:    200,
		Data:    product,
		Message: "Succes Update Data Product",
	}

	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

}
func deletProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	vars := mux.Vars(r)
	productID := vars["id"]
	db.First(&product, productID)
	db.Delete(&product)

	res := Result{
		Code:    200,
		Data:    product,
		Message: "Succes Delete  Product",
	}

	result, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
