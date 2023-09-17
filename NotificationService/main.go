package main

import (
	"NotificationService/internal/repo"
	"NotificationService/internal/route"
	"NotificationService/internal/utils"
	"context"
	"fmt"
	"net/http"

	_ "NotificationService/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@Title			Notification service by Ivan Sidorenko
//	@Version		1.0
//	@description	Server's API

//	@host		localhost:8181
//	@BasePath	/

func main() {
	utils.InitGenerator()
	utils.PMan = utils.NewPman("application.properties")
	repo.ConnectToMongo(context.TODO(),
		fmt.Sprintf("mongodb://%s:%s",
			utils.PMan.Get("mongo_host").(string),
			utils.PMan.Get("mongo_port").(string)),
		utils.PMan.Get("mongo_db_name").(string))
	repo.DropDb()
	router := mux.NewRouter()
	route.NewClientRouter(router)
	route.NewDistributionRouter(router)
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "admin/index.html")
	}).Methods(http.MethodGet)
	http.Handle("/", router)
	go route.HandleLoop()
	err := http.ListenAndServe(":8181", nil)
	panic(err)
}
