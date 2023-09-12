package main

import (
	"NotificationService/internal/repo"
	"NotificationService/internal/route"
	"NotificationService/internal/utils"
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

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
	http.Handle("/", router)
	go route.HandleLoop()
	err := http.ListenAndServe(":8181", nil)
	panic(err)
}
