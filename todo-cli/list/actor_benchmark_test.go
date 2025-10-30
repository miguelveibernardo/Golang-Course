package list

import (
	"fmt"
	"testing"
)

// Benchmark adding items sequencially
func BenchmarkActorAdd(b *testing.B) {
	actor := NewListActor([]Item{})
	defer actor.Stop()

	for i := 0; i < b.N; i++ {
		_, err := actor.Add(fmt.Sprintf("Benchmark Task %d", i))
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// Benchmark adding items sequencially
func BenchmarkActorConcurrentAddGet(b *testing.B) {
	actor := NewListActor([]Item{})
	defer actor.Stop()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				actor.Add(fmt.Sprintf("Parallel Task %d", i))
			} else {
				actor.GetAll()
			}
			i++
		}
	})
}

// benchmark full lifecycle: Add - Update - Get - Delete
func BenchmarkActorFullLifecycle(b *testing.B) {
	actor := NewListActor([]Item{})
	defer actor.Stop()

	for i := 0; i < b.N; i++ {
		actor.Add(fmt.Sprintf("Task %d", i))
		actor.UpdateStatus(i%5, "completed")
		actor.UpdateDescription(i%5, "Updated Task")
		actor.Delete(i % 3)
		actor.GetAll()
	}
}
