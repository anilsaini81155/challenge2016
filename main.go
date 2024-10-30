package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	loadCities()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nSelect an option:")
		fmt.Println("1) Create distributor")
		fmt.Println("2) Add permission")
		fmt.Println("3) Link distributor")
		fmt.Println("4) Check permissions")
		fmt.Println("5) Exit")
		fmt.Print("Enter choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Print("Enter distributor name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			createDistributor(name)
		case "2":
			fmt.Print("Enter distributor name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			distributor := distributors[name]
			if distributor == nil {
				fmt.Println("Distributor not found.")
				continue
			}
			fmt.Print("Enter permission type (include/exclude): ")
			permTypeInput, _ := reader.ReadString('\n')
			permTypeInput = strings.TrimSpace(strings.ToLower(permTypeInput))
			var permType PermissionType
			if permTypeInput == "include" {
				permType = Include
			} else if permTypeInput == "exclude" {
				permType = Exclude
			} else {
				fmt.Println("Invalid permission type.")
				continue
			}
			fmt.Print("Enter regions (comma-separated, e.g., INDIA, UNITEDSTATES): ")
			regionsInput, _ := reader.ReadString('\n')
			regions := strings.Split(regionsInput, ",")
			addPermissions(distributor, permType, regions)
			fmt.Printf("Permissions added to %s\n", distributor.Name)
		case "3":
			fmt.Print("Enter parent distributor: ")
			parent, _ := reader.ReadString('\n')
			fmt.Print("Enter child distributor: ")
			child, _ := reader.ReadString('\n')
			err := linkDistributor(strings.TrimSpace(parent), strings.TrimSpace(child))
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "4":
			fmt.Print("Enter distributor name to check: ")
			name, _ := reader.ReadString('\n')
			distributor := distributors[strings.TrimSpace(name)]
			if distributor == nil {
				fmt.Println("Distributor not found.")
				continue
			}
			showEffectivePermissions(distributor)
		case "5":
			fmt.Println("Exiting.")
			return
		default:
			fmt.Println("Invalid choice.")
		}
	}
}
