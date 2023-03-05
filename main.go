package main

import (
	"encoding/json"
	"fmt"
	scl "hack-msb/SQLClient"
	"log"
	"net/http"
	"time"
)

func term(endp string) {
	log.Printf("Requested : <%s>\n", endp)
}

var sqlc *scl.SQLClient = scl.OpenConn()

func hello(w http.ResponseWriter, req *http.Request) {
	term("Hello")
	fmt.Fprintf(w, "hello\n")
}

type User struct {
	UserId, UserName string
	UserRole         uint
}

type PickingElement struct {
	Emplacement, Article string
	Quantite, Capacite   int
}

type StockElement struct {
	Emplacement, Article, Lot, Sous_lot, US string
	Quantite                                int
}

type StockRes struct {
	Produits []StockElement
}

type PickingRes struct {
	Produits []PickingElement
}

func getPicking(w http.ResponseWriter, req *http.Request) {
	term("Picking")

	row := sqlc.GetRows(
		fmt.Sprintf(
			"SELECT * from %s",
			scl.TBN_PICKING,
		),
	)
	defer row.Close()
	res := PickingRes{}

	for row.Next() {
		new := PickingElement{}
		row.Scan(&new.Emplacement, &new.Article, &new.Quantite, &new.Capacite)
		res.Produits = append(res.Produits, new)
	}

	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}

func getStock(w http.ResponseWriter, req *http.Request) {
	term("Stock")

	row := sqlc.GetRows(
		fmt.Sprintf(
			"SELECT * from %s",
			scl.TBN_STOCK,
		),
	)
	defer row.Close()
	res := StockRes{}

	for row.Next() {
		new := StockElement{}
		row.Scan(&new.Emplacement, &new.Article, &new.Lot, &new.Sous_lot, &new.Quantite, &new.US, nil)
		res.Produits = append(res.Produits, new)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}

const ALGORITHM string = "SELECT %s FROM EtatStock WHERE Article = '%s' and Statut = 'A' ORDER BY SUBSTRING(Lot, 7, 2), SUBSTRING(Lot, 9, 3), Quantite"

type Crate struct {
	Next string
}

func getNextCrate(w http.ResponseWriter, req *http.Request) {
	term("Next Crate")
	itemId := req.URL.Query().Get("itemId")

	row := sqlc.GetRows(
		fmt.Sprintf(
			ALGORITHM,
			"Emplacement",
			itemId,
		),
	)
	defer row.Close()
	row.Next()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := Crate{}
	row.Scan(&res.Next)
	json.NewEncoder(w).Encode(res)
}

type IncomingLog struct {
	Qte, State int
	Produit    string
	ProduitId  string
	From       string
}

type Log struct {
	Qte, State                   int
	Produit, Lot, Slot, StockPos string
	PrduitId, PickPos            string
	Time                         time.Time
}

func getLogs(w http.ResponseWriter, req *http.Request) {}
func logger(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var inc IncomingLog
	err := json.NewDecoder(req.Body).Decode(&inc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getUser(w http.ResponseWriter, req *http.Request) {
	term("User")
	userId := req.URL.Query().Get("userId")
	row := sqlc.GetRows(
		fmt.Sprintf(
			"SELECT * from %s where UserId = %s",
			scl.TBN_USERS,
			userId,
		),
	)
	defer row.Close()
	row.Next()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	user := User{}
	row.Scan(&user.UserId, &user.UserName, &user.UserRole)
	json.NewEncoder(w).Encode(user)
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/users", getUser)
	http.HandleFunc("/picking", getPicking)
	http.HandleFunc("/nextcrate", getNextCrate)
	http.HandleFunc("/stock", getStock)
	println("zebi")

	if err := http.ListenAndServe("localhost:5050", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
