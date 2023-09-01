package core

import (
	"fmt"
	"strconv"

	"syscall"

	"containerManger.com/pkg/container"
)

// ...  utility functions that only admin can execute

type User struct {
	ID    string
	State string
}

func getContainers() map[int]*container.ContainerInfo {
	err, containers := container.LoadContainersFromFile()
	if err != nil {
		fmt.Println("Failed to load containers fro json")
		return nil
	}
	return containers
}
func ListContainers() {
	containers := getContainers()
	fmt.Println("\n<><><><><><><><><><><><><><><><><><><><><><><>")
	for id, con := range containers {
		if con.Status == "Removed" {
			id, _ := strconv.Atoi(con.ID)
			delete(containers, id)
			container.WriteContainersToFile()
			continue
		}
		fmt.Printf("ID: %v, Status: %s, Parent-PID: %d, Child-PID: %d\n", id, con.Status, con.ParentPID, con.ChildPID)
	}

}

func RemoveContainer(ID int) {
	containers := getContainers()
	con, exists := containers[ID]
	if !exists {
		fmt.Println("Container with ID", ID, "does not exist.")
		return
	}
	err := syscall.Kill(con.ChildPID, syscall.SIGKILL)
	if err != nil {
		fmt.Printf("Failed to kill process with ID %d: %s\n", ID, err)
		return
	}
	container.DeleteContainer(ID, containers)
	fmt.Println("Container with ID", ID, "has been removed.")

}
func StopContainer(ID int) {
	containers := getContainers()
	con, exists := containers[ID]
	if !exists {
		fmt.Println("Container with ID", ID, "does not exist.")
		return
	}
	if con.Status == "Stopped" {
		fmt.Println("The Container with ID", ID, "is already stopped...")
		return

	}
	if er := syscall.Kill(con.ParentPID, syscall.SIGUSR1); er != nil {
		fmt.Printf("Failed to stop process with ID %d: %s\n", ID, er)
	}
	err := syscall.Kill(con.ChildPID, syscall.SIGSTOP)
	con.Status = "Stopped"
	if err != nil {
		fmt.Printf("Failed to kill process with ID %d: %s\n", ID, err)
		return
	}

	fmt.Println("Container with ID", ID, "has been stopped.")
	container.WriteContainersToFile()

}
func ContinueContainer(ID int) {
	containers := getContainers()
	con, exists := containers[ID]
	if !exists {
		fmt.Println("Container with ID", ID, "does not exist.")
		return
	}
	if con.Status == "Running" {
		fmt.Println("The Container with ID", ID, "is already running...")
		return

	}
	if er := syscall.Kill(con.ParentPID, syscall.SIGUSR2); er != nil {
		fmt.Printf("Failed to stop process with ID %d: %s\n", ID, er)
	}
	err := syscall.Kill(con.ChildPID, syscall.SIGCONT)
	con.Status = "Running"
	if err != nil {
		fmt.Printf("Failed to kill process with ID %d: %s\n", ID, err)
		return
	}
	fmt.Println("Now Container with ID", ID, "is Running...")
	container.WriteContainersToFile()
}

func RemoveAllContainers() {
	containers := getContainers()
	for id, con := range containers {
		err := syscall.Kill(con.ChildPID, syscall.SIGKILL)
		if err != nil {
			fmt.Printf("Failed to kill process with ID %d: %s\n", id, err)

		}
		container.DeleteContainer(id, containers)

	}
	fmt.Println("All Container have been Removed.")

}
