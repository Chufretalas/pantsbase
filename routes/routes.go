package routes

import (
	"io/fs"
	"net/http"

	"github.com/Chufretalas/pantsbase/controllers"
	"github.com/gorilla/mux"
)

var StaticFS fs.FS
var Router *mux.Router

func LoadRoutes() {
	Router = mux.NewRouter()
	Router.HandleFunc("/", controllers.Index)
	Router.HandleFunc("/table_view", controllers.TableView)
	Router.HandleFunc("/form_handlers/new_table", controllers.NewTable).Methods("POST")
	Router.HandleFunc("/form_handlers/new_row", controllers.NewRow).Methods("POST")
	Router.HandleFunc("/form_handlers/update_row", controllers.UpdateRow).Methods("POST")
	Router.HandleFunc("/api/tables", controllers.Tables).Methods("GET")
	Router.HandleFunc("/api/query/{table_name}", controllers.Query).Methods("POST")
	Router.HandleFunc("/api/delete_one/{table_name}/{id}", controllers.DeleteOne).Methods("DELETE")
	Router.HandleFunc("/api/delete_table/{table_name}", controllers.DeleteTable).Methods("DELETE")
	Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(StaticFS)))) // I love you gorilla mux ❤
	http.Handle("/", Router)
}
