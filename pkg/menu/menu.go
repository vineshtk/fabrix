package menu

import (
	"fmt"
)

func ShowMainMenu() {
	fmt.Println("Welcome to Fabrix - The helper tool for chaincode developers to create fabric network, it does all the heavy lifting for you!!!")
	fmt.Println("You will be guided during throughout the process. Let's start...")
	fmt.Print("\n\n")
	fmt.Println("MENU")
	fmt.Println("N - New network")
	fmt.Println("S - Select an existing network")
	fmt.Println("D - Docker status")
	fmt.Println("C - Clean all Docker resources")
	fmt.Println("Q - Quit")
	fmt.Print("\n")
	fmt.Print("Please select from the options: N, S, D, C, Q : ")

	// color.Cyan("Hello, this text is in cyan!")

}
