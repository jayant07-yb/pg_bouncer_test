package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func test_db() {

	psqlInfo := "host=localhost port=5400 user=user1 password=user1pass dbname=db1 binary_parameters=yes sslmode=disable"

	// Other connection configs are read from the standard environment variables:
	// PGSSLMODE, PGSSLROOTCERT, and so on.
	db, err := sql.Open("postgres", psqlInfo)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Define a prepared statement. You'd typically define the statement
	// elsewhere and save it for use in functions such as this one.
	stmt, err := db.Prepare("SELECT * FROM test_table WHERE test_name = ?")
	if err != nil {
		log.Fatal(err)
	}

	// Execute the prepared statement, passing in an id value for the
	// parameter whose placeholder is ?
	err2 := stmt.QueryRow("Prepared_Statement")
	if err2 != nil {
		fmt.Println(err2)
	}

}

func main() {

	test_db()
}
