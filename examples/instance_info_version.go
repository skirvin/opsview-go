package main

import (
	"github.com/skirvin/opsview-go/opsview"
	"log"
	"os"
)

func main() {
	log.Println("Get Opsview Instance Info")
	var instance = opsview.NewClient(os.Getenv("OPSVIEW_USERNAME"), os.Getenv("OPSVIEW_PASSWORD"), os.Getenv("OPSVIEW_BASEURI"), true)

	log.Printf("Opsview Instance Version: %s", instance.Version)
}