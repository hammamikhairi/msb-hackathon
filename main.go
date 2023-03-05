package main

import (
	"encoding/json"
	"fmt"
	scl "hack-msb/SQLClient"
	"log"
	"net/http"
	"time"
	// "time"
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
	Picking []PickingElement
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
		res.Picking = append(res.Picking, new)
	}

	w.Header().Set("Content-Type", "application/json")
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

func nextCrate(prodId string) StockElement {
	row := sqlc.GetRows(
		fmt.Sprintf(
			ALGORITHM,
			" * ",
			prodId,
		),
	)
	defer row.Close()
	row.Next()

	res := StockElement{}
	row.Scan(&res.Emplacement, &res.Article, &res.Lot, &res.Sous_lot, &res.Quantite, &res.US, nil)
	return res
}

func getNextCrate(w http.ResponseWriter, req *http.Request) {
	term("Next Crate")
	itemId := req.URL.Query().Get("itemId")
	elem := nextCrate(itemId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := Crate{elem.Emplacement}
	json.NewEncoder(w).Encode(res)
}

type IncomingLog struct {
	Qte     int
	State   string
	Produit string
}

type Log struct {
	Article, Op, PosEnStock, Lot string
	Sous_lot                     string
	Qte                          int
	OpTime                       time.Time
}

func getCurrentTime() string {
	now := time.Now()
	return now.UTC().Format("2006-01-02 15:04:05")
}

type Logs struct {
	Logs []Log
}

func getLogs(w http.ResponseWriter, req *http.Request) {
	term("Piping Logs")
	rows := sqlc.GetRows(
		fmt.Sprintf(
			"select * from %s",
			scl.TBN_HISTORY,
		),
	)
	defer rows.Close()

	logs := Logs{}
	for rows.Next() {
		log := Log{}
		rows.Scan(&log.Article, &log.Op, &log.Qte, &log.PosEnStock, &log.Lot, &log.Sous_lot, &log.OpTime)
		logs.Logs = append(logs.Logs, log)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(logs)
}

func logger(w http.ResponseWriter, req *http.Request) {
	term("logger")
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

	next := nextCrate(inc.Produit)

	logRow := []any{
		inc.Produit,
		inc.State,
		inc.Qte,
		next.Emplacement,
		next.Lot,
		next.Sous_lot,
		fmt.Sprint(getCurrentTime()),
	}

	sqlc.PushRows([][]any{logRow}, scl.TBN_HISTORY)
	_ = logRow
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

	http.HandleFunc("/koudoular", KOUDOULAR)

	http.HandleFunc("/log", logger)
	http.HandleFunc("/logs", getLogs)

	if err := http.ListenAndServe("localhost:5051", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func KOUDOULAR(w http.ResponseWriter, req *http.Request) {
	term("Picking")

	row := sqlc.GetRows(
		fmt.Sprintf(
			"SELECT * from %s",
			scl.TBN_PICKING,
		),
	)
	defer row.Close()
	var res []PickingElement
	for row.Next() {
		new := PickingElement{}
		row.Scan(&new.Emplacement, &new.Article, &new.Quantite, &new.Capacite)
		res = append(res, new)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}
