package thread

import (
	"log"
	"sync"
)

type ThreadManager struct {
	wg sync.WaitGroup
}

func NewThreadManager() *ThreadManager {
	return &ThreadManager{}
}

func (tm *ThreadManager) New(threadHandler func()) {
	tm.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v", err)
			}
			tm.wg.Done()
		}()
		threadHandler()
	}()
}

func (tm *ThreadManager) Wait() {
	tm.wg.Wait()
}
