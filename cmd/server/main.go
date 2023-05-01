package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/renatospaka/library/configs"
	"github.com/renatospaka/library/internal/entity"
	"github.com/renatospaka/library/internal/infra/database"
	"github.com/renatospaka/library/internal/infra/webservers/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	log.Println("iniciando a aplicação...")
	configs, err := configs.LoadConfig(".")
	if err != nil {
		log.Panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	productDB := database.NewProduct(db)
	ProductHandler := handlers.NewProductHandler(productDB)
	r.Post("/products", ProductHandler.CreateProduct)
	r.Get("/products", ProductHandler.GetProducts)
	r.Get("/products/{id}", ProductHandler.GetProduct)
	r.Put("/products/{id}", ProductHandler.UpdateProduct)
	r.Delete("/products/{id}", ProductHandler.DeleteProduct)
	
	userDB := database.NewUser(db)
	UserHandler := handlers.NewUserHandler(userDB, configs.TokenAuth, configs.JWTExpiresIn)
	r.Post("/users", UserHandler.Create)
	r.Post("/users/generate_token", UserHandler.GetJWT)

	log.Println("servidor escutando porta:", 8000)
	http.ListenAndServe(":8000", r)
}
