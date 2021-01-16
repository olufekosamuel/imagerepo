package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/olufekosamuel/imagerepo/helpers"
	"github.com/olufekosamuel/imagerepo/users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Register(w http.ResponseWriter, r *http.Request) {

	// add cors support
	helpers.SetupResponse(&w, r)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if (*r).Method == "POST" {
		var user models.User
		_ = json.NewDecoder(r.Body).Decode(&user)

		if len(user.Email) == 0 {
			http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "email can not be empty"), 400)
			return
		}

		if user.Password == "" {
			http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "Password must be provided"), 400)
			return
		}

		collection, err := helpers.GetDBCollection("users")

		if err != nil {
			logger.Log(err)
			http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "couldn't connect to the mongo collection"), 500)
			return
		}

		var query_result primitive.M

		err = collection.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Decode(&query_result)

		if len(query_result) != 0 {
			http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "User already exists, please login"), 401)
			return
		}

		hash, err := helpers.HashPassword(user.Password)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "error occured while hashing password"), 401)
			return
		}

		user.CreatedAt = time.Now().Format(time.RFC3339)
		user.UpdatedAt = time.Now().Format(time.RFC3339)
		user.Password = hash

		_, err = collection.InsertOne(context.TODO(), user)

		json.NewEncoder(w).Encode(models.Response{
			Status: "success",
			Error:  false,
			Msg:    "User account created successfully",
		})
		return
	}
	http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "method not allowed"), 400)
	return
}
