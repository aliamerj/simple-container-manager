package container

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

type ContainerInfo struct {
	ID        string
	Status    string
	ParentPID int
	ChildPID  int
}

var mu sync.Mutex
var containers = make(map[int]*ContainerInfo)

func StartContainer() {
	fmt.Println("Starting container...")

	// Locking mutex to safely update shared data
	newID := len(containers) + 1

	// Start listening for admin messages
	go adminMessagesListener()
	// new container
	cmd := exec.Command("/bin/sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Create a pipe
	pr, pw := io.Pipe()
	cmd.Stdin = pr
	// quit the terminal
	go func() {
		defer pw.Close()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "q" {
				containers[newID].Status = "Removed"
				WriteContainersToFile()
				syscall.Kill(containers[newID].ParentPID, syscall.SIGKILL)

			} else if line != "exit" {
				pw.Write([]byte(line + "\n"))
			} else {
				fmt.Println("Please use 'q' to exit the container.")
			}
		}
	}()

	// Start the process (non-blocking) to get the PID
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting the /bin/sh command: %v\n", err)
		os.Exit(1)
	}

	// Create a new ContainerInfo instance
	newContainer := &ContainerInfo{
		ID:        fmt.Sprintf("%d", newID),
		Status:    "Running",
		ChildPID:  cmd.Process.Pid,
		ParentPID: os.Getpid(),
	}

	// Add the new container to the containers map
	addContainer(newID, newContainer)

	// Now block and keep the terminal interactive
	err := cmd.Wait()
	mu.Lock()
	defer mu.Unlock()
	if err != nil {
		if err.Error() == "signal: killed" {
			fmt.Println("You have been removed from the server")
			return
		}
		fmt.Printf("Failed waiting for container to complete: %s\n", err.Error())
		newContainer.Status = "Failed"
	} else {
		newContainer.Status = "Stopped"
	}
	WriteContainersToFile()
}

func adminMessagesListener() {
	// Create a channel to receive signals
	fmt.Printf("Listening for signals in PID: %d\n", os.Getpid())
	sigs := make(chan os.Signal, 1)

	// Listen for SIGSTOP and SIGCONT
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		// Wait for a signal
		sig := <-sigs

		switch sig {
		case syscall.SIGUSR1:
			fmt.Println("Container has been paused by the admin.")
		case syscall.SIGUSR2:
			fmt.Println("Container has been resumed by the admin.")
		}
	}
}

func WriteContainersToFile() error {
	mu.Lock()
	defer mu.Unlock()
	data, err := json.Marshal(containers)
	if err != nil {
		return err
	}
	err = os.WriteFile("../data/containers.json", data, 0644)
	return err
}
func LoadContainersFromFile() (error, map[int]*ContainerInfo) {
	mu.Lock()
	defer mu.Unlock()
	data, err := os.ReadFile("../data/containers.json")
	if err != nil {
		return err, containers
	}
	if err := json.Unmarshal(data, &containers); err != nil {
		return err, containers
	}
	return nil, containers

}
func addContainer(id int, container *ContainerInfo) {
	containers[id] = container
	WriteContainersToFile()
}
func DeleteContainer(id int, containers map[int]*ContainerInfo) {
	delete(containers, id)
	WriteContainersToFile()
}
