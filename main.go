package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
)

var (
	projectID = "cheerlights-hrd"
)

type Data struct {
	SourceIp string
	Color    string
	Date     time.Time
}

func index(w http.ResponseWriter, r *http.Request) {
	defer func(t time.Time) { fmt.Fprintf(w, "Stored result in %s.", time.Since(t)) }(time.Now())

	ctx := r.Context()
	// Create new datastore client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(w, "failed to create datastore client: %v", err)
		log.Fatalf("failed to create datastore client: %v", err)
	}

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

	k := datastore.IncompleteKey("Data", nil)
	if _, err := client.Put(ctx, k, &data); err != nil {
		fmt.Fprintf(w, "failed to save data to store: %v", err)
		log.Fatalf("failed to save data to store: %v", err)
	}

	fmt.Fprintf(w, "saved key/values to store. %v ", data)
}

func report(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Create new datastore client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(w, "failed to create datastore client: %v", err)
		log.Fatalf("failed to create datastore client: %v", err)
	}

	// limit report to requestor's data only (by src ip)
	sip := r.RemoteAddr

	// Construct full sql query, limit to src-ip, unless admin is set.
	query := datastore.NewQuery("Data").
		Filter("SourceIp=", sip).
		Order("Date").
		Limit(100)

	// Return full data, if admin is set.
	if admin := r.URL.Query().Get("admin"); admin != "" {
		query = datastore.NewQuery("Data").
			Order("Date").
			Limit(1000)
	}

	var results []Data
	if _, err := client.GetAll(ctx, query, &results); err != nil {
		fmt.Fprintf(w, "failed to retrieve data from the datastore: %v", err)
		log.Fatalf("failed to retrieve data from the datastore: %v", err)
	}
	reportTmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		fmt.Fprintf(w, "failed to parse the template: %v", err)
		log.Fatalf("failed to parse the template: %v", err)
	}
	if err = reportTmpl.Execute(w, results); err != nil {
		fmt.Fprintf(w, "failed to execute the template: %v", err)
		log.Fatalf("failed to execute the template: %v", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("using default port: %v", port)
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/report", report)

	log.Printf("Listening on port: %v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
