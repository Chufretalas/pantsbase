package main

import (
	"fmt"
	"net/http"

	"github.com/Chufretalas/pantsbase/db"
	"github.com/Chufretalas/pantsbase/routes"
)

func main() {
	db.ConnectDB()
	defer db.DB.Close()

	routes.LoadRoutes()
	fmt.Println("Listening on http://localhost:8000")
	http.ListenAndServe(":8000", routes.Router)
}
