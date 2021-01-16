package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ichtrojan/thoth"
	"github.com/olufekosamuel/imagerepo/auth"
	"github.com/olufekosamuel/imagerepo/helpers"
	"github.com/olufekosamuel/imagerepo/users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	logger, _ = thoth.Init("log")
)

func Login(w http.ResponseWriter, r *http.Request) {

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

		if user.Password == "" {
			http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "Password not provided"), 400)
			return
		}

		if user.Email == "" {
			http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "Email not provided"), 400)
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

		if len(query_result) == 0 {
			http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "Wrong Login Credentials"), 400)
			return
		}

		if helpers.CheckPasswordHash(user.Password, query_result["password"].(string)) == true {
			tokenStr, err := auth.CreateUserToken(query_result["_id"].(primitive.ObjectID).Hex(), user.Email)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, err.Error()), 401)
				return
			}

			json.NewEncoder(w).Encode(models.Response{
				Status: "success",
				Error:  false,
				Data: models.DataResponse{
					"token": tokenStr,
				},
			})
			return
		}

		http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "Wrong login credentials"), 400)
		return

	}
	http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "method not allowed"), 400)
	return
}
