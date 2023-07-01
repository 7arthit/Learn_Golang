package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
)

type Credentials struct {
	User     string `db:"user" json:"user"`
	Password string `db:"pass" json:"pass"`
}

type loginHandler interface {
	Login(http.ResponseWriter, *http.Request)
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func NewLoginHandler(db *sqlx.DB) loginHandler {
	return repoDB{db: db}
}

func (repo repoDB) Login(w http.ResponseWriter, r *http.Request) {
	// รับ Credentials
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No USER Credentials!"))
		return
	}
	// ตรวจสอบผู้ใช้
	query := "SELECT * FROM member WHERE user = ? AND pass = ? ;"
	result := Member{}
	err = repo.db.Get(&result, query, creds.User, creds.Password)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("USER or PASSWORD Invalid!"))
		return
	}

	// สร้าง JWT claims
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := Claims{
		Username: creds.User,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "suksan.group",
		},
	}

	// สร้าง JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var signKey string = "TSM2023"

	// ลงลายมือชื่อ Server ด้วย signKey
	signedToken, err := token.SignedString([]byte(signKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	type tokenResponse struct {
		Token       string `json: "token"`
		ExpiresDate string `json:"expire_date"`
		Etc         string `json:"etc"`
	}
	response := tokenResponse{
		Token:       signedToken,
		ExpiresDate: expirationTime.String(),
		Etc:         "You can access resource By use Authorization:Bearer HEADER",
	}
	// ส่ง Token ให้ Client
	// http.SetCookie(w, &http.Cookie{
	//  Name:    "token",
	//  Value:   signedToken,
	//  Expires: expirationTime,
	// })
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func validateAuthen(w http.ResponseWriter, r *http.Request) (bool, jwt.Claims) {
	t := r.Header.Get("Authorization")
	var signKey string = "TSM2023"

	if t == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return false, nil
	}
	jwtText := strings.TrimPrefix(t, "Bearer ")
	customClaims := Claims{}
	token, err := jwt.ParseWithClaims(
		jwtText,
		&customClaims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(signKey), nil
		},
	)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))

		return false, nil
	}
	claims := token.Claims
	return true, claims
}
