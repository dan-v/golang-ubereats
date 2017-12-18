package main

import (
	"log"
	"flag"
	"html/template"
	"net/http"
	"github.com/dan-v/golang-ubereats/ubereats"
)

type PopularItems struct {
	StoreTitle string
	Items      ubereats.MenuItems
}

type AllStorePopularItems []*PopularItems

func (a *AllStorePopularItems) serve(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.ParseFiles("template.html"))
	templates.ExecuteTemplate(w, "template.html", a)
}

func main() {
	// parse flags
	latitudePtr := flag.String("lat", "47.6062", "Latitude")
	longitudePtr := flag.String("long", "-122.3321", "Longitude")
	portPtr := flag.String("port", "3000", "Port to start web server on")
	flag.Parse()

	// initialize api
	client, err := ubereats.NewClient(*latitudePtr, *longitudePtr)
	if err != nil {
		log.Fatal(err)
	}

	// get store list
	storeList, err := client.GetStoreList()
	if err != nil {
		log.Fatal(err)
	}

	// loop through all stores and get most popular items
	all := AllStorePopularItems{}
	concurrency := 25
	sem := make(chan bool, concurrency)
	for _, store := range storeList {
		log.Println(store.Title)
		sem <- true
		go func(store ubereats.Store) {
			defer func() { <-sem }()
			sd, err := client.GetStoreDetails(store.UUID)
			if err != nil {
				log.Println("Failed to get store details for " + store.Title, err.Error())
				return
			}
			popularItems := sd.GetPopularItems()
			if popularItems.String() != "" {
				all = append(all, &PopularItems{StoreTitle: store.Title, Items: popularItems})
			}
		}(store)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	http.HandleFunc("/", all.serve)
	log.Println("Web page at: http://localhost:" + *portPtr)
	http.ListenAndServe(":" + *portPtr, nil)
}
