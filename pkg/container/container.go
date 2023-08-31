package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func StartContainer() {
	fmt.Printf("container is ON .....")
	// The command to run inside the container
	// For example, running /bin/sh
	cmd := exec.Command("/bin/sh")

	// Setting up namespaces using the Go syscall package

	// We're setting the SysProcAttr attribute of the Cmd structure to a new SysProcAttr object.
	// This object has a field called Cloneflags that is set to combination of various flags:
	//* CLONE_NEWUTS: Creates a new UTS (UNIX Time Sharing) namespace.
	//* This isolates hostname and NIS domain name.
	//* CLONE_NEWPID: Creates a new PID (Process ID) namespace.
	//* This isolates the process ID number space,
	//* meaning the process IDs in the new namespace start at 1 as if this is a new system.
	//* CLONE_NEWNS: Creates a new mount namespace.
	//* This isolates the set of filesystem mount points seen by the processes in the new namespace.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	// Standard IO setup : These flags are Linux-specific system calls for creating new namespaces,
	// and they're what allow our container to isolate certain resources from the host system.
	//* This means that the child process will take its input from the same place as the parent process.
	//* If you type something into the parent process's terminal, the child process will receive that input.
	cmd.Stdin = os.Stdin
	//*  This means that the output and errors from the child process will
	//* be displayed in the same terminal as the parent process.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the /bin/sh command - %s\n", err)
		os.Exit(1)
	}

}
func StopContainer() {
	fmt.Printf("Container is OFF ... ")

}
