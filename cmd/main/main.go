package main

import (
	"fmt"

	"github.com/Kaushik1766/LibraryManagement/internal/app"
	"github.com/Kaushik1766/LibraryManagement/internal/db"
)

func main() {
	dbCon := db.GetDB()

	err := db.CreateTables(dbCon)
	if err != nil {
		fmt.Println(err.Error())
	}

	App := app.NewApp(dbCon)
	App.Run()
}
