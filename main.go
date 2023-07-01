package main

import (
	"STMBACKEND/handler"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func main() {
	ict, err := time.LoadLocation("Asia/Bangkok")

	if err != nil {
		panic(ict)
	}
	time.Local = ict
	db := InitDatabase()
	//fmt.Println("Arhtit")
	router := mux.NewRouter()
	member := handler.NewMemberHandler(db)
	login := handler.NewLoginHandler(db)

	router.HandleFunc("/members", member.NewMember).Methods(http.MethodPost)
	router.HandleFunc("/members", member.GetMember).Methods(http.MethodGet)
	router.HandleFunc("/members", member.EditMember).Methods(http.MethodPut)
	router.HandleFunc("/members/{user}", member.GetMemberByUser).Methods(http.MethodGet)
	router.HandleFunc("/members/{user}", member.DeleteMember).Methods(http.MethodDelete)

	router.HandleFunc("/logins", login.Login).Methods(http.MethodPost)

	fmt.Println("Server Start at Port : 500...")
	http.ListenAndServe(":5500", router)

}

func InitDatabase() *sqlx.DB {
	dbUser := "root"
	dbPassword := "art12345"
	dbHost := "localhost"
	dbName := "tsmbackend"
	db, err := sqlx.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v",
		dbUser,
		dbPassword,
		dbHost,
		dbName,
	))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(1 * time.Minute)
	db.SetMaxOpenConns(12)
	db.SetMaxIdleConns(12)
	return db
}
