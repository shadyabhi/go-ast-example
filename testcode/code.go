package examplecode

import (
	"sync"

	log "github.com/Sirupsen/logrus"
)

func bob() {
	log.Infof("Hello World from Bob")
}

func main() {
	log.Infof("Hello World")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		bob()
		wg.Done()
	}()
	wg.Wait()
}
