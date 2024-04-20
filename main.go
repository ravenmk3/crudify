package main

import (
	"log"

	"crudify/db/mysql"
)

func main() {
	provider, err := mysql.NewMySqlSchemaProvider("misaka", 3306, "root", "mysql")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = provider.Close()
	}()

	tables, err := provider.GetTables("canal_manager")
	if err != nil {
		log.Fatal(err)
	}

	println(tables)
}
