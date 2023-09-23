package routes

import (
	"net/http"

	"github.com/Chufretalas/pantsbase/controllers"
	"github.com/gorilla/mux"
)

var Router *mux.Router

func LoadRoutes() {
	Router = mux.NewRouter()
	Router.HandleFunc("/", controllers.Index)
	Router.HandleFunc("/table_view", controllers.TableView)
	Router.HandleFunc("/new_table", controllers.NewTable).Methods("POST")
	Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))) // I love you gorilla mux ❤
	http.Handle("/", Router)
}
