package menu

import (
	"fmt"

	"github.com/fatih/color"
)

func ShowMainMenu() {

	asciiArt := " _____     _          _      \n" +
		"|  ___|_ _| |__  _ __(_)_  __\n" +
		"| |_ / _' | '_ \\| '__| \\ \\/ /\n" +
		"|  _| (_| | |_) | |  | |>  < \n" +
		"|_|  \\__,_|_.__/|_|  |_/_/\\_\\\n"

	fmt.Print("\n\n")
	color.Blue(asciiArt)
	fmt.Print("\n")

	color.Green("The helper tool for chaincode developers to create fabric network, it does all the heavy lifting for you!!!")
	color.Green("You will be guided throughout the process.")
	color.Green("Prompt will ask you to give all the details of your network, please provide details acordingly.")
	color.Yellow("Let's start...")

	fmt.Print("\n")
}
