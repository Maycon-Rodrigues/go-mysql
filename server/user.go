package server

import (
	"encoding/json"
	"fmt"
	"go-mysql/db"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// StoreUser will create a User on DB
func StoreUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Error on get body"))
		return
	}

	var user user
	if err = json.Unmarshal(body, &user); err != nil {
		w.Write([]byte("Error to convert user on struct"))
		return
	}

	db, err := db.Connection()
	if err != nil {
		w.Write([]byte("Error to connect on data base"))
		return
	}
	defer db.Close()

	statemant, err := db.Prepare("INSERT INTO Users (name, email) VALUES (?, ?)")
	if err != nil {
		w.Write([]byte("Error to prepare statemant"))
		return
	}

	defer statemant.Close()

	result, err := statemant.Exec(&user.Name, &user.Email)
	if err != nil {
		w.Write([]byte("Error to execute statemant"))
		return
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		w.Write([]byte("Error getting id inserted "))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("User created on success! ID: %d", insertedID)))
}

// GetAllUsers will select all Users on DB
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db, err := db.Connection()
	if err != nil {
		w.Write([]byte("Error to connect on data base"))
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Users")
	if err != nil {
		w.Write([]byte("Error on get users"))
		return
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var user user

		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Error on scan user"))
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.Write([]byte("Error to convert user from json"))
		return
	}
}

// GetUser will select a User on DB
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		w.Write([]byte("Error on convert id to int"))
		return
	}

	db, err := db.Connection()
	if err != nil {
		w.Write([]byte("Error to connect on data base"))
		return
	}
	defer db.Close()

	row, err := db.Query("SELECT * FROM Users WHERE id = ?", ID)

	var user user
	if row.Next() {
		if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.Write([]byte("Error on scan user"))
			return
		}
	}

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.Write([]byte("Error to convert user from json"))
		return
	}
}

// UpdateUser will change user's data on DB
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 54)
	if err != nil {
		w.Write([]byte("Error on convert id to int"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Error on get body"))
		return
	}

	var user user
	if err := json.Unmarshal(body, &user); err != nil {
		w.Write([]byte("Error to convert user on struct"))
		return
	}

	db, err := db.Connection()
	if err != nil {
		w.Write([]byte("Error to connect on data base"))
		return
	}
	defer db.Close()

	statemant, err := db.Prepare("UPDATE Users SET name = ?, email = ? WHERE id = ?")
	if err != nil {
		w.Write([]byte("Error to prepare statemant"))
		return
	}
	defer statemant.Close()

	if _, err := statemant.Exec(user.Name, user.Email, ID); err != nil {
		w.Write([]byte("Error to execute statemant"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteUser will remove user on DB
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		w.Write([]byte("Error on convert id to int"))
		return
	}

	db, err := db.Connection()
	if err != nil {
		w.Write([]byte("Error to connect on data base"))
		return
	}
	defer db.Close()

	statemant, err := db.Prepare("DELETE FROM Users WHERE id = ?")
	if err != nil {
		w.Write([]byte("Error to prepare statemant"))
		return
	}
	defer statemant.Close()

	if _, err := statemant.Exec(ID); err != nil {
		w.Write([]byte("Error to execute statemant"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
