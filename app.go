package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	Db     *sql.DB
}

func (app *App) Initialize() error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", DbUser, DbPassword, Host, Port, Dbname)
	var err error
	app.Db, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	error_msg := map[string]string{"error": err}
	sendResponse(w, statusCode, error_msg)
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.Db)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, products)
}

func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	p := product{ID: id}

	err = p.getProductByID(app.Db)

	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, http.StatusNotFound, "Product not found")
		} else {
			sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	sendResponse(w, http.StatusOK, p)
}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	p := product{}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Inavlid payload")
		return
	}

	err = p.createProduct(app.Db)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusCreated, p)
}

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	p := product{ID: id}

	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Inavlid payload")
		return
	}

	err = p.updateProduct(app.Db)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, p)
}

func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	p := product{ID: id}

	err = p.deleteProduct(app.Db)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, map[string]string{"message": "Deleted Successfully"})
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")
}
