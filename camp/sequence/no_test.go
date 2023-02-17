package sequence

import (
	"fmt"
	"sync"
	"testing"
)

func TestUserNO(t *testing.T) {
	i := [50]int{}
	wait := sync.WaitGroup{}
	for _ = range i {
		wait.Add(1)
		go func() {
			no, _ := ServicerNo.NewNo()
			fmt.Println("no", no)
			wait.Done()
		}()
	}
	wait.Wait()
}
