package data

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const createClientsTableQuery = `
DROP TABLE IF EXISTS clients;
CREATE TABLE clients(id INTEGER PRIMARY KEY, name TEXT, revenue INT, ukbased BOOL);
INSERT INTO clients(name, revenue, ukbased) VALUES('Apple',1111,false);
INSERT INTO clients(name, revenue, ukbased) VALUES('BT',2222,true);
INSERT INTO clients(name, revenue, ukbased) VALUES('Costco',3333,false);
INSERT INTO clients(name, revenue, ukbased) VALUES('Delta Airlines',4444,false);
`

func CreateSampleTable() error {
	db, err := sql.Open("sqlite 3", "sample.db")
	if err != nil {
		return errors.Wrap(err, "unable to create new table sample.db")
	}

	fmt.Println("sample.db created")
	defer db.Close()

	_, err = db.Exec(createClientsTableQuery)
	if err != nil {
		return errors.Wrap(err, "unable to create new table clients")
	}
	fmt.Println("table clients created")

	return nil
}
