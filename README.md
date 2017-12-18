An unofficial Golang client for [Uber Eats](https://www.ubereats.com/) API. I've used this client to create an example application that grabs the most popular items from all restaurants in a given location.

## Most Popular Items Example
This will grab all restaurant popular items and display them on a beautifully styled web page at http://localhost:3000
```
go run main.go -lat "47.6062" -long "-122.3321"
```

![](/screenshot.png?raw=true)

## Unofficial Golang Client
```
import "github.com/dan-v/golang-ubereats/ubereats"

// given latitude and longitude of location
client, err := ubereats.NewClient("47.6062", "-122.3321")

// get store list
storeList, err := client.GetStoreList()

// get store details
for _, store := range storeList {
    fmt.Println(store.Title)
    details, err := client.GetStoreDetails(store.UUID)
    fmt.Println(details)
}
```