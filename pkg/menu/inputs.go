package menu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/vineshtk/fabrix/pkg/configs"
)

var domainName string

// var channelName string

func GetInputsFromUser(channelName string, version string) {
	fmt.Print("\n")
	fmt.Print("Enter the domain name (eg: example.com): ")
	fmt.Scan(&domainName)
	fmt.Print("Enter the number of organizations: ")
	numOrganizations, err := getInputAsInt()
	if err != nil || numOrganizations <= 0 {
		fmt.Println("Invalid input for the number of organizations.")
		return
	}

	// Create a map to store the number of peers for each organization
	OrganizationPeers := make(map[string]int)

	// Get organization names and number of peers for each organization
	for i := 1; i <= numOrganizations; i++ {
		fmt.Printf("Enter the name of organization %d: ", i)
		orgName := getInputAsString()

		fmt.Printf("Enter the number of peers for organization %s: ", orgName)
		numPeers, err := getInputAsInt()
		if err != nil || numPeers < 0 {
			fmt.Println("Invalid input for the number of peers.")
			return
		}
		// Store the values in the map
		OrganizationPeers[orgName] = numPeers
	}

	if channelName == "" {
		fmt.Print("Enter the channel name (eg: mychannel): ")
		fmt.Scan(&channelName)
	}

	configs.CreateConfigs(domainName, OrganizationPeers, channelName, version)
	
}

// getInputAsString reads a line of input from the terminal and returns it as a string.
func getInputAsString() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// getInputAsInt reads a line of input from the terminal and converts it to an integer.
func getInputAsInt() (int, error) {
	inputStr := getInputAsString()
	return strconv.Atoi(inputStr)
}
