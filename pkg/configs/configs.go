package configs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func CreateConfigs(orgPeers map[string]int) {

	CreateDockerComposeCA(orgPeers)
	CreateDockerComposeMembers(orgPeers)

}

// The CreateDockerComposeCA is used to create the CAs for all the organisations and orderer
func CreateDockerComposeCA(orgPeers map[string]int) {

	// set the file name, type and path
	viper.SetConfigName("docker-compose-ca")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("pkg/configs/generated/docker")

	// set values for each fields in docker compose file
	viper.Set("version", "3.7")
	viper.Set("networks.test.name", "fabric_test")

	i := 1
	ports := []int{7054, 17054}

	// Create the configuration for orderer CA
	viper.Set("services.ca_orderer.image", "hyperledger/fabric-ca:1.5.7")
	viper.Set("services.ca_orderer.labels.service", "hyperledger-fabric")

	// Set environment variables as a slice of strings for addinf it as seperate fields"
	envSlice := []string{
		"FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server",
		"FABRIC_CA_SERVER_CA_NAME=ca-orderer",
		"FABRIC_CA_SERVER_TLS_ENABLED=true",
		fmt.Sprintf("FABRIC_CA_SERVER_PORT=%v", ports[0]),
		fmt.Sprintf("FABRIC_CA_SERVER_OPERATIONS_LISTENADDRESS=0.0.0.0:%v", ports[1]),
	}
	viper.Set("services.ca_orderer.environment", envSlice)

	//** this nee to be changed since we need to add those as strings
	portSlice := []string{
		// fmt.Sprintf("%d:%d", ports[0], ports[0]),
		fmt.Sprintf("%v:%v", ports[1], ports[1]),
		fmt.Sprintf("%v:%v", ports[1], ports[1]),
	}

	viper.Set("services.ca_orderer.ports", portSlice)
	viper.Set("services.ca_orderer.command", "sh -c 'fabric-ca-server start -b admin:adminpw -d'")

	volumeSlice := [1]string{
		"../organizations/fabric-ca/ordererOrg:/etc/hyperledger/fabric-ca-server",
	}
	viper.Set("services.ca_orderer.volumes", volumeSlice)

	viper.Set("services.ca_orderer.container_name", "ca_orderer")
	networkSlice := [1]string{
		"test",
	}
	viper.Set("services.ca_orderer.networks", networkSlice)

	// create configs for all the organisations
	for org := range orgPeers {
		viper.Set(fmt.Sprintf("services.ca_%v.image", org), "hyperledger/fabric-ca:1.5.7")
		viper.Set(fmt.Sprintf("services.ca_%v.labels.service", org), "hyperledger-fabric")

		// Set environment variables as a slice of strings for addinf it as seperate fields"
		envSlice := []string{
			"FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server",
			fmt.Sprintf("FABRIC_CA_SERVER_CA_NAME=ca-%v", org),
			"FABRIC_CA_SERVER_TLS_ENABLED=true",
			fmt.Sprintf("FABRIC_CA_SERVER_PORT=%v", ports[0]+i*1000),
			fmt.Sprintf("FABRIC_CA_SERVER_OPERATIONS_LISTENADDRESS=0.0.0.0:%v", ports[1]+i*1000),
		}

		viper.Set(fmt.Sprintf("services.ca_%v.environment", org), envSlice)

		//** this nee to be changed since we need to add those as strings
		portSlice := []string{
			fmt.Sprintf("%d:%d", ports[0]+i*1000, ports[0]+i*1000),
			fmt.Sprintf("%d:%d", ports[1]+i*1000, ports[1]+i*1000),
		}

		viper.Set(fmt.Sprintf("services.ca_%v.ports", org), portSlice)
		viper.Set(fmt.Sprintf("services.ca_%v.command", org), "sh -c 'fabric-ca-server start -b admin:adminpw -d'")

		volumeSlice := [1]string{
			fmt.Sprintf("../organizations/fabric-ca/%v:/etc/hyperledger/fabric-ca-server", org),
		}

		viper.Set(fmt.Sprintf("services.ca_%v.volumes", org), volumeSlice)
		viper.Set(fmt.Sprintf("services.ca_%v.container_name", org), fmt.Sprintf("ca_%v", org))

		networkSlice := [1]string{
			"test",
		}

		viper.Set(fmt.Sprintf("services.ca_%v.networks", org), networkSlice)

		err := viper.SafeWriteConfig()
		if err != nil {
			if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
				err = viper.WriteConfig()
				if err != nil {
					log.Fatalf("Error while updating config file %s", err)
				}
			} else {
				log.Fatalf("Error while creating config file %s", err)
			}
		}

		i += 1
	}

	fmt.Println("docker-compose-ca.yaml Configuration file created/updated successfully!")
}

func CreateRegisterEnroll(orgPeers map[string]int) {
	// this need to be addresed
	// how to impliment since viper dont support script file
}

func CreateDockerComposeMembers(orgPeers map[string]int) {

	//viper.KeyDelimiter(":") to adjest the key delimiter from "." to ":"
	// for adding keys like "orderer.example.com"
	var custom_viper = viper.NewWithOptions(viper.KeyDelimiter(":"))

	custom_viper.SetConfigName("docker-compose-orgs")
	custom_viper.SetConfigType("yaml")
	custom_viper.AddConfigPath("pkg/configs/generated/docker")
	custom_viper.Set("version", "3.7")
	custom_viper.Set("networks:test:name", "fabric_test")

	// volumeMap := map[string]string{
	// 	"orderer.example.com":            "",
	// 	"peer0.manufacturer.example.com": "",
	// 	"peer0.dealer.example.com":       "",
	// }

	// creating configs for ordering service
	custom_viper.Set("volumes:orderer.example.com", "")
	custom_viper.Set("services:orderer.example.com:container_name", "orderer.example.com")
	custom_viper.Set("services:orderer.example.com:image", "hyperledger/fabric-orderer:2.5.4")
	custom_viper.Set("services:orderer.example.com:labels:service", "hyperledger-fabric")

	envSlice := []string{
		"FABRIC_LOGGING_SPEC=INFO",
		"ORDERER_GENERAL_LISTENADDRESS=0.0.0.0",
		"ORDERER_GENERAL_LISTENPORT=7050",
		"ORDERER_GENERAL_LOCALMSPID=OrdererMSP",
		"ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp",

		"ORDERER_GENERAL_TLS_ENABLED=true",
		"ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key",
		"ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt",
		"ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]",
		"ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=/var/hyperledger/orderer/tls/server.crt",
		"ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=/var/hyperledger/orderer/tls/server.key",
		"ORDERER_GENERAL_CLUSTER_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]",
		"ORDERER_GENERAL_BOOTSTRAPMETHOD=none",
		"ORDERER_CHANNELPARTICIPATION_ENABLED=true",
		"ORDERER_ADMIN_TLS_ENABLED=true",
		"ORDERER_ADMIN_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt",
		"ORDERER_ADMIN_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key",
		"ORDERER_ADMIN_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]",
		"ORDERER_ADMIN_TLS_CLIENTROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]",
		"ORDERER_ADMIN_LISTENADDRESS=0.0.0.0:7053",
		"ORDERER_OPERATIONS_LISTENADDRESS=orderer.example.com:9443",
		"ORDERER_METRICS_PROVIDER=prometheus",
	}
	custom_viper.Set("services:orderer.example.com:environment", envSlice)
	custom_viper.Set("services:orderer.example.com:working_dir", "/root")
	custom_viper.Set("services:orderer.example.com:command", "orderer")

	// correct the domain name or keep the example.com
	ordererVolumeSlice := []string{
		"../organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp:/var/hyperledger/orderer/msp",
		"../organizations/ordererOrganizations/example.com/orderers/orderer.example.com/tls/:/var/hyperledger/orderer/tls",
		"orderer.example.com:/var/hyperledger/production/orderer",
	}
	custom_viper.Set("services:orderer.example.com:volumes", ordererVolumeSlice)

	orderePortSlice := []string{
		"7050:7050",
		"7053:7053",
		"9443:9443",
	}
	custom_viper.Set("services:orderer.example.com:ports", orderePortSlice)

	networkSlice := []string{
		"automobile",
	}
	custom_viper.Set("services:orderer.example.com:networks", networkSlice)

	// configs for CLI
	custom_viper.Set("services:cli:container_name", "cli")
	custom_viper.Set("services:cli:image", "hyperledger/fabric-tools:2.5.4")
	custom_viper.Set("services:cli:labels:service", "hyperledger-fabric")
	custom_viper.Set("services:cli:tty", "true")
	custom_viper.Set("services:cli:stdin_open", "true")

	envSliceCLI := []string{
		"GOPATH=/opt/gopath",
		"FABRIC_LOGGING_SPEC=INFO",
		"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
	}
	custom_viper.Set("services:cli:environment", envSliceCLI)

	custom_viper.Set("services:cli:working_dir", "/opt/gopath/src/github.com/hyperledger/fabric/peer")
	custom_viper.Set("services:cli:command", "/bin/bash")

	volumeSliceCLI := []string{
		"/var/run/docker.sock:/host/var/run/docker.sock",
		"../organizations:/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations",
	}
	custom_viper.Set("services:cli:volumes", volumeSliceCLI)

	networkSliceCLI := []string{
		"test",
	}
	custom_viper.Set("services:cli:networks", networkSliceCLI)

	//  this need to be added in for loop
	//   depends_on:
	//   - peer0.manufacturer.example.com
	//   - peer0.dealer.example.com
	//   - peer0.mvd.example.com

	// custom_viper.Set("services:cli:depends_on", "cli")
	// custom_viper.Set("services:orderer.example.com:environment", envSlice)

	// for creating port numbers dynamically asw well keeping the peer count
	i := 0
	peerPorts := []int{
		5984,
	}

	// creating couchdb and peers for all the orgs
	for org, peers := range orgPeers {

		for peer := 0; peer < peers; peer++ {

			// couchdb configs
			custom_viper.Set(fmt.Sprintf("services:%vpeer%vdb:container_name", org, peer), fmt.Sprintf("%vpeer%vdb", org, peer))
			custom_viper.Set(fmt.Sprintf("services:%vpeer%vdb:image", org, peer), "couchdb:3.3.2")
			custom_viper.Set(fmt.Sprintf("services:%vpeer%vdb:labels:service", org, peer), "hyperledger-fabric")

			envCouch := []string{
				"COUCHDB_USER=admin",
				"COUCHDB_PASSWORD=adminpw",
			}
			custom_viper.Set(fmt.Sprintf("services:%vpeer%vdb:environment", org, peer), envCouch)

			portsCouch := []string{
				fmt.Sprintf("%v:5984", peerPorts[0] + i*2000),
			}
			custom_viper.Set(fmt.Sprintf("services:%vpeer%vdb:ports", org, peer), portsCouch)
			custom_viper.Set(fmt.Sprintf("services:%vpeer%vdb:networks", org, peer), networkSlice)

			// peer config
			

































			err := custom_viper.SafeWriteConfig()
			if err != nil {
				if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
					err = custom_viper.WriteConfig()
					if err != nil {
						log.Fatalf("Error while updating config file %s", err)
					}
				} else {
					log.Fatalf("Error while creating config file %s", err)
				}
			}
			i++
		}
	}
	fmt.Println("docker-compose-orgs.yaml Configuration file created/updated successfully!")
}

func CreateConfigTx() {

}

// func WriteDockerCa() {
// 	viper.Set("someKey", "newValue")

// 	// Marshal the configuration back to YAML
// 	configContent, err := yaml.Marshal(viper.AllSettings())
// 	if err != nil {
// 		fmt.Println("Error marshaling config to YAML", err)
// 		return
// 	}

// 	// Define the new file path
// 	newPath := "pkg/configs/generated/docker/docker-compose-ca.yaml"

// 	// Write the modified configuration to a new file
// 	err = os.WriteFile(newPath, configContent, 0644) // Adjust permissions as needed
// 	if err != nil {
// 		fmt.Println("Error writing config to new file", err)
// 		return
// 	}

// 	fmt.Println("Configuration written to new file successfully.")
// }

// func ReadConfig() {
// 	// Set the file name of the configurations file
// 	viper.SetConfigName("docker-compose-ca-default") // name of config file (without extension)

// 	// Set the type of the configuration file
// 	viper.SetConfigType("yaml")

// 	// Set the path to look for the configurations file
// 	viper.AddConfigPath("pkg/configs/defaults/docker") // path to look for the config file in

// 	// Find and read the config file
// 	err := viper.ReadInConfig()

// 	if err != nil { // Handle errors reading the config file
// 		log.Fatalf("Error while reading config file %s", err)
// 	}
// }

// func CreateCA(orgPeers map[string]int) {

// 	// Reading from default configs

// 	// Setting up some configurations
// 	viper.Set("version", "3.7")
// 	viper.Set("networks.test.network", "fabric_test")
// 	viper.Set("website", "golang.org")
// 	viper.Set("services.ca_org1", "org1")

// 	// Creating and updating a configuration file
// 	viper.SetConfigName("docker-compose-ca")            // name of config file (without extension)
// 	viper.SetConfigType("yaml")                         // specifying the config type
// 	viper.AddConfigPath("pkg/configs/generated/docker") // path to look for the config file in

// 	// viper.SetConfigFile("pkg/configs/generated/docker/docker-compose-ca.yaml")
// 	err := viper.SafeWriteConfig()
// 	if err != nil {
// 		if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
// 			err = viper.WriteConfig()
// 			if err != nil {
// 				log.Fatalf("Error while updating config file %s", err)
// 			}
// 		} else {
// 			log.Fatalf("Error while creating config file %s", err)
// 		}
// 	}

// 	fmt.Println("Configuration file created/updated successfully!")
// }
