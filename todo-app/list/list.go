package list

import (
	"encoding/json"
	"fmt"
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

type ItemList struct {
	Items  []Item
	nextID int //ID tracker
	File   string
}

// Creates a new ItemList and loads data from file
func NewItemList(filename string) *ItemList {
	l := &ItemList{File: filename}
	l.Load()
	return l
}

// Reads the items from the JSON file
func (l *ItemList) Load() {
	data, err := os.ReadFile(l.File)
	if err != nil {
		l.Items = []Item{}
		return
	}

	if err := json.Unmarshal(data, &l.Items); err != nil {
		fmt.Println("Error loading the items: ", err)
		l.Items = []Item{}
		return
	}

	l.nextID = 0

	for _, item := range l.Items {
		if item.ID >= l.nextID {
			l.nextID = item.ID + 1
		}
	}
}

// Saves items to JSON file
func (l *ItemList) Save() {

	data, err := json.MarshalIndent(l.Items, "", " ")
	if err != nil {
		fmt.Println("Error saving items: ", err)
		return
	}

	if err := os.WriteFile(l.File, data, 0644); err != nil {
		fmt.Println("Error writing to file: ", err)
	}
}

func (l *ItemList) Add(description string) {

	item := Item{
		ID:          l.nextID,
		Description: description,
		Status:      StatusNotStarted,
	}
	l.nextID++
	l.Items = append(l.Items, item)
	l.Save()
}

func (l *ItemList) FindByID(id int) (*Item, int) {
	for i, item := range l.Items {
		if item.ID == id {
			return &l.Items[i], i
		}
	}
	return nil, -1
}

func (l *ItemList) Delete(id int) bool {
	_, index := l.FindByID(id)
	if index == -1 {
		return false
	}
	l.Items = append(l.Items[:index], l.Items[index+1:]...)
	l.Save()
	return true
}

func (l *ItemList) UpdateDescription(id int, desc string) bool {
	item, _ := l.FindByID(id)
	if item == nil {
		return false
	}
	item.Description = desc
	l.Save()
	return true
}

func (l *ItemList) UpdateStatus(id int, status string) bool {
	item, _ := l.FindByID(id)
	if item == nil {
		return false
	}
	switch status {
	case StatusNotStarted, StatusStarted, StatusCompleted:
		item.Status = status
		l.Save()
		return true
	default:
		return false
	}
}
