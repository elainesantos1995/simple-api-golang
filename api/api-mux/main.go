package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var usersSlice []User

func main() {
	fmt.Println("REST API using mux is running...")

	fmt.Println("REST API users registered", usersSlice)

	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/users", getUsers).Methods("GET")
	muxRouter.HandleFunc("/users/{id}", getUserByID).Methods("GET")
	muxRouter.HandleFunc("/users", createUser).Methods("POST")
	muxRouter.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", muxRouter))
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getUsers(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	fmt.Println("Returning results...")
	json.NewEncoder(writer).Encode(usersSlice)
}

func getUserByID(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	params := mux.Vars(request)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(writer, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// finding user for id into slice
	foundUser := User{}
	for _, user := range getUsersList() {
		if user.ID == userID {
			foundUser = user
			break
		}
	}

	if foundUser.ID == 0 {
		http.Error(writer, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(writer).Encode(foundUser)
}

// Assist function for to get users slice
func getUsersList() []User {
	return usersSlice
}

func createUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	var newUser User
	err := json.NewDecoder(request.Body).Decode(&newUser)
	if err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	// adding new user into slice users
	addUser(&newUser)

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(newUser)
}

func addUser(newUser *User) {
	usersSlice = append(usersSlice, *newUser)

	fmt.Println("New users registered", len(usersSlice))
}

func deleteUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	params := mux.Vars(request)
	userID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(writer, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// search user list index
	index := -1
	for i, user := range usersSlice {
		if user.ID == userID {
			index = i
			break
		}
	}

	// Returning error if user not found
	if index == -1 {
		http.Error(writer, "User not found", http.StatusNotFound)
		return
	}

	// Removing user from user sslice
	usersSlice = append(usersSlice[:index], usersSlice[index+1:]...)

	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "User with ID %d deleted", userID)
	fmt.Println("User slice after removing: ", usersSlice)
}
