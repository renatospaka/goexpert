package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
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
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("JWTExpiresIn", configs.JWTExpiresIn))
	
	userDB := database.NewUser(db)
	UserHandler := handlers.NewUserHandler(userDB)
	r.Post("/users", UserHandler.Create)
	r.Post("/users/generate_token", UserHandler.GetJWT)

	productDB := database.NewProduct(db)
	ProductHandler := handlers.NewProductHandler(productDB)
	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/", ProductHandler.CreateProduct)
		r.Get("/", ProductHandler.GetProducts)
		r.Get("/{id}", ProductHandler.GetProduct)
		r.Put("/{id}", ProductHandler.UpdateProduct)
		r.Delete("/{id}", ProductHandler.DeleteProduct)		
	})

	log.Println("servidor escutando porta:", 8000)
	http.ListenAndServe(":8000", r)
}
