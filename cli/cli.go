package cli

import (
	"flag"
	"fmt"
	"os"

	"containerManger.com/pkg/admin"
	"containerManger.com/pkg/container"
)

func Run() {
	// Define flags here
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	adminCmd := flag.NewFlagSet("admin", flag.ExitOnError)
	adminPassword := adminCmd.String("p", "", "Admin password")

	// Parse and handle commands
	//* first element is the main file, the second is the command
	if len(os.Args) < 2 {
		fmt.Println("expected 'start','admin' subcommands")
		os.Exit(1)

	}

	switch os.Args[1] {
	case "admin":
		adminCmd.Parse(os.Args[2:])
		if *adminPassword != "password" {
			fmt.Println("Invalid admin password")
			os.Exit(1)

		}
		admin.AdminInterface()

	case "start":
		startCmd.Parse(os.Args[2:])

		container.StartContainer()

	default:
		fmt.Println("Unknown command, expected 'start' or 'stop' ")
		os.Exit(1)

	}

}
