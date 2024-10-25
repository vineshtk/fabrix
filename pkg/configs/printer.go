package configs

import (
	"encoding/json"
	"fmt"
	"os"
)

// Function to print organisation details neatly
func printOrganisation(org Organisation) {
	fmt.Printf("Organisation Name: %s\n", org.Name)
	fmt.Printf("MSP ID: %s\n", org.MSPId)

	fmt.Printf("\nCA Info:\n")
	fmt.Printf("  Name: %s\n", org.Ca.Name)
	fmt.Printf("  Port: %d\n", org.Ca.Port)

	fmt.Printf("\nPeers:\n")
	for i, peer := range org.Peers {
		fmt.Printf("  Peer %d:\n", i+1)
		fmt.Printf("    Name: %s\n", peer.Name)
		fmt.Printf("    Port: %d\n", peer.Port)
		fmt.Printf("    CouchDB Name: %s\n", peer.CouchDbName)
		fmt.Printf("    CouchDB Port: %d\n", peer.CouchDbPort)
	}

	// fmt.Printf("\nAdmin:\n")
	// fmt.Printf("  Name: %s\n", org.Admin.Name)

	// fmt.Printf("\nUser:\n")
	// fmt.Printf("  Name: %s\n", org.User.Name)
}

// Example to print the entire NetworkInfo
func printNetworkInfo(networkInfo *NetworkInfo) {
	fmt.Printf("Network Name: %s\n", networkInfo.NetworkName)
	fmt.Printf("Domain Name: %s\n", networkInfo.DomainName)
	fmt.Printf("Number of Organisations: %d\n", networkInfo.NumberOfOrganisations)
	fmt.Printf("Channel Name: %s\n", info.ChannelName)

	fmt.Printf("\nOrderer Info:\n")
	printOrganisation(networkInfo.Orderer)
	fmt.Printf("\nOrganisations Info:")
	for _, org := range networkInfo.Organisations {
		printOrganisation(org)
		fmt.Println("---------------------------------")
	}
}

// Save NetworkInfo to a JSON file
func SaveNetworkInfoToFile(info *NetworkInfo, filePath string) error {
	// Convert the NetworkInfo struct to JSON
	data, err := json.MarshalIndent(info, "", "  ") // Pretty-print JSON
	if err != nil {
		return fmt.Errorf("failed to marshal network info: %v", err)
	}

	// Create or open the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write the JSON data to the file
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	fmt.Println("Network info saved to", filePath)
	return nil
}


