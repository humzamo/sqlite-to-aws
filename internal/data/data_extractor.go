package data

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

type Database struct {
	client       *sql.DB
	DatabaseName string
	TableName    string
}

func NewDatabase(client *sql.DB, databaseName string, tableName string) Database {
	return Database{
		client:       client,
		DatabaseName: databaseName,
		TableName:    tableName,
	}
}

func (d *Database) GetColumnData() ([]ColumnData, error) {
	rows, err := d.client.Query("PRAGMA table_info(" + d.TableName + ")")
	if err != nil {
		return nil, errors.Wrap(err, "error querying table structure")
	}
	defer rows.Close()

	var columns []ColumnData
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
			return nil, errors.Wrap(err, "error scanning row")
		}

		column := ColumnData{
			Name: name,
			Type: datatype,
		}
		columns = append(columns, column)
	}

	return columns, nil
}

func (d *Database) GetClientData() ([]ClientData, error) {
	rows, err := d.client.Query(fmt.Sprintf("SELECT * FROM %s", d.TableName))
	if err != nil {
		return nil, errors.Wrap(err, "error querying table data")
	}
	defer rows.Close()

	var clientData []ClientData
	for rows.Next() {
		row := ClientData{}
		err = rows.Scan(&row.ID, &row.Name, &row.Revenue, &row.UkBased)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning row")
		}

		clientData = append(clientData, row)
	}

	return clientData, nil
}
