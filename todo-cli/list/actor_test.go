package list

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListActor_ConcurrentAccess(t *testing.T) {
	actor := NewListActor([]Item{})
	defer actor.Stop()

	var wg sync.WaitGroup

	//run concurrent writers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			desc := "Task " + string(rune('A'+n))
			items, err := actor.Add(desc)
			assert.NoError(t, err, "Add() should not return error")
			assert.GreaterOrEqual(t, len(items), 1, "Items should not be empty after Add()")
		}(i)
	}

	//Run concurrent readers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			items, err := actor.GetAll()
			assert.NoError(t, err, "GetAll() should not return error")
			assert.NotNil(t, items, "GetAll() should return a non-nil slice")
		}()
	}

	wg.Wait()

	items, err := actor.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, items, "Expected items after concurrent adds")
}

func TestListActor_Parallel(t *testing.T) {

	//create one actor for both subtests
	actor := NewListActor([]Item{})

	//ensure it only stops *after" all subtests finish
	t.Cleanup(func() {
		actor.Stop()
	})

	t.Run("AddParallel", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 100; i++ {
			_, err := actor.Add(fmt.Sprintf("Concurrent Task %d", i))
			if err != nil && err != ErrActorStopped {
				t.Errorf("unexpected error: %v", err)
			}
		}
		items, err := actor.GetAll()
		assert.NoError(t, err)
		assert.NotEmpty(t, items, "Expected items after AddParallel")
	})

	t.Run("GetParallel", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 10; i++ {
			items, err := actor.GetAll()
			if err != nil && err != ErrActorStopped {
				t.Errorf("unexpected error: %v", err)
			}
			if items != nil {
				assert.NotNil(t, items)
			}
		}
	})
}

// Tests graceful shutdown behavior
func TestListActor_Stop(t *testing.T) {
	actor := NewListActor([]Item{})
	_, err := actor.Add("Initial Task")
	assert.NoError(t, err)

	actor.Stop()

	_, err = actor.Add("Should fail after stop")
	assert.ErrorIs(t, err, ErrActorStopped, "Ading after stop should return ErrActorStopped")
}

func TestListActor_CommandSequence(t *testing.T) {
	actor := NewListActor([]Item{})
	defer actor.Stop()

	//Adding multiple(5) items
	for i := 0; i < 5; i++ {
		_, err := actor.Add("Item " + string(rune('A'+i)))
		assert.NoError(t, err)
	}

	//Update status and description sequentially
	items, err := actor.UpdateStatus(2, "completed")
	assert.NoError(t, err)
	assert.Equal(t, "completed", items[2].Status)

	items, err = actor.UpdateDescription(3, "Updated C")
	assert.NoError(t, err)
	assert.Equal(t, "Updated C", items[3].Description)

}

func TestListActor_Race(t *testing.T) {
	actor := NewListActor([]Item{})
	defer actor.Stop()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(n int) {
			defer wg.Done()
			actor.Add("Race Task " + string(rune('A'+n%26)))
		}(i)
		go func() {
			defer wg.Done()
			actor.GetAll()
		}()
	}
	wg.Wait()

	time.Sleep(50 * time.Millisecond) //to let go routines to settle
	items, err := actor.GetAll()
	assert.NoError(t, err)
	assert.True(t, len(items) > 0)
}
