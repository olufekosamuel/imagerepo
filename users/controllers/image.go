package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/olufekosamuel/imagerepo/auth"
	"github.com/olufekosamuel/imagerepo/helpers"
	"github.com/olufekosamuel/imagerepo/users/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 20 // 20MB
/*
Endpoint to upload image
*/
func ImageUpload(w http.ResponseWriter, r *http.Request) {

	// add cors support
	helpers.SetupResponse(&w, r)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("content-type", "application/json")

	if (*r).Method == "POST" {

		id_available := false
		id, _, err := auth.ExtractUserIDEmail(r)

		if err != nil {

		} else {
			id_available = true
		}

		new_id, _ := primitive.ObjectIDFromHex(id)
		text := r.FormValue("text")
		image_type := r.FormValue("type")
		characteristics := r.Form["characteristics"]
		files := r.MultipartForm.File["image"]

		if len(text) < 1 {
			http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "Image text must be provided"), 400)
			return
		}

		if len(image_type) < 1 || image_type != "private" && image_type != "public" {
			http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "Image type must be provided, either private or public"), 400)
			return
		}

		if len(characteristics) < 1 {
			http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "Image characteristics must be provided"), 400)
			return
		}

		if !id_available {
			if image_type != "public" {
				http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "Image type must be made public or create an account to make it private"), 400)
				return
			}
		}

		for _, fileHeader := range files {

			if fileHeader.Size > MAX_UPLOAD_SIZE {
				http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 20MB in size", fileHeader.Filename), 400)
				return
			}

			// Open the file
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer file.Close()

			buff := make([]byte, 512)
			_, err = file.Read(buff)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			filetype := http.DetectContentType(buff)
			if filetype != "image/jpeg" && filetype != "image/png" {
				http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", 400)
				return
			}

			// create an AWS session which can be
			// reused if we're uploading many files
			s, err := session.NewSession(&aws.Config{
				Region: aws.String("eu-west-2"),
				Credentials: credentials.NewStaticCredentials(
					os.Getenv("S3_ACCESS_ID"), // id
					os.Getenv("S3_SECRET_ID"), // secret
					""),                       // token can be left blank for now
			})

			if err != nil {
				logger.Log(err)
				http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "error occured while trying to upload"), 400)
				return
			}

			fileName, err := helpers.UploadFileToS3(s, file, fileHeader, "images/")
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "error occured while trying to upload file"), 400)
				return
			}

			var image models.Image

			collection, err := helpers.GetDBCollection("images")

			if err != nil {
				logger.Log(err)
				http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "couldn't connect to the mongo collection"), 500)
				return
			}

			if id_available {
				image.UserID = new_id
			}

			image.Image = fileName
			image.Characteristics = characteristics
			image.Text = text
			image.Type = image_type
			image.CreatedAt = time.Now().Format(time.RFC3339)
			image.UpdatedAt = time.Now().Format(time.RFC3339)

			_, err = collection.InsertOne(context.TODO(), image)

			if err != nil {
				logger.Log(err)
				http.Error(w, fmt.Sprintf(`{"status":"error","msg":"error occured while trying to saving image"}`), 400)
				return
			}

		}

		json.NewEncoder(w).Encode(models.Response{
			Status: "success",
			Error:  false,
			Msg:    "User image(s) uploaded successfully",
		})
		return

	}
	http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "method not allowed"), 400)
	return
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

/*
Endpoint to search image
*/
func SearchImage(w http.ResponseWriter, r *http.Request) {

	// add cors support
	helpers.SetupResponse(&w, r)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("content-type", "application/json")

	if (*r).Method == "GET" {

		id_available := false
		id, _, err := auth.ExtractUserIDEmail(r)

		if err == nil {
			id_available = true
		}

		new_id, _ := primitive.ObjectIDFromHex(id)

		query := bson.M{}
		collection, err := helpers.GetDBCollection("images")
		var query_result primitive.M
		//check for filtering
		text, text_exists := r.URL.Query()["text"]
		characteristics, characteristics_exists := r.URL.Query()["characteristics"]
		image_id, id_eixsts := r.URL.Query()["id"]

		if id_eixsts && len(image_id[0]) > 1 { //check similar images of a specified image

			id, _ := primitive.ObjectIDFromHex(image_id[0])

			err = collection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&query_result)

			if err != nil {
				logger.Log(err)
				http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "image with that id not found"), 404)
				return
			}

			query["type"] = query_result["type"]

		}

		if text_exists && len(text[0]) > 1 { //if image text is supplied do a check
			query["text"] = primitive.Regex{Pattern: text[0], Options: "i"}
		}

		imageList := make([]models.Image, 0)
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "couldn't connect to the mongo collection"), 500)
			return
		}

		findOptions := options.Find()
		findOptions.SetSort(bson.D{{"createdat", -1}})
		cursor, err := collection.Find(context.TODO(), query, findOptions)

		for cursor.Next(ctx) {
			// Declare a result BSON object
			var image models.Image

			//var cover models.Cover
			err := cursor.Decode(&image)
			if err != nil {
				logger.Log(err)
				http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "error occured while getting a image"), 500)
				return
			}

			if image.Type == "private" { //means a authenticated user is making a search, so lets show him images that are privately owned by him/her
				if !id_available && image.UserID != new_id { //check if the user owns the image
					break
				}
			}

			image.Image = "https://imgrepo-test.s3.eu-west-2.amazonaws.com/" + image.Image

			if characteristics_exists && len(characteristics[0]) > 1 { //if characteristics are supplied, then do a check
				for i := range characteristics {
					if contains(image.Characteristics, characteristics[i]) { //check if this particular image contains any of the characteristics requested for.
						imageList = append(imageList, image)
						break
					}
				}
			} else {
				imageList = append(imageList, image)
			}

		}

		json.NewEncoder(w).Encode(models.Response{
			Status: "success",
			Error:  false,
			Data: models.DataResponse{
				"images": imageList,
			},
		})
		return

	}
	http.Error(w, fmt.Sprintf(`{"status":"error","error":true,"msg":"%s"}`, "method not allowed"), 400)
	return
}
