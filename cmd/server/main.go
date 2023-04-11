package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/renatospaka/library/configs"
	"github.com/renatospaka/library/internal/dto"
	"github.com/renatospaka/library/internal/entity"
	"github.com/renatospaka/library/internal/infra/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	log.Println("iniciando a aplicação...")
	_, err := configs.LoadConfig(".")
	if err != nil {
		log.Panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProduct(db)
	ProductHandler := NewProductHandler(productDB)

	log.Println("servidor escutando porta:", 8000)
	http.HandleFunc("/products", ProductHandler.CreateProduct)
	http.ListenAndServe(":8000", nil)
}


type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (ph *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// em geral esse acesso é feito pelo usercase
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = ph.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
