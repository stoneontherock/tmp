package main

import (
	"aa/db"
	"aa/httpapi"
)

func main() {
	db.InitDataBase()
	httpapi.InitModel(db.DB)
	httpapi.Serve()
}
