package main

import (
	"aa/aalog"
	"aa/config"
	"aa/db"
	"aa/httpapi"
)

func main() {
	c := config.GetConf()
	aalog.InitLog(&c.LOG)
	db.InitDataBase(&c.DB)
	httpapi.InitModel(db.DB)
	httpapi.Serve(c.HTTP.ListenAddr)
}
