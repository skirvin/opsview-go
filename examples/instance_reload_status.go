package main

import (
	"github.com/skirvin/opsview-go/opsview"
	"log"
	"os"
)

func main() {
	log.Println("Get Opsview Instance Reload State")
	var instance = opsview.NewClient(os.Getenv("OPSVIEW_USERNAME"), os.Getenv("OPSVIEW_PASSWORD"), os.Getenv("OPSVIEW_BASEURI"), true)

	log.Println(instance.ReloadStatus)
}