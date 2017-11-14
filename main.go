package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var (
	projectID = "cheerlights-hrd"
)

type Data struct {
	SourceIp string
	Color    string
	Date     time.Time
}

func handle(w http.ResponseWriter, r *http.Request) {
	defer func(t time.Time) { fmt.Fprintf(w, "Stored result in %s.", time.Since(t)) }(time.Now())
	ctx := appengine.NewContext(r)

	// create a safe data struct.
	data := Data{SourceIp: "127.0.0.1",
		Color: "white",
		Date:  time.Now().Round(0),
	}

	// replace the data struct's contents (ip/color) from Get vars.
	if color := r.URL.Query().Get("color"); color != "" {
		data.Color = color
	}
	if ra := r.RemoteAddr; ra != "" {
		data.SourceIp = ra
	}

	k := datastore.NewIncompleteKey(ctx, "Data", nil)
	if _, err := datastore.Put(ctx, k, &data); err != nil {
		log.Fatalf("failed to save data to store: %v", err)
	}

	fmt.Fprintf(w, "saved key/values to store. %v ", data)
}

func report(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// limit report to requestor's data only (by src ip)
	sip := r.RemoteAddr

	// Construct full sql query.
	query := datastore.NewQuery("Data").
		Filter("SourceIp=", sip).
		Order("Date").
		Limit(100)

	if admin := r.URL.Query().Get("admin"); admin != "" {
		query = datastore.NewQuery("Data").
			Order("Date").
			Limit(100)
	}

	var results []Data
	if _, err := query.GetAll(ctx, &results); err != nil {
		log.Fatalf("failed to retrieve data from the datastore: %v", err)
	}
	reportTmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		log.Fatalf("failed to parse the template: %v", err)
	}
	if err = reportTmpl.Execute(w, results); err != nil {
		log.Fatalf("failed to execute the template: %v", err)
	}
}

func main() {
	http.HandleFunc("/", handle)
	http.HandleFunc("/report", report)

	appengine.Main()
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
