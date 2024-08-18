package configs

import "fmt"

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

	fmt.Printf("\nAdmin:\n")
	fmt.Printf("  Name: %s\n", org.Admin.Name)

	fmt.Printf("\nUser:\n")
	fmt.Printf("  Name: %s\n", org.User.Name)
}

// Example to print the entire NetworkInfo
func printNetworkInfo(networkInfo *NetworkInfo) {
	fmt.Printf("Network Name: %s\n", networkInfo.NetworkName)
	fmt.Printf("Domain Name: %s\n", networkInfo.DomainName)
	fmt.Printf("Number of Organisations: %d\n", networkInfo.NumberOfOrganisations)

	fmt.Printf("\nOrderer Info:\n")
	printOrganisation(networkInfo.Orderer)

	fmt.Printf("\nOrganisations Info:\n")
	for _, org := range networkInfo.Organisations {
		printOrganisation(org)
		fmt.Println("---------------------------------")
	}
}
