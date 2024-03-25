package configs

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func CreateConfigs(domainName string, orgPeers map[string]int) {

	CreateDockerComposeCA(orgPeers)
	CreateDockerComposeMembers(domainName, orgPeers)
	// ReadConfig()
	CreateConfigTx(domainName, orgPeers)
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

	type Env struct {
	}

	// create configs for all the organisations
	for org := range orgPeers {
		org := strings.ToLower(org)
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

func CreateDockerComposeMembers(domainName string, orgPeers map[string]int) {

	//viper.KeyDelimiter(":") to adjest the key delimiter from "." to ":"
	// for adding keys like "orderer.example.com"
	var custom_viper = viper.NewWithOptions(viper.KeyDelimiter(":"))

	custom_viper.SetConfigName("docker-compose-orgs")
	custom_viper.SetConfigType("yaml")
	custom_viper.AddConfigPath("pkg/configs/generated/docker")
	custom_viper.Set("version", "3.7")
	custom_viper.Set("networks:test:name", "fabric_test")
	// volumes will be added when the peers are created

	// creating configs for ordering service
	custom_viper.Set(fmt.Sprintf("volumes:orderer.%v", domainName), "")
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:container_name", domainName), fmt.Sprintf("orderer.%v", domainName))
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:image", domainName), "hyperledger/fabric-orderer:2.5.4")
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:labels:service", domainName), "hyperledger-fabric")

	ordererEnv := []string{
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
		fmt.Sprintf("ORDERER_OPERATIONS_LISTENADDRESS=orderer.%v:9443", domainName),
		"ORDERER_METRICS_PROVIDER=prometheus",
	}
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:environment", domainName), ordererEnv)
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:working_dir", domainName), "/root")
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:command", domainName), "orderer")

	// correct the domain name or keep the example.com
	ordererVolumes := []string{
		fmt.Sprintf("../organizations/ordererOrganizations/%v/orderers/orderer.%v/msp:/var/hyperledger/orderer/msp", domainName, domainName),
		fmt.Sprintf("../organizations/ordererOrganizations/%v/orderers/orderer.%v/tls/:/var/hyperledger/orderer/tls", domainName, domainName),
		fmt.Sprintf("orderer.%v:/var/hyperledger/production/orderer", domainName),
	}
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:volumes", domainName), ordererVolumes)

	orderePorts := []string{
		"7050:7050",
		"7053:7053",
		"9443:9443",
	}
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:ports", domainName), orderePorts)

	networkSlice := []string{
		"test",
	}
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:networks", domainName), networkSlice)

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

	CLIVolumes := []string{
		"/var/run/docker.sock:/host/var/run/docker.sock",
		"../organizations:/opt/gopath/src/github.com/hyperledger/fabric/peer/organizations",
	}
	custom_viper.Set("services:cli:volumes", CLIVolumes)
	custom_viper.Set("services:cli:networks", networkSlice)

	// CLI depends will be added from the for loop
	CLIDepends := []string{}

	// for creating port numbers dynamically as well keeping the peer count
	i := 0
	ports := []int{
		5984,
		7051,
		9444,
	}
	caser := cases.Title(language.English)

	// creating couchdb and peers for all the orgs
	for org, peers := range orgPeers {

		org := strings.ToLower(org)
		orgMSP := fmt.Sprintf("%vMSP", caser.String(org))

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
				fmt.Sprintf("%v:5984", ports[0]+i*2000),
			}
			custom_viper.Set(fmt.Sprintf("services:%vpeer%vdb:ports", org, peer), portsCouch)
			custom_viper.Set(fmt.Sprintf("services:%vpeer%vdb:networks", org, peer), networkSlice)

			// peer config
			custom_viper.Set(fmt.Sprintf("services:peer%v.%v.%v:container_name", peer, org, domainName), fmt.Sprintf("services:peer%v.%v.%v", peer, org, domainName))
			custom_viper.Set(fmt.Sprintf("services:peer%v.%v.%v:image", peer, org, domainName), "hyperledger/fabric-peer:2.5.4")
			custom_viper.Set(fmt.Sprintf("services:peer%v.%v.%v:labels:service", peer, org, domainName), "hyperledger-fabric")

			peerEnv := []string{
				"FABRIC_LOGGING_SPEC=INFO",
				"#- FABRIC_LOGGING_SPEC=DEBUG",
				"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
				"CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=fabric_test",
				"CORE_PEER_TLS_ENABLED=true",
				"CORE_PEER_PROFILE_ENABLED=false",
				"CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt",
				"CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key",
				"CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt",
				"# Peer specific variables",
				fmt.Sprintf("CORE_PEER_ID=peer%v.%v.%v", peer, org, domainName),
				fmt.Sprintf("CORE_PEER_ADDRESS=peer%v.%v.%v:%v", peer, org, domainName, ports[1]+i*2000),
				fmt.Sprintf("CORE_PEER_LISTENADDRESS=0.0.0.0:%v", ports[1]+i*2000),
				fmt.Sprintf("CORE_PEER_CHAINCODEADDRESS=peer%v.%v.%v:%v", peer, org, domainName, ports[1]+i*2000+1),
				fmt.Sprintf("CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:%v", ports[1]+i*2000+1),
				fmt.Sprintf("CORE_PEER_GOSSIP_BOOTSTRAP=peer%v.%v.%v:%v", peer, org, domainName, ports[1]+i*2000),
				fmt.Sprintf("CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer%v.%v.%v:%v", peer, org, domainName, ports[1]+i*2000),
				fmt.Sprintf("CORE_PEER_LOCALMSPID=%v", orgMSP),
				"CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/msp",
				fmt.Sprintf("CORE_OPERATIONS_LISTENADDRESS=peer%v.%v.%v:%v", peer, org, domainName, ports[2]+i*1),
				"CORE_METRICS_PROVIDER=prometheus",
				fmt.Sprintf("CHAINCODE_AS_A_SERVICE_BUILDER_CONFIG={'peername':'peer%v%v'}", peer, org),
				"CORE_CHAINCODE_EXECUTETIMEOUT=300s",
				"CORE_LEDGER_STATE_STATEDATABASE=CouchDB",
				fmt.Sprintf("CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=%v:5984", fmt.Sprintf("%vpeer%vdb", org, peer)),
				"CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=admin",
				"CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=adminpw",
			}
			custom_viper.Set(fmt.Sprintf("services:peer%v.%v.%v:environment", peer, org, domainName), peerEnv)

			peerVolumes := []string{
				"/var/run/docker.sock:/host/var/run/docker.sock",
				fmt.Sprintf("../organizations/peerOrganizations/%v.%v/peers/peer%v.%v.%v:/etc/hyperledger/fabric", org, domainName, peer, org, domainName),
				fmt.Sprintf("peer%v.%v.%v:/var/hyperledger/production", peer, org, domainName),
			}
			custom_viper.Set(fmt.Sprintf("services:peer%v.%v.%v:volumes", peer, org, domainName), peerVolumes)

			custom_viper.Set(fmt.Sprintf("services:peer%v.%v.%v:working_dir", peer, org, domainName), "/root")
			custom_viper.Set(fmt.Sprintf("services:peer%v.%v.%v:command", peer, org, domainName), "peer node start")
			peerPorts := []string{
				fmt.Sprintf("%v:%v", ports[1]+i*2000, ports[1]+i*2000),
				fmt.Sprintf("%v:%v", ports[2]+i*1, ports[2]+i*1),
			}
			custom_viper.Set(fmt.Sprintf("services:peer%v.%v.%v:ports", peer, org, domainName), peerPorts)

			peerDepends := []string{
				fmt.Sprintf("%vpeer%vdb", org, peer),
			}
			custom_viper.Set(fmt.Sprintf("services:peer%v.%v.%v:depends_on", peer, org, domainName), peerDepends)

			custom_viper.Set(fmt.Sprintf("services:peer%v.%v.%v:networks", peer, org, domainName), networkSlice)

			// adding the peer volumes
			custom_viper.Set(fmt.Sprintf("volumes:peer%v.%v.%v", peer, org, domainName), "")

			// adding peers to depends field of CLI - may be improved
			CLIDepends = append(CLIDepends, fmt.Sprintf("peer%v.%v.%v", peer, org, domainName))
			custom_viper.Set("services:cli:depends_on", CLIDepends)

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

func CreateConfigTx(domainName string, orgPeers map[string]int) {
	
	viper.Reset()
	type Policies struct {
		Readers string `yaml:"Readers"`
		Writers string `yaml:"Writers"`
		Admins  string `yaml:"Admins"`
	}

	type Organization struct {
		Name             string   `yaml:"Name"`
		ID               string   `yaml:"ID"`
		MSPDir           string   `yaml:"MSPDir"`
		Policies         Policies `yaml:"Policies"`
		OrdererEndpoints []string `yaml:"OrdererEndpoints,omitempty"`
	}

	// set the file name, type and path
	viper.SetConfigName("configtx")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("pkg/configs/generated/config")

	// orgsConfig := []map[string]string{
	// 	{
	// 		"name":   "Org1",
	// 		"domain": "org1.example.com",
	// 	},
	// }

	org := Organization{

		Name:   "OrdererOrg",
		ID:     "OrdererMSP",
		MSPDir: "../organizations/ordererOrganizations/auto.com/msp",
		Policies: Policies{
			Readers: "OR('OrdererMSP.member')",
			Writers: "OR('OrdererMSP.member')",
			Admins:  "OR('OrdererMSP.admin')",
		},
		OrdererEndpoints: []string{"orderer.auto.com:7050"},
	}

	// Set the organization configuration in Viper
	viper.Set("Organizations", []Organization{org})

	// viper.Set("Organizations.&OrdererOrg", orgsConfig)

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
	fmt.Println("configtx.yaml Configuration file created/updated successfully!")
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
// 	viper.SetConfigName("configtx") // name of config file (without extension)

// 	// Set the type of the configuration file
// 	viper.SetConfigType("yaml")

// 	// Set the path to look for the configurations file
// 	viper.AddConfigPath("pkg/configs/defaults/config") // path to look for the config file in

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

// volumeMap := map[string]string{
// 	"orderer.example.com":            "",
// 	"peer0.manufacturer.example.com": "",
// 	"peer0.dealer.example.com":       "",
// }
