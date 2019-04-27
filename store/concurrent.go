package store

import (
	"log"
	"sync"
)

// ConcurrentSave is a simple wrapper around DataStoreWrapper#Save. It is meant to ease calling from a `go` statement. The WaitGroup must be Added to before calling this function.
func ConcurrentSave(n map[string]int, dataStore *DataStoreWrapper, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Starting save")

	log.Println("Clearing targets")
	errorChannel := make(chan error)
	quitChannel := make(chan interface{}, 1)

	go func(e <-chan error, q <-chan interface{}) {
		for {
			select {
			case err := <-e:
				if err != nil {
					log.Printf("Error deleting element: %v", err)
				}
			case <-q:
				return
			}
		}
	}(errorChannel, quitChannel)

	wwg := new(sync.WaitGroup)
	for name := range n {
		wwg.Add(1) // TODO: extract this into something like `wg.Add(len(n))`
		go func(elemName string) {
			defer wwg.Done()
			if err := dataStore.Delete(elemName); err != nil {
				errorChannel <- err
			}
		}(name)
	}
	wwg.Wait()
	quitChannel <- struct{}{}
	close(errorChannel)
	close(quitChannel)
	log.Println("Delete succeded")

	err := dataStore.Save(n)
	if err != nil {
		log.Printf("Save failed with error: %v\n", err)
		log.Println("Operations were rolled back")
	} else {
		log.Println("Save succeded")
	}
}
