package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Member struct {
	User     string `db:"user" json:"user"`
	Password string `db:"pass" json:"pass"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Tel      string `db:"tel" json:"tel"`
}

type MemberResponse struct {
	User  string `db:"user" json:"user"`
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
	Tel   string `db:"tel" json:"tel"`
}

type repoDB struct {
	db *sqlx.DB
}

type memberHandler interface {
	NewMember(http.ResponseWriter, *http.Request)
	GetMember(http.ResponseWriter, *http.Request)
	EditMember(http.ResponseWriter, *http.Request)
	GetMemberByUser(http.ResponseWriter, *http.Request)
	DeleteMember(http.ResponseWriter, *http.Request)
}

func NewMemberHandler(db *sqlx.DB) memberHandler {
	return repoDB{db: db}
}

// NewMember implements memberHandler
func (repo repoDB) NewMember(w http.ResponseWriter, r *http.Request) {
	var userData Member
	err := json.NewDecoder(r.Body).Decode(&userData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Data Not Formate"))
	}
	query := "INSERT INTO member (user, pass, name, email, tel) VALUES (?,?,?,?,?)"
	arg := []any{}
	arg = append(arg, userData.User)
	arg = append(arg, userData.Password)
	arg = append(arg, userData.Name)
	arg = append(arg, userData.Email)
	arg = append(arg, userData.Tel)

	result, err := repo.db.Exec(query, arg...)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("database Error"))
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected < 1 {
		fmt.Println("Roweffact < 0")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("database Not Insert"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (repo repoDB) GetMember(w http.ResponseWriter, r *http.Request) {
	searchData := r.URL.Query()
	name := searchData.Get("name")
	query := "SELECT user, name, email, tel FROM member ;"
	arg := []any{}
	if name != "" {
		query += "WHERE name LIKE ? ;"
		arg = append(arg, "%"+name+"%")
	}

	UserResponse := []MemberResponse{}
	err := repo.db.Select(&UserResponse, query, arg...)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Database Error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserResponse)
}

func (repo repoDB) EditMember(w http.ResponseWriter, r *http.Request) {
	var userData Member
	err := json.NewDecoder(r.Body).Decode(&userData)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Data Not Formate"))
	}
	query := "UPDATE member set pass=?, name=?, email=?, tel=? WHERE user=?"
	arg := []any{}
	arg = append(arg, userData.Password)
	arg = append(arg, userData.Name)
	arg = append(arg, userData.Email)
	arg = append(arg, userData.Tel)
	arg = append(arg, userData.User)

	result, err := repo.db.Exec(query, arg...)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("database Error"))
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected < 1 {
		fmt.Println("Roweffact < 0")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("database Not Edit"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userData)

}

func (repo repoDB) GetMemberByUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user"]
	query := "SELECT user, name, email, tel FROM member WHERE user = ?;"
	UserResponse := MemberResponse{}
	err := repo.db.Get(&UserResponse, query, userID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Database Error"))
		return
	}
	if err == sql.ErrNoRows {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No User"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserResponse)
}

func (repo repoDB) DeleteMember(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user"]
	query := "DELETE FROM member WHERE user = ?"
	UserResponse := MemberResponse{}
	err := repo.db.Select(&UserResponse, query, userID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Delete Sucsess"))
		return
	}
	if err == sql.ErrNoRows {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Database Not Edit"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
