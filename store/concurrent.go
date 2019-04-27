package store

import (
	"log"
	"sync"
)

// ConcurrentSave is a simple wrapper around DataStoreWrapper#Save. It is meant to ease calling from a `go` statement. The WaitGroup must be Added to before calling this function.
func ConcurrentSave(n map[string]int, dataStore *DataStoreWrapper, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Starting save")

	err := dataStore.Save(n)
	if err != nil {
		log.Printf("Save failed with error: %v\n", err)
		log.Println("Operations were rolled back")
	} else {
		log.Println("Save succeded")
	}
}
