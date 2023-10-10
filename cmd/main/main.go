package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/humzamo/sqlite-to-aws/internal/data"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// err := dummydata.CreateSampleTable()
	// if err != nil {
	// 	log.Fatalf("CreateSampleTable: %v", err)
	// }

	db, err := sql.Open("sqlite3", "sample.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("PRAGMA table_info(" + "clients" + ")")
	if err != nil {
		log.Fatalf("error querying table structure: %v", err)
	}
	defer rows.Close()

	var columns []data.ColumnData
	for rows.Next() {
		var (
			cid         int
			name        string
			datatype    string
			notnull     int
			dfltValue   sql.NullString
			primary_key int
		)
		err := rows.Scan(&cid, &name, &datatype, &notnull, &dfltValue, &primary_key)
		if err != nil {
			log.Fatalf("error scanning row: %v", err)
		}

		column := data.ColumnData{
			Name: name,
			Type: datatype,
		}
		columns = append(columns, column)
	}

	rows, err = db.Query("SELECT * FROM clients")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var clientData []data.ClientData
	for rows.Next() {
		row := data.ClientData{}
		err = rows.Scan(&row.ID, &row.Name, &row.Revenue, &row.UkBased)
		if err != nil {
			log.Fatal(err)
		}

		clientData = append(clientData, row)
	}

	tableData := data.TableData{
		ColumnData: columns,
		ClientData: clientData,
	}

	jsonResult, err := json.Marshal(tableData)
	if err != nil {
		log.Fatalf("error marshalling json: %v", err)
	}

	fmt.Println(string(jsonResult))
}
