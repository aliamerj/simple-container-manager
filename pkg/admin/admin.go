package admin

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"containerManger.com/pkg/core"
)

func AdminInterface() {
	fmt.Println("Welcome to Simple container management ")
	for {
		fmt.Println("_______________________________________________")
		fmt.Println("Container Management System - Admin Interface ")
		fmt.Println("Commands:")
		fmt.Println("1. List Containers")
		fmt.Println("2. Stop Container")
		fmt.Println("3. Continue a paused Container")
		fmt.Println("4. Remove Container ")
		fmt.Println("5. Remove All container")
		fmt.Println("6. Exit")
		fmt.Println("_______________________________________________")
		fmt.Print("Enter your choice: ")
		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = choice[:len(choice)-1] // remove trailing newline

		switch choice {
		case "1":
			fmt.Println("Listing all containers...")
			// Call function to list containers
			core.ListContainers()

		case "2":
			exitStopInterface := false
			for !exitStopInterface {
				// Get container ID from user
				fmt.Print("Enter the ID of the container you wish to stop or 'q' to quit: ")
				reader := bufio.NewReader(os.Stdin)
				idStr, _ := reader.ReadString('\n')
				idStr = strings.TrimSpace(idStr) // Remove trailing newline and any spaces

				// Check if user wants to quit
				if idStr == "q" {
					fmt.Println("Returning to the admin controller.")
					exitStopInterface = true
					continue
				}

				// Convert string ID to integer
				id, err := strconv.Atoi(idStr)
				if err != nil {
					fmt.Println("Invalid ID. Please enter a numeric container ID.")
					continue
				}

				// Call function to stop container
				core.StopContainer(id)
			}

		case "3":
			exitContinueInterface := false
			for !exitContinueInterface {
				// Get container ID from user
				fmt.Print("Enter the ID of the container you wish to restart or 'q' to quit: ")
				reader := bufio.NewReader(os.Stdin)
				idStr, _ := reader.ReadString('\n')
				idStr = strings.TrimSpace(idStr) // Remove trailing newline and any spaces

				// Check if user wants to quit
				if idStr == "q" {
					fmt.Println("Returning to the admin controller.")
					exitContinueInterface = true
					continue
				}

				// Convert string ID to integer
				id, err := strconv.Atoi(idStr)
				if err != nil {
					fmt.Println("Invalid ID. Please enter a numeric container ID.")
					continue
				}

				// Call function to stop container
				core.ContinueContainer(id)
			}
		case "4":
			exitStopInterface := false
			for !exitStopInterface {
				// Get container ID from user
				fmt.Print("Enter the ID of the container you wish to remove or 'q' to quit: ")
				reader := bufio.NewReader(os.Stdin)
				idStr, _ := reader.ReadString('\n')
				idStr = strings.TrimSpace(idStr) // Remove trailing newline and any spaces

				// Check if user wants to quit
				if idStr == "q" {
					fmt.Println("Returning to the admin controller.")
					exitStopInterface = true
					continue
				}

				// Convert string ID to integer
				id, err := strconv.Atoi(idStr)
				if err != nil {
					fmt.Println("Invalid ID. Please enter a numeric container ID.")
					continue
				}

				core.RemoveContainer(id)
			}

		case "5":
			core.RemoveAllContainers()

		case "6":
			fmt.Println("Exiting admin interface.")
			return

		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}

}
