package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/ichtrojan/thoth"
	"github.com/joho/godotenv"
	"github.com/olufekosamuel/imagerepo/helpers"
	"github.com/olufekosamuel/imagerepo/users/controllers"
	"github.com/olufekosamuel/imagerepo/users/models"
)

var (
	logger, _ = thoth.Init("log")
)

func Index(w http.ResponseWriter, r *http.Request) {
	// add cors support
	helpers.SetupResponse(&w, r)
	w.Header().Set("content-type", "application/json")

	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	if (*r).Method == "GET" {
		json.NewEncoder(w).Encode(models.Response{
			Status: "success",
			Error:  false,
			Msg:    "Server working fine",
		})
		return
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "Endpoint Not Found"), 404)
	}
}

func main() {

	if err := godotenv.Load(); err != nil {
		logger.Log(errors.New("no .env file found"))
		log.Fatal("No .env file found")
	}

	port, exist := os.LookupEnv("PORT")

	if !exist {
		logger.Log(errors.New("PORT not set in .env"))
		log.Fatal("PORT not set in .env")
	}

	//Test Endpoint is working fine
	http.HandleFunc("/", Index)

	//Endpoints
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/image", controllers.ImageUpload)
	http.HandleFunc("/search", controllers.SearchImage)

	port = fmt.Sprintf(":%s", port)

	fmt.Println(fmt.Sprintf("application is running on port %s", port))

	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Log(err)
	}

}
