package main

import (
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/appengine"
)

type Data struct {
	SourceIp string
	Color    string
	Date     time.Time
}

func handle(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	data := &Data{SourceIp: "127.0.0.1",
		Color: "white",
		Date:  time.Now(),
	}

	k := datastore.NameKey("Data", data.Date, nil)

}

func main() {
	http.HandleFunc("/", handle)
	http.HandleFunc("/report", report)

	appengine.Main()
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
