package list_test

import (
	"os"
	"testing"
	"todo-cli/list"
)

func sampleItems() []list.Item {
	return []list.Item{
		{ID: 1, Description: "Task 1", Status: list.StatusNotStarted},
		{ID: 2, Description: "Task 2", Status: list.StatusStarted},
	}
}

func TestAdd(t *testing.T) {
	items := sampleItems()
	items = list.Add(items, "New Task")

	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}

	last := items[len(items)-1]
	if last.Description != "New Task" {
		t.Errorf("Expected 'New Task', got %q", last.Description)
	}
	if last.Status != list.StatusNotStarted {
		t.Errorf("Expected status 'not started', got %q", last.Status)
	}
}

func TestDelete(t *testing.T) {
	items := sampleItems()
	items = list.Delete(items, 1)

	if len(items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(items))
	}

	if items[0].ID == 1 {
		t.Errorf("Item with ID 1 should have been deleted")
	}
}

func TestUpdateDescription(t *testing.T) {
	items := sampleItems()
	updated, err := list.UpdateDescription(items, 1, "Updated Task 1")

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if updated[0].Description != "Updated Task 1" {
		t.Errorf("Expected 'Updated Task 1', got %q", updated[0].Description)
	}
	//Invalid ID
	_, err = list.UpdateDescription(items, 999, "No Task")
	if err == nil {
		t.Errorf("Expected error for invalid ID, got nil")
	}
}

func TestUpdateStatus(t *testing.T) {
	items := sampleItems()
	updated, err := list.UpdateStatus(items, 2, list.StatusCompleted)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if updated[1].Status != list.StatusCompleted {
		t.Errorf("Expected 'Completed', got %q", updated[1].Status)
	}
	//Invalid Status
	_, err = list.UpdateStatus(items, 1, "Doing Something weird")
	if err == nil {
		t.Errorf("Expected error for invalid status, got nil")
	}

	_, err = list.UpdateStatus(items, 999, list.StatusStarted)
	if err == nil {
		t.Errorf("Expected error for invalid ID, got nil")
	}
}

func TestGetNextID(t *testing.T) {
	items := sampleItems()
	next := list.GetNextID(items)
	if next != 3 {
		t.Errorf("Expected next ID to be 3, got %d", next)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpFile := "test_items.json"
	defer os.Remove(tmpFile)

	items := sampleItems()
	list.SaveToFile(tmpFile, items)

	loaded := list.LoadFromFile(tmpFile)
	if len(loaded) != len(items) {
		t.Errorf("Expected %d items, got %d", len(items), len(loaded))
	}

	for i := range items {
		if items[i].Description != loaded[i].Description {
			t.Errorf("Mismatch in description: expected %q, got %q", items[i].Description, loaded[i].Description)
		}
	}
}
