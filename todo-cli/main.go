package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"todo-cli/api"
	"todo-cli/list"
)

func printItems(items []list.Item) {
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
	fmt.Println()
}

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	filename := list.DefaultDataFile
	items := list.LoadFromFile(filename)
	scanner := bufio.NewScanner(os.Stdin)

	slog.Info("Application Started")
	fmt.Println("Welcome to the To-Do List Application")
	fmt.Println("Type 'help' to see available commands")

	//Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		slog.Info("Graceful shutdown signal received - Closing Application...")
		stop()
		os.Exit(0)
	}()

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
	update <id> status <new status>				- Update item status (started, not started, completed)
	delete <id>						- Delete an item
	server							- Start HTTP Json API on port 8080
	exit							- Exit the application
		`)

		case "server":
			fmt.Println("Starting HTTP server on http://localhost:8080")
			go api.StartServer() //Starts the server
			fmt.Println("Press Crtl+c to stop the server gracefully.")
			select {} //keeps running until interrupted

		case "add":
			if len(args) < 2 {
				fmt.Println("Usage: add <description>")
				continue
			}
			desc := strings.Join(args[1:], " ")
			items = list.Add(items, desc)
			list.SaveToFile(filename, items)
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
				slog.Error("Invalid ID", "input", args[1], "error", err)
				continue
			}

			field := args[2]
			value := strings.ToLower(strings.Join(args[3:], " "))

			switch field {
			case "description":
				updatedItems, err := list.UpdateDescription(items, id, value)
				if err != nil {
					fmt.Println("Error: ", err)
					slog.Error("Update description failed", "id", id, "error", err)
					continue
				}
				items = updatedItems
				list.SaveToFile(filename, items)
				fmt.Println("Description updated")

			case "status":
				updateItems, err := list.UpdateStatus(items, id, value)
				if err != nil {
					fmt.Println("Error: ", err)
					slog.Error("Update status failed", "id", id, "error", err)
					continue
				}
				items = updateItems
				list.SaveToFile(list.DefaultDataFile, items)
				fmt.Println("Status updated")

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
				slog.Error("Invalid delete ID", "input", args[1], "error", err)
				continue
			}

			items = list.Delete(items, id)
			list.SaveToFile(filename, items)
			fmt.Println("Item deleted")

		case "exit":
			fmt.Println("Closing application...")
			return

		default:
			fmt.Println("Unknown command. Use 'help' for more information")
			slog.Warn("Unknown command", "command", command)

		}
	}
}
