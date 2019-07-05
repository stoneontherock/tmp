package main

import (
	"aa/aalog"
	"aa/config"
	"aa/db"
	"aa/grpc"
	"aa/httpapi"
)

func main() {
	config.GetConf()
	aalog.InitLog(&config.C.LOG)
	dbp := db.InitDataBase(&config.C.DB)
	httpapi.Init(dbp)
	go httpapi.Serve()
	grpc.Serve()
}
