package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todo-app/list"
)

func printItems(l *list.ItemList) {
	if len(l.Items) == 0 {
		fmt.Println("No items found")
		return
	}
	fmt.Println("")
	fmt.Println("To-Do list:")
	fmt.Printf("%-5s %-12s %s\n", "ID", "STATUS", "Description")
	fmt.Println(strings.Repeat("-", 50))
	for _, item := range l.Items {
		fmt.Printf("%-5d %-12s %s\n", item.ID, item.Status, item.Description)
	}
	fmt.Println("")
}

func main() {
	items := list.NewItemList(list.DefaultDataFile)
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
	add <Descriptio>					- Add a new to-do item
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
			desc := strings.Join(args[1:], " ")
			items.Add(desc)
			fmt.Println("Item added")

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
			field := args[2]
			value := strings.Join(args[3:], " ")

			switch field {
			case "description":
				if items.UpdateDescription(id, value) {
					fmt.Println("Description updated")
				} else {
					fmt.Printf("Item with ID %d not found.\n", id)
				}

			case "status":
				if items.UpdateStatus(id, value) {
					fmt.Println("Status updated")
				} else {
					fmt.Println("Invalid status or ID")
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

			if items.Delete(id) {
				fmt.Println("Item Deleted")
			} else {
				fmt.Printf("Item with ID %d not found.\n", id)
			}

		case "exit":
			fmt.Println("Closing application...")
			return

		default:
			fmt.Println("Unknown command. Use 'help' for more information")

		}
	}
}
