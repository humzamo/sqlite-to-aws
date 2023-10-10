package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/mattn/go-sqlite3"

	"github.com/humzamo/sqlite-to-aws/internal/data"
)

const (
	awsRegion  = "eu-west-2"
	bucketName = "humza-mo-sqlite-to-aws"
	path       = "client-data"
	timeFormat = "2006-01-02-15:04:05.000"
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

	// Create an AWS S3 configuration and client
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		fmt.Println("Error loading AWS configuration:", err)
		return
	}
	client := s3.NewFromConfig(cfg)

	fileName := fmt.Sprintf("client-data-%s.json", time.Now().Format(timeFormat))

	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(path + "/" + fileName),
		Body:        aws.ReadSeekCloser(bytes.NewReader(jsonResult)),
		ContentType: aws.String("application/json"),
	}

	// upload the JSON data to S3.
	_, err = client.PutObject(context.Background(), uploadInput)
	if err != nil {
		fmt.Println("Error uploading JSON to S3:", err)
		return
	}

	fmt.Printf("JSON file %s uploaded to S3 successfully\n", fileName)
}
