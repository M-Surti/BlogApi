package handlers

import (
	"BlogApi/database"
	"BlogApi/middleware"
	"BlogApi/models"
	"BlogApi/utils"
	"encoding/json"
	"fmt"

	"net/http"
	"strconv"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// collect the details of the user as request body
	var req *models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// check if user already exist
	var user models.User
	err = database.Db.Where("email = ?", req.Email).First(&user).Error
	if err == nil {
		http.Error(w, "user already exits", http.StatusBadRequest)
		return
	}

	// Hash the Password
	HashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "unable to hash password", http.StatusBadRequest)
		return
	}

	req.Password = HashPassword

	// add user to the database
	err = database.Db.Create(&req).Error
	if err != nil {
		http.Error(w, "unable to create database", http.StatusBadRequest)
	}

	// send a response
	w.WriteHeader(http.StatusOK)

}

func Login(w http.ResponseWriter, r *http.Request) {
	// decode request body

	var login models.User
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if user exist

	var user models.User
	err = database.Db.Where("email = ?", login.Email).First(&user).Error
	if err != nil {
		http.Error(w, "this user does not exist", http.StatusBadRequest)
		return
	}

	// check if password matches what we have in our database

	err = utils.ComparePassword(login.Password, user.Password)
	if err != nil {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	// uint ---> int ---> string
	idstr := strconv.Itoa(int(user.ID))

	// generating a token

	token, err := middleware.GenerateJWT(idstr)
	if err != nil {
		http.Error(w, "unable to generate token", http.StatusInternalServerError)
		return
	}

	// send a response
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(token)

}

func CreatePost(w http.ResponseWriter, r *http.Request){
	// decode the request body
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	post.UserID = userID

	err = database.Db.Create(&post).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "post created successfully")
}

func GetALLPost(w http.ResponseWriter, r *http.Request){
	var posts []models.Post
	err := database.Db.Find(&posts).Error
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
	
}


func GetPostByID(){

}

func DeletePOst(w http.ResponseWriter, r *http.Request){

}
