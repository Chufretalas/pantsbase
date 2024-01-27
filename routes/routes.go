package routes

import (
	"net/http"

	"github.com/Chufretalas/pantsbase/controllers"
	"github.com/gorilla/mux"
)

var Router *mux.Router

// TODO: make the api more consistent, like making delete_one be delete_one/{table_name}/{id}
func LoadRoutes() {
	Router = mux.NewRouter()
	Router.HandleFunc("/", controllers.Index)
	Router.HandleFunc("/table_view", controllers.TableView)
	Router.HandleFunc("/new_table", controllers.NewTable).Methods("POST")
	Router.HandleFunc("/new_row", controllers.NewRow).Methods("POST")
	Router.HandleFunc("/query", controllers.Query).Methods("POST")
	Router.HandleFunc("/delete_one", controllers.DeleteOne).Methods("DELETE")
	Router.HandleFunc("/update_row", controllers.UpdateRow).Methods("POST")
	Router.HandleFunc("/delete_table/{table_name}", controllers.DeleteTable).Methods("DELETE")
	Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))) // I love you gorilla mux ❤
	http.Handle("/", Router)
}
