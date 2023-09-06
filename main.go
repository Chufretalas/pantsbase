package main

import (
	"net/http"

	"github.com/Chufretalas/pantsbase/db"
	"github.com/Chufretalas/pantsbase/routes"
)

func main() {
	db.ConnectDB()
	defer db.DB.Close()

	routes.LoadRoutes()
	http.ListenAndServe(":8000", routes.Router)
}
