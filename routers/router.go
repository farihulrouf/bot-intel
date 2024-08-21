package router

import (
	"bot_intel/controllers" // Menggunakan nama module dan path yang sesuai

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/system/devices", controllers.CreateDevice).Methods("GET")
	return r
}
