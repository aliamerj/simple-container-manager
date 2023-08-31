package core

import (
	"fmt"

	"os/exec"
	"sync"
	"syscall"
)

type ContainerInfo struct {
	ID     string
	Status string
	PID    int
}

var mu sync.Mutex
var containers = make(map[int]*ContainerInfo)

// ...  utility functions that only admin can execute

func ListContainers() {
	fmt.Printf("number of container is : %d \n<><><><><><><><><><><><><><><><><><><><><><><>\n", len(containers))
	mu.Lock()
	defer mu.Unlock()
	for id, container := range containers {
		fmt.Printf("ID: %v, Status: %s, PID: %d\n", id, container.Status, container.PID)
	}

}

func CreateNewContainer() (string, error) {
	newID := len(containers) + 1

	// Initialize a command to start a new container
	cmd := exec.Command("/bin/sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	// Redirect the output somewhere instead of the terminal
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	// Start the container and run it in the background
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("Failed to start container: %v", err)
	}

	// Create a new ContainerInfo instance
	newContainer := &ContainerInfo{
		ID:     fmt.Sprintf("%d", newID),
		Status: "Running",
		PID:    cmd.Process.Pid,
	}

	// Add the new container to the containers map
	addContainer(newID, newContainer)

	return newContainer.ID, nil
}

func RemoveContainer(ID int) {
	mu.Lock()
	defer mu.Unlock()

	container, exists := containers[ID]
	if !exists {
		fmt.Println("Container with ID", ID, "does not exist.")
		return
	}
	err := syscall.Kill(container.PID, syscall.SIGTERM)
	if err != nil {
		fmt.Printf("Failed to kill process with ID %d: %s\n", ID, err)
		return
	}
	delete(containers, ID)
	fmt.Println("Container with ID", ID, "has been removed.")

}
func StopContainer(ID int) {
	mu.Lock()
	defer mu.Unlock()
	container, exists := containers[ID]
	if !exists {
		fmt.Println("Container with ID", ID, "does not exist.")
		return
	}
	if container.Status == "Stopped" {
		fmt.Println("The Container with ID", ID, "is already stopped...")
		return

	}
	err := syscall.Kill(container.PID, syscall.SIGSTOP)
	container.Status = "Stopped"
	if err != nil {
		fmt.Printf("Failed to kill process with ID %d: %s\n", ID, err)
		return
	}

	fmt.Println("Container with ID", ID, "has been stopped.")

}
func ContinueContainer(ID int) {
	mu.Lock()
	defer mu.Unlock()
	container, exists := containers[ID]
	if !exists {
		fmt.Println("Container with ID", ID, "does not exist.")
		return
	}
	if container.Status == "Running" {
		fmt.Println("The Container with ID", ID, "is already running...")
		return

	}
	err := syscall.Kill(container.PID, syscall.SIGCONT)
	container.Status = "Running"
	if err != nil {
		fmt.Printf("Failed to kill process with ID %d: %s\n", ID, err)
		return
	}
	fmt.Println("Now Container with ID", ID, "is Running...")
}

func RemoveAllContainers() {
	mu.Lock()
	defer mu.Unlock()
	for id, container := range containers {
		err := syscall.Kill(container.PID, syscall.SIGKILL)
		if err != nil {
			fmt.Printf("Failed to kill process with ID %d: %s\n", id, err)

		}
		delete(containers, id)
	}
	fmt.Println("All Container have been Removed.")

}

func addContainer(id int, container *ContainerInfo) {
	mu.Lock()
	defer mu.Unlock()
	containers[id] = container
}
