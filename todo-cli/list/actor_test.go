package list

import (
	"fmt"
	"sync"
	"testing"
)

func TestListActor_ConcurrentAccess(t *testing.T) {
	actor := NewListActor([]Item{})
	defer actor.Stop()

	var wg sync.WaitGroup

	//run 10 concurrent writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_, err := actor.Add("Task " + string(rune('A'+n)))
			if err != nil {
				t.Errorf("Add failed %v", err)
			}
		}(i)
	}

	//Run 10 concurrent readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := actor.GetAll()
			if err != nil {
				t.Errorf("GetAll failed: %v", err)
			}
		}()
	}

	wg.Wait()

	items, _ := actor.GetAll()
	if len(items) == 0 {
		t.Fatal("Expected items after concurrent adds, got none")
	}
}

func TestListActor_Parallel(t *testing.T) {
	actor := NewListActor([]Item{})
	defer actor.Stop()

	t.Run("AddParallel", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			t.Run(fmt.Sprintf("Add-%d", i), func(t *testing.T) {
				t.Parallel()
				actor.Add("Concurrent Task")
			})
		}
	})

	t.Run("GetParallel", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 10; i++ {
			t.Run(fmt.Sprintf("Get-%d", i), func(t *testing.T) {
				t.Parallel()
				actor.GetAll()
			})
		}
	})
}
