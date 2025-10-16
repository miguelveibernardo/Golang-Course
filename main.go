package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const dataFile = "items.json"

const (
	StatusNotStarted = "not started"
	StatusStarted    = "started"
	StatusCompleted  = "completed"
)

// the to-do list
type Item struct {
	ID          int    `json: "id"`
	Description string `json: "description"`
	Status      string `json: "status"`
}

type ItemList []Item

var nextID int //ID tracker

func loadItems() ItemList {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		return ItemList{}
	}

	var items ItemList
	if err := json.Unmarshal(file, &items); err != nil {
		fmt.Println("Error loading the items: ", err)
		return ItemList{}
	}

	for _, item := range items {
		if item.ID >= nextID {
			nextID = item.ID + 1
		}
	}

	return items
}

func saveItems(items ItemList) {
	data, err := json.MarshalIndent(items, "", " ")
	if err != nil {
		fmt.Println("Error saving items: ", err)
		return
	}

	if err := os.WriteFile(dataFile, data, 0644); err != nil {
		fmt.Println("Error writing to file: ", err)
	}
}

func printItems(items ItemList) {
	if len(items) == 0 {
		fmt.Println("No items found")
		return
	}
	fmt.Println("")
	fmt.Println("To-Do list:")
	fmt.Printf("%-5s %-12s %s\n", "ID", "STATUS", "Description")
	fmt.Println(strings.Repeat("-", 50))
	for _, item := range items {
		fmt.Printf("%-5d %-12s %s\n", item.ID, item.Status, item.Description)
	}
	fmt.Println("")
}

func findItemByID(items ItemList, id int) (*Item, int) {
	for index, item := range items {
		if item.ID == id {
			return &items[index], index
		}
	}
	return nil, -1
}

func main() {
	items := loadItems()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to the To-Do List Application")
	fmt.Println("Type 'help' to see available commands")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		args := strings.Fields(input)
		command := strings.ToLower(args[0])

		switch command {
		case "help":
			fmt.Println(`
Available Commands:
	add <Descriptio>					- Add a new to-do item with its desciption
	list							- Shows the entire list
	update <id> description <new descriptio>		- Update item description
	update <id> status <new status>				- Update item status
	delete <id>						- Delete an item
	exit							- Exit the application
		`)

		case "add":
			if len(args) < 2 {
				fmt.Println("Usage: add <description>")
				continue
			}
			description := strings.Join(args[1:], " ")
			item := Item{
				ID:          nextID,
				Description: description,
				Status:      StatusNotStarted,
			}
			nextID++
			items = append(items, item)
			fmt.Println("Item added")
			saveItems(items)

		case "list":
			printItems(items)

		case "update":
			if len(args) < 4 {
				fmt.Println("Usage: update <id> description|status <value>")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid ID")
				continue
			}
			item, _ := findItemByID(items, id)
			if item == nil {
				fmt.Printf("Item with ID %d not found.\n", id)
				continue
			}
			field := args[2]
			value := strings.Join(args[3:], " ")

			switch field {
			case "description":
				item.Description = value
				fmt.Println("Description updated")
				saveItems(items)
			case "status":
				switch value {
				case StatusNotStarted, StatusStarted, StatusCompleted:
					item.Status = value
					fmt.Println("Status updated")
					saveItems(items)
				default:
					fmt.Println("Invalid status. Please use one of the following options: `not started`, `started`, or `completed`")
				}
			default:
				fmt.Println("Invalid update field. Please use one of the following options: `description`, or `status`")
			}

		case "delete":
			if len(args) != 2 {
				fmt.Println("Usage: delete <id>")
				continue
			}
			id, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("Invalid ID")
				continue
			}

			_, index := findItemByID(items, id)
			if index == -1 {
				fmt.Println("Item with ID %d not found.\n", id)
				continue
			}
			items = append(items[:index], items[index+1:]...)
			fmt.Println("Item deleted")
			saveItems(items)

		case "exit":
			saveItems(items)
			fmt.Println("Closing application...")

		default:
			fmt.Println("Unknown command. Use 'help' for more information")

		}
	}
}
