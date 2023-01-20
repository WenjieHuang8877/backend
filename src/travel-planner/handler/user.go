package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"travel-planner/model"

	"time"

	"travel-planner/backend"
	"travel-planner/service"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
)

func signinHandler(w http.ResponseWriter, r *http.Request){
  fmt.Println("Received sign request")
  w.Header().Set("Content-Type","text/plain")

  // Get User infor from client
  decoder := json.NewDecoder(r.Body)
  var user model.User
  if err := decoder.Decode(&user); err != nil {
	http.Error(w, "Cannot decode user data from client", http.StatusBadRequest);
	fmt.Printf("Cannot decode user data from client %v\n", err)
	return
  }
  // check user exist or not
  exist, err := service.CheckUser(user.Email, user.Password)
  if err != nil {
	http.Error(w, "Failed to read user from database",http.StatusInternalServerError)
	fmt.Printf("Failed to read user from database %v\n", err)
	return
  }

  if !exist {
	http.Error(w, "User doesn't exist or wrong password", http.StatusUnauthorized)
	fmt.Printf("User doesn't exists or wrong password\n")
    return
  }

  // generate token
  token := jwt.NewWithClaims(jwt.SigningMethodES256,jwt.MapClaims{
	"email": user.Email,
	"exp" :   time.Now().Add(time.Hour * 24).Unix(),
  })
  // sign and get the complete encoded token as a string using the secret
  tokenString, err := token.SignedString(mySigningKey)

  if err != nil {
	http.Error(w, "Failed to generate token", http.StatusInternalServerError)
	fmt.Printf("Failed to generate token %v\n", err)
	return
  }

  w.Write([]byte(tokenString))

}

func signupHandlerer(w http.ResponseWriter, r *http.Request){
   fmt.Println("Received a signup request")
   w.Header().Set("Content-Type","text/plain")
   
   decoder := json.NewDecoder(r.Body)

   var user model.User

   if err := decoder.Decode(&user); err != nil{
	http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
	fmt.Printf("Cannot decode user data from client %v\n", err)
	return
   }

   if user.Email == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9.@]$`).MatchString(user.Email){
     http.Error(w,"Invalid email address or password", http.StatusBadRequest)
	 fmt.Printf("Invalid email address or password\n");
   }

   user.Id = uuid.New()
   // check if the user already existed
   userFound, _ := backend.DB.ReadUserByEmail(user.Email)

  if userFound != nil {
	http.Error(w, "The user has been registered before", http.StatusBadRequest)
	fmt.Println("The user has been registered before")
    return
  }

  success, err := service.AddUser(&user)

  if err != nil {
	http.Error(w, "Failed to save user to database", http.StatusInternalServerError)
	fmt.Printf("Failed to save to databse %v\n", err)
	return 
  }

  if !success {
	 http.Error(w, "User already exists", http.StatusBadRequest)
        fmt.Println("User already exists")
        return
  }

   fmt.Printf("User added successfully: %s.\n", user.Email)
}

////? 传送*user可？
func getUserHandler(w http.ResponseWriter, r *http.Request){
  fmt.Println("Received a get user information request")
  w.Header().Set("Content-Type","application/json")

  
  id := mux.Vars(r)["id"]

  
  user, err := service.CheckUserInfo(id)

  if err != nil {
	http.Error(w, "Failed to read user info from backend", http.StatusInternalServerError)
	return
  }
//? 传送*user可？
  js, err := json.Marshal(user)

  if err != nil {
	http.Error(w, "Failed to parse User into JSON format", http.StatusInternalServerError)
	return
  }

  w.Write(js)
  
}
// update interface has no return value in gorm?
func UpdateUserHander(w http.ResponseWriter, r *http.Request){
	fmt.Println("Received a request for updating user's information")

	user := r.Context().Value("user")
	fmt.Println(user)
	
   
	
	id := mux.Vars(r)["id"]
	fmt.Println(id)
	password := r.FormValue("password")
	username := r.FormValue("username")
	gender := r.FormValue("gender")
    age, _:=strconv.ParseInt(r.FormValue("age"), 10, 64)
     
	success, err :=service.UpdateUserInfo(id, password, username,gender, age)
	
	if !success {
		http.Error(w, "Failed to update user to backend",http.StatusInternalServerError)
		fmt.Printf("Failed to update user to backend %v\n ", err)
	}

	fmt.Println("User is updated successfully")
	fmt.Fprintf(w, "Update request received %s\n", id)

}