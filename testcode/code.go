package examplecode

import (
	"sync"

	log "github.com/Sirupsen/logrus"
)

func bob() {
	log.Debugf("Hello World from Bob: %d", 0)
}

func main() {
	log.Infof("Hello World from main")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		bob()
		wg.Done()
	}()
	wg.Wait()
}
