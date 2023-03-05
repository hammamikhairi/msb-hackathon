// Go connection Sample Code:
package SQLClient

import (
	// "context"
	"database/sql"
	"strings"
	// "errors"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"

	"encoding/csv"
	"os"
)

var (
	server   = "stal-server.database.windows.net"
	port     = 1433
	user     = "khairi"
	password = "@Azure2023"
	database = "stal"
)

type SQLClient struct {
	db *sql.DB
}

const (
	TBN_PICKING string = "Picking"
	TBN_STOCK   string = "EtatStock"
	TBN_USERS   string = "Users"
	TBN_HISTORY string = "History"

	TB_PICKING string = "Picking (Emplacement, Article, Quantite, Capicite) "
	TB_STOCK   string = "EtatStock (Emplacement, Article, Lot, Sous_lot, Quantite, US, Statut) "
	TB_USERS   string = "Users (UserId, UserName, UserRole) "
	TB_HISTORY string = "History (Article, Op, Qte, PosEnStock, Lot, Sous_lot) "
)

func mainExamp() {
	sqlc := OpenConn()
	sqlc.Push("TEST")
}

func mainPushRowsTest() {

	sqlc := OpenConn()

	rows, err := LoadCSVData("test.csv")
	if err != nil {
		panic(err)
	}
	sqlc.PushRows(rows, TB_USERS)

}

func LoadCSVData(filename string) ([][]any, error) {
	// Open the CSV filename
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Convert the slice of slices of strings to a slice of slices of interface{}
	var result [][]interface{}
	for index, record := range records {
		if index == 0 {
			continue
		}
		row := make([]any, 0, len(record))
		for _, value := range record {
			row = append(row, value)
		}
		result = append(result, row)
	}

	return result, nil
}

func OpenConn() *SQLClient {

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	var err error
	// Create connection pool
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connected!")
	return &SQLClient{
		db,
	}
}

func (sc *SQLClient) GetRows(query string) *sql.Rows {
	rows, err := sc.db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	return rows
}

func fmtValues(values []any) string {
	base := "("
	for _, val := range values {
		base = base + ", '" + fmt.Sprint(val) + "'"
	}
	base += ")"
	return strings.Replace(base, ",", " ", 1)
}

func (sc *SQLClient) PushRows(rows [][]any, table string) {
	query := fmt.Sprintf("INSERT INTO %s VALUES ", table)
	for _, row := range rows {
		query += fmtValues(row) + ","

	}
	query = query[:len(query)-1]
	// fmt.Println(query)
	sc.Push(query)
}

func (sc *SQLClient) Push(query string) {
	_, err := sc.db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}
