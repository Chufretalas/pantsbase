package routes

import (
	"io/fs"
	"net/http"

	"github.com/Chufretalas/pantsbase/controllers"
)

var StaticFS fs.FS
var Router *http.ServeMux

func LoadRoutes(router *http.ServeMux) {
	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(StaticFS))))

	router.HandleFunc("GET /", controllers.Index)
	router.HandleFunc("GET /table_view", controllers.TableView)
	router.HandleFunc("POST /form_handlers/new_table", controllers.NewTable)

	router.HandleFunc("GET /api/tables", controllers.Tables)

	router.HandleFunc("GET /api/tables/{table_name}", controllers.Query)
	router.HandleFunc("GET /api/tables/{table_name}/{id}", controllers.QueryOne)
	router.HandleFunc("POST /api/tables/{table_name}", controllers.PostRows)
	router.HandleFunc("DELETE /api/tables/{table_name}/{id}", controllers.DeleteOne)
	router.HandleFunc("PUT /api/tables/{table_name}/{id}", controllers.UpdateOne)
	router.HandleFunc("PATCH /api/tables/{table_name}/{id}", controllers.UpdateOne)
	router.HandleFunc("DELETE /api/delete_table/{table_name}", controllers.DeleteTable)
}
