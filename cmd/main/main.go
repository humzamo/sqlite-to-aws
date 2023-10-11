package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	awsupload "github.com/humzamo/sqlite-to-aws/internal/aws-upload"
	"github.com/humzamo/sqlite-to-aws/internal/data"
)

const (
	awsRegion  = "eu-west-2"
	bucketName = "humza-mo-sqlite-to-aws"
	path       = "client-data"
	timeFormat = "2006-01-02-15:04:05.000"
)

func main() {
	// check if the command-line arguments for the database and table names are provided
	if len(os.Args) < 3 {
		fmt.Println("missing database name and/or table name")
		fmt.Println("Usage: go run main.go <database_name> <table_name>")
		os.Exit(1)
	}

	databaseName := os.Args[1]
	tableName := os.Args[2]

	db, err := sql.Open("sqlite3", databaseName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("successfully connected to database", databaseName)

	database := data.NewDatabase(db, databaseName, tableName)

	columns, err := database.GetColumnData()
	if err != nil {
		log.Fatalf("error getting column data from database %s and table %s: %v", database.DatabaseName, database.TableName, err)
	}
	fmt.Println("successfully retrieved column data for table", database.TableName)

	clientData, err := database.GetClientData()
	if err != nil {
		log.Fatalf("error getting column data from database %s and table %s: %v", database.DatabaseName, database.TableName, err)
	}
	fmt.Println("successfully retrieved table data for table", database.TableName)

	tableData := data.TableData{
		ColumnData: columns,
		ClientData: clientData,
	}

	jsonResult, err := json.Marshal(tableData)
	if err != nil {
		log.Fatalf("error marshalling json: %v", err)
	}

	fmt.Println("Extracted data from table:")
	fmt.Println(string(jsonResult))

	fileName := fmt.Sprintf("client-data-%s.json", time.Now().Format(timeFormat))

	err = awsupload.UploadToS3Bucket(fileName, jsonResult)
	if err != nil {
		log.Fatalf("error uploading file %s to S3: %v", fileName, err)
	}

	fmt.Printf("JSON file %s uploaded to S3 successfully\n", fileName)
}
