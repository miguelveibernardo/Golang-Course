package list

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

const DefaultDataFile = "items.json"

const (
	StatusNotStarted = "not started"
	StatusStarted    = "started"
	StatusCompleted  = "completed"
)

// the to-do list structure
type Item struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func LoadFromFile(filename string) []Item {
	data, err := os.ReadFile(filename)
	if err != nil {
		slog.Warn("No existing data file found, starting with emplty list", "file", filename, "error", err)
		return []Item{}
	}

	var items []Item

	if err := json.Unmarshal(data, &items); err != nil {
		slog.Error("Error loading items", "file", filename, "error", err)
		return []Item{}
	}

	slog.Info("Items loaded from file", "file", filename, "count", len(items))
	return items
}

func GetNextID(items []Item) int {
	maxID := 0
	for _, i := range items {
		if i.ID >= maxID {
			maxID = i.ID + 1
		}
	}
	return maxID
}

// Saves items to JSON file
func SaveToFile(filename string, items []Item) {

	data, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		slog.Error("Error marshlling items for save", "file", filename, "error", err)
		return
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		slog.Error("Error writing items for file", "file", filename, "error", err)
		return
	}

	slog.Info("Items save successfully", "file", filename, "count", len(items))
}

func Add(items []Item, description string) []Item {

	newItem := Item{
		ID:          GetNextID(items),
		Description: description,
		Status:      StatusNotStarted,
	}
	slog.Info("Items added", "id", newItem.ID, "description", newItem.Description)
	return append(items, newItem)
}

func Delete(items []Item, id int) []Item {
	newItems := []Item{}
	found := false
	for _, item := range items {
		if item.ID != id {
			newItems = append(newItems, item)
		} else {
			found = true
			slog.Info("Item deleted", "id", id)
		}
	}

	if !found {
		slog.Warn("Attempted to delete non-existing item", "id", id)
	}

	return newItems
}

func UpdateDescription(items []Item, id int, desc string) ([]Item, error) {
	for i, item := range items {
		if item.ID == id {
			items[i].Description = desc
			slog.Info("Item description updated", "id", id, "new_description", desc)
			return items, nil
		}
	}
	slog.Warn("Update description failed: item not found", "id", id)
	return items, fmt.Errorf("item with ID %d not found", id)
}

func UpdateStatus(items []Item, id int, status string) ([]Item, error) {
	for i, item := range items {
		if item.ID == id {
			switch status {
			case StatusStarted, StatusCompleted, StatusNotStarted:
				items[i].Status = status
				slog.Info("Item status updated", "id", id, "new_status", status)
				return items, nil
			default:
				slog.Warn("Invalid status value", "status", status)
				return items, fmt.Errorf("invalid satus: %s. Please use Started, Not Started or Completed", status)
			}
		}
	}
	slog.Warn("Update status failed: item not found", "id", id)
	return items, fmt.Errorf("item with ID %d not found", id)
}
