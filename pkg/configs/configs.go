package configs

import (
	_ "embed"
	"fmt"
	"sort"

	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type NetworkInfo struct {
	NumberOfOrganisations int            `json:"numberOfOrganisations,omitempty"`
	DomainName            string         `json:"domainName,omitempty"`
	NetworkName           string         `json:"networkName,omitempty"`
	Organisations         []Organisation `json:"organisations,omitempty"`
	Orderer               Organisation   `json:"orderer,omitempty"`
}

type Organisation struct {
	Name  string `json:"name,omitempty"`
	Ca    CA     `json:"ca,omitempty"`
	Peers []Peer `json:"peers,omitempty"`
	MSPId string `json:"MSPId,omitempty"`
	Admin Admin  `json:"admin,omitempty"`
	User  User   `json:"user,omitempty"`
}

type CA struct {
	Name string `json:"name,omitempty"`
	Port int    `json:"port,omitempty"`
}

type Peer struct {
	Name        string `json:"name,omitempty"`
	Port        int    `json:"port,omitempty"`
	CouchDbName string `json:"couchDbName,omitempty"`
	CouchDbPort int    `json:"couchDbPort,omitempty"`
}

type Admin struct {
	Name string `json:"name,omitempty"`
}

type User struct {
	Name string `json:"name,omitempty"`
}

var info *NetworkInfo

// Create a slice to store the map keys - this is to preserve the looping order
var keys []string

//go:embed defaults/peercfg/core.yaml
var coreYaml []byte

func CreateConfigs(domainName string, orgPeers map[string]int) {
	CreateFolders(domainName)
	CreateDockerComposeCA(domainName, orgPeers)
	CreateDockerComposeMembers(domainName, orgPeers)
	CreateConfigTx(domainName, orgPeers)
	createPeercfg(domainName)
	CreateRegisterEnroll(domainName, orgPeers)
	CreateStartNetwork(domainName, orgPeers)
	createStopNetwork(domainName)

	// ReadCaConfig(domainName)
}

// The CreateDockerComposeCA is used to create the CAs for all the organisations and orderer
func CreateDockerComposeCA(domainName string, orgPeers map[string]int) {

	// set the file name, type and path
	viper.SetConfigName("docker-compose-ca")
	viper.SetConfigType("yaml")
	path := fmt.Sprintf("fabrix/%v/Network/docker", domainName)
	viper.AddConfigPath(path)

	// updating network Info variable when the file is being generated
	// for reusing later for generating certificates and script files

	// set values for each fields in docker compose file
	viper.Set("version", "3.7")
	viper.Set("networks.test.name", "fabric_test")

	info = &NetworkInfo{}
	// change this value dynamically later
	info.NumberOfOrganisations = 3

	info.DomainName = domainName
	info.NetworkName = "fabric_test"

	// fmt.Println("this is network info", Info)

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
	info.Orderer.Ca.Name = "ca-orderer"
	info.Orderer.Ca.Port = ports[0]
	info.Orderer.Name = "orderer"

	//** this nee to be changed since we need to add those as strings
	portSlice := []string{
		// fmt.Sprintf("%d:%d", ports[0], ports[0]),
		fmt.Sprintf("%v:%v", ports[0], ports[0]),
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

	// Append map keys to the slice for preserving looping order
	for key := range orgPeers {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// create configs for all the organisations
	for _, org := range keys {
		fmt.Println("this is range order at CA", org)
		org := strings.ToLower(org)
		viper.Set(fmt.Sprintf("services.ca_%v.image", org), "hyperledger/fabric-ca:1.5.7")
		viper.Set(fmt.Sprintf("services.ca_%v.labels.service", org), "hyperledger-fabric")

		// Set environment variables as a slice of strings for addinf it as seperate fields"
		envSlice := []string{
			"FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server",
			fmt.Sprintf("FABRIC_CA_SERVER_CA_NAME=ca-%v", org),
			"FABRIC_CA_SERVER_TLS_ENABLED=true",
			fmt.Sprintf("FABRIC_CA_SERVER_PORT=%v", ports[0]+i*10),
			fmt.Sprintf("FABRIC_CA_SERVER_OPERATIONS_LISTENADDRESS=0.0.0.0:%v", ports[1]+i*10),
		}

		viper.Set(fmt.Sprintf("services.ca_%v.environment", org), envSlice)

		//** this need to be changed since we need to add those as strings
		portSlice := []string{
			fmt.Sprintf("%d:%d", ports[0]+i*10, ports[0]+i*10),
			fmt.Sprintf("%d:%d", ports[1]+i*10, ports[1]+i*10),
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

		if err := viper.SafeWriteConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
				if err = viper.WriteConfig(); err != nil {
					log.Fatalf("Error while updating config file %s", err)
				}
			} else {
				log.Fatalf("Error while creating config file %s", err)
			}
		}
		organisation := Organisation{
			Name: org,
			Ca: CA{
				Name: fmt.Sprintf("ca-%s", org),
				Port: ports[0] + i*10,
			},
		}

		info.Organisations = append(info.Organisations, organisation)

		i += 1
	}

	fmt.Println("docker-compose-ca.yaml Configuration file created/updated successfully!")
}

func CreateDockerComposeMembers(domainName string, orgPeers map[string]int) {

	//viper.KeyDelimiter(":") to adjust the key delimiter from "." to ":"
	// for adding keys like "orderer.example.com"
	var custom_viper = viper.NewWithOptions(viper.KeyDelimiter(":"))

	custom_viper.SetConfigName("docker-compose-orgs")
	custom_viper.SetConfigType("yaml")

	path := fmt.Sprintf("fabrix/%v/Network/docker", domainName)
	custom_viper.AddConfigPath(path)

	custom_viper.Set("version", "3.7")
	custom_viper.Set("networks:test:name", "fabric_test")
	// volumes will be added when the peers are created

	// creating configs for ordering service
	custom_viper.Set(fmt.Sprintf("volumes:orderer.%v", domainName), map[string]string{})
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
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:environment", domainName), ordererEnv)
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:working_dir", domainName), "/root")
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:command", domainName), "orderer")

	ordererVolumes := []string{
		fmt.Sprintf("../organizations/ordererOrganizations/%v/orderers/orderer.%v/msp:/var/hyperledger/orderer/msp", domainName, domainName),
		fmt.Sprintf("../organizations/ordererOrganizations/%v/orderers/orderer.%v/tls/:/var/hyperledger/orderer/tls", domainName, domainName),
		fmt.Sprintf("orderer.%v:/var/hyperledger/production/orderer", domainName),
	}
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:volumes", domainName), ordererVolumes)

	ordererPorts := []string{
		"7050:7050",
		"7053:7053",
		"9443:9443",
	}
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:ports", domainName), ordererPorts)

	networkSlice := []string{
		"test",
	}
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:networks", domainName), networkSlice)
	custom_viper.Set(fmt.Sprintf("services:orderer.%v:networks", domainName), networkSlice)

	// info for orderer details
	info.Orderer.MSPId = "OrdererMSP"
	ordererPeer := Peer{
		Name: "orderer",
		Port: 7050,
	}
	info.Orderer.Peers = append(info.Orderer.Peers, ordererPeer)

	// configs for CLI
	custom_viper.Set("services:cli:container_name", "cli")
	custom_viper.Set("services:cli:image", "hyperledger/fabric-tools:2.5.4")
	custom_viper.Set("services:cli:labels:service", "hyperledger-fabric")
	custom_viper.Set("services:cli:tty", true)
	custom_viper.Set("services:cli:stdin_open", true)

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
	for index, org := range keys {
		peers := orgPeers[org]
		fmt.Println("this is range order", org)
		org := strings.ToLower(org)
		orgMSP := fmt.Sprintf("%vMSP", caser.String(org))
		peerList := []Peer{}

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
				fmt.Sprintf("%v:5984", ports[0]+i*20),
			}
			custom_viper.Set(fmt.Sprintf("services:%vpeer%vdb:ports", org, peer), portsCouch)
			custom_viper.Set(fmt.Sprintf("services:%vpeer%vdb:networks", org, peer), networkSlice)

			// peer config
			custom_viper.Set(fmt.Sprintf("services:peer%v.%v.%v:container_name", peer, org, domainName), fmt.Sprintf("peer%v.%v.%v", peer, org, domainName))
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
				fmt.Sprintf("CORE_PEER_ADDRESS=peer%v.%v.%v:%v", peer, org, domainName, ports[1]+i*20),
				fmt.Sprintf("CORE_PEER_LISTENADDRESS=0.0.0.0:%v", ports[1]+i*20),
				fmt.Sprintf("CORE_PEER_CHAINCODEADDRESS=peer%v.%v.%v:%v", peer, org, domainName, ports[1]+i*20+1),

				fmt.Sprintf("CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:%v", ports[1]+i*20+1),
				fmt.Sprintf("CORE_PEER_GOSSIP_BOOTSTRAP=peer%v.%v.%v:%v", peer, org, domainName, ports[1]+i*20),
				fmt.Sprintf("CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer%v.%v.%v:%v", peer, org, domainName, ports[1]+i*20),

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
				fmt.Sprintf("%v:%v", ports[1]+i*20, ports[1]+i*20),
				fmt.Sprintf("%v:%v", ports[2]+i*1, ports[2]+i*1),
			}
			custom_viper.Set(fmt.Sprintf("services:peer%v.%v.%v:ports", peer, org, domainName), peerPorts)

			peerDepends := []string{
				fmt.Sprintf("%vpeer%vdb", org, peer),
			}
			custom_viper.Set(fmt.Sprintf("services:peer%v.%v.%v:depends_on", peer, org, domainName), peerDepends)

			custom_viper.Set(fmt.Sprintf("services:peer%v.%v.%v:networks", peer, org, domainName), networkSlice)

			// adding the peer volumes
			custom_viper.Set(fmt.Sprintf("volumes:peer%v.%v.%v", peer, org, domainName), map[string]string{})

			// adding peers to depends field of CLI - may be improved
			CLIDepends = append(CLIDepends, fmt.Sprintf("peer%v.%v.%v", peer, org, domainName))
			custom_viper.Set("services:cli:depends_on", CLIDepends)

			if err := custom_viper.SafeWriteConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
					if err = custom_viper.WriteConfig(); err != nil {
						log.Fatalf("Error while updating config file %s", err)
					}
				} else {
					log.Fatalf("Error while creating config file %s", err)
				}
			}
			// info for orgs peer details
			orgPeer := Peer{
				Name:        fmt.Sprintf("peer%v.%v.%v", peer, org, domainName),
				Port:        ports[1] + i*20,
				CouchDbName: fmt.Sprintf("%vpeer%vdb", org, peer),
				CouchDbPort: ports[0] + i*20,
			}

			peerList = append(peerList, orgPeer)
			i++
		}

		info.Organisations[index].MSPId = orgMSP
		info.Organisations[index].Peers = peerList
	}
	fmt.Println("docker-compose-orgs.yaml Configuration file created/updated successfully!")
}

func CreateConfigTx(domainName string, orgPeers map[string]int) {

	viper.Reset()

	type Rule struct {
		Type string `yaml:"Type"`
		Rule string `yaml:"Rule"`
	}

	type Policies struct {
		Readers              Rule `yaml:"Readers"`
		Writers              Rule `yaml:"Writers"`
		Admins               Rule `yaml:"Admins"`
		LifecycleEndorsement Rule `yaml:"LifecycleEndorsement,omitempty"`
		Endorsement          Rule `yaml:"Endorsement,omitempty"`
		BlockValidation      Rule `yaml:"BlockValidation,omitempty"`
	}

	type Organization struct {
		Name             string   `yaml:"Name"`
		ID               string   `yaml:"ID"`
		MSPDir           string   `yaml:"MSPDir"`
		Policies         Policies `yaml:"Policies"`
		OrdererEndpoints []string `yaml:"OrdererEndpoints,omitempty"`
	}

	type BatchStruct struct {
		MaxMessageCount   int
		AbsoluteMaxBytes  string
		PreferredMaxBytes string
	}

	type ApplicationDefaults struct {
		Organizations map[string]string          `yaml:"Organizations"`
		Policies      Policies                   `yaml:"Policies"`
		Capabilities  map[string]map[string]bool `yaml:"Capabilities"`
	}

	type OrderDefaults struct {
		Addresses     []string          `yaml:"Addresses"`
		BatchTimeout  string            `yaml:"BatchTimeout"`
		BatchSize     BatchStruct       `yaml:"BatchSize"`
		Organizations map[string]string `yaml:"Organizations"`
		Policies      Policies          `yaml:"Policies"`
	}

	type ChannelDefaults struct {
		Policies     Policies                   `yaml:"Policies"`
		Capabilities map[string]map[string]bool `yaml:"Capabilities"`
	}
	// set the file name, type and path
	viper.SetConfigName("configtx")
	viper.SetConfigType("yaml")
	path := fmt.Sprintf("fabrix/%v/Network/config", domainName)
	viper.AddConfigPath(path)

	ordererOrg := Organization{
		Name:   "OrdererOrg",
		ID:     "OrdererMSP",
		MSPDir: fmt.Sprintf("../organizations/ordererOrganizations/%v/msp", domainName),
		Policies: Policies{
			Readers: Rule{Type: "Signature", Rule: `OR('OrdererMSP.member')`},
			Writers: Rule{Type: "Signature", Rule: "OR('OrdererMSP.member')"},
			Admins:  Rule{Type: "Signature", Rule: "OR('OrdererMSP.member')"},
		},
		OrdererEndpoints: []string{fmt.Sprintf("orderer.%v:7050", domainName)},
	}

	caser := cases.Title(language.English)
	orgs := []Organization{ordererOrg}
	for org := range orgPeers {

		org := strings.ToLower(org)
		orgMSP := fmt.Sprintf("%vMSP", caser.String(org))

		otherOrg := Organization{
			Name:   orgMSP,
			ID:     orgMSP,
			MSPDir: fmt.Sprintf("../organizations/peerOrganizations/%v.%v/msp", org, domainName),
			Policies: Policies{
				Readers:     Rule{Type: "Signature", Rule: fmt.Sprintf("OR('%v.admin', '%v.peer', '%v.client')", orgMSP, orgMSP, orgMSP)},
				Writers:     Rule{Type: "Signature", Rule: fmt.Sprintf("OR('%v.admin','%v.client')", orgMSP, orgMSP)},
				Admins:      Rule{Type: "Signature", Rule: fmt.Sprintf("OR('%v.admin')", orgMSP)},
				Endorsement: Rule{Type: "Signature", Rule: fmt.Sprintf("OR('%v.peer')", orgMSP)},
			},
		}
		orgs = append(orgs, otherOrg)
	}

	// Set the organization configuration in Viper
	viper.Set("Organizations", orgs)

	Capabilities := map[string]map[string]bool{
		"Channel":     {"V2_0": true},
		"Orderer":     {"V2_0": true},
		"Application": {"V2_5": true},
	}

	viper.Set("Capabilities.Channel", Capabilities["Channel"])
	viper.Set("Capabilities.Orderer", Capabilities["Orderer"])
	viper.Set("Capabilities.Application", Capabilities["Application"])

	// viper.Set("Application.Organizations", "")
	applicationDefaults := ApplicationDefaults{
		Organizations: map[string]string{},
		Policies:      Policies{Readers: Rule{Type: "ImplicitMeta", Rule: "ANY Readers"}, Writers: Rule{Type: "ImplicitMeta", Rule: "ANY Writers"}, Admins: Rule{Type: "ImplicitMeta", Rule: "MAJORITY Admins"}, LifecycleEndorsement: Rule{Type: "ImplicitMeta", Rule: "MAJORITY Endorsement"}, Endorsement: Rule{Type: "ImplicitMeta", Rule: "MAJORITY Endorsement"}},
		Capabilities:  map[string]map[string]bool{"<<": Capabilities["Application"]},
	}
	viper.Set("Application", applicationDefaults)

	orderDefaults := OrderDefaults{
		Addresses:     []string{fmt.Sprintf("orderer.%v:7050", domainName)},
		BatchTimeout:  "2s",
		BatchSize:     BatchStruct{MaxMessageCount: 10, AbsoluteMaxBytes: "99 MB", PreferredMaxBytes: "512 KB"},
		Organizations: map[string]string{},
		Policies:      Policies{Readers: Rule{Type: "ImplicitMeta", Rule: "ANY Readers"}, Writers: Rule{Type: "ImplicitMeta", Rule: "ANY Writers"}, Admins: Rule{Type: "ImplicitMeta", Rule: "MAJORITY Admins"}, BlockValidation: Rule{Type: "ImplicitMeta", Rule: "ANY Writers"}},
	}
	viper.Set("Orderer", orderDefaults)

	channelDefaults := ChannelDefaults{
		Policies:     Policies{Readers: Rule{Type: "ImplicitMeta", Rule: "ANY Readers"}, Writers: Rule{Type: "ImplicitMeta", Rule: "ANY Writers"}, Admins: Rule{Type: "ImplicitMeta", Rule: "MAJORITY Admins"}, BlockValidation: Rule{Type: "ImplicitMeta", Rule: "ANY Writers"}},
		Capabilities: map[string]map[string]bool{"<<": {"V2_0": true}},
	}
	viper.Set("Channel", channelDefaults)

	type Consenters struct {
		Host          string `yaml:"Host"`
		Port          string `yaml:"Port"`
		ClientTLSCert string `yaml:"ClientTLSCert"`
		ServerTLSCert string `yaml:"ServerTLSCert"`
	}

	consenters := []Consenters{{Host: fmt.Sprintf("orderer.%v", domainName), Port: "7050", ClientTLSCert: fmt.Sprintf("../organizations/ordererOrganizations/%v/orderers/orderer.%v/tls/server.crt", domainName, domainName), ServerTLSCert: fmt.Sprintf("../organizations/ordererOrganizations/%v/orderers/orderer.%v/tls/server.crt", domainName, domainName)}}

	profiles := map[string]map[string]interface{}{
		"ThreeOrgsChannel": {"<<": channelDefaults, "orderer": map[string]interface{}{"<<": orderDefaults, "OrdererType": "etcdraft", "EtcdRaft": map[string]interface{}{"Consenters": consenters}, "Organizations": ordererOrg, "Capabilities": Capabilities["Orderer"]}, "Application": map[string]interface{}{"<<": applicationDefaults, "Organizations": orgs[1:], "Capabilities": Capabilities["Application"]}}}

	// set profile config
	viper.Set("Profiles", profiles)

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
	// fmt.Println("this is info", info)
	// printNetworkInfo(info)
}

func CreateFolders(domainName string) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}
	folder1 := fmt.Sprintf("fabrix/%v/Network/config", domainName)
	folder2 := fmt.Sprintf("fabrix/%v/Network/docker", domainName)
	folder3 := fmt.Sprintf("fabrix/%v/Network/peercfg", domainName)
	// Create the full path for the new folder
	folderPath1 := filepath.Join(cwd, folder1)
	folderPath2 := filepath.Join(cwd, folder2)
	folderPath3 := filepath.Join(cwd, folder3)

	// Create the folder
	err = os.MkdirAll(folderPath1, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}
	// Create the folder
	err = os.MkdirAll(folderPath2, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}

	// Create the folder
	err = os.MkdirAll(folderPath3, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating folder:", err)
		return
	}

	fmt.Println("Folder created at:", domainName)
}

func CreateRegisterEnroll(domainName string, orgPeers map[string]int) {

	// Define the script file path
	filePath := fmt.Sprintf("./fabrix/%v/Network/registerEnroll.sh", domainName)

	ordererCaPort := info.Orderer.Ca.Port

	scriptContent0 := fmt.Sprintf(`
#!/bin/bash

function createOrderer() {
  export DOMAIN_NAME="%s"
  export ORDERER_CA_PORT=%d
  echo "Enrolling the CA admin"
  mkdir -p organizations/ordererOrganizations/${DOMAIN_NAME}

  export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}

  set -x
  fabric-ca-client enroll -u https://admin:adminpw@localhost:${ORDERER_CA_PORT} --caname ca-orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-${ORDERER_CA_PORT}-ca-orderer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-${ORDERER_CA_PORT}-ca-orderer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-${ORDERER_CA_PORT}-ca-orderer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-${ORDERER_CA_PORT}-ca-orderer.pem
    OrganizationalUnitIdentifier: orderer" > "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/msp/config.yaml"

  # Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

  # Copy orderer org's CA cert to orderer org's /msp/tlscacerts directory (for use in the channel MSP definition)
  mkdir -p "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/msp/tlscacerts"
  cp "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem" "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/msp/tlscacerts/tlsca.${DOMAIN_NAME}-cert.pem"

  # Copy orderer org's CA cert to orderer org's /tlsca directory (for use by clients)
  mkdir -p "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/tlsca"
  cp "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem" "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/tlsca/tlsca.${DOMAIN_NAME}-cert.pem"

  echo "Registering orderer"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Registering the orderer admin"
  set -x
  fabric-ca-client register --caname ca-orderer --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  echo "Generating the orderer msp"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:${ORDERER_CA_PORT} --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/msp/config.yaml"

  echo "Generating the orderer-tls certificates, use --csr.hosts to specify Subject Alternative Names"
  set -x
  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:${ORDERER_CA_PORT} --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/tls" --enrollment.profile tls --csr.hosts orderer.${DOMAIN_NAME} --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  # Copy the tls CA cert, server cert, server keystore to well known file names in the orderer's tls directory that are referenced by orderer startup config
  cp "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/tls/ca.crt"
  cp "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/tls/signcerts/"* "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/tls/server.crt"
  cp "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/tls/keystore/"* "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/tls/server.key"

  # Copy orderer org's CA cert to orderer's /msp/tlscacerts directory (for use in the orderer MSP definition)
  mkdir -p "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/msp/tlscacerts"
  cp "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/tls/tlscacerts/"* "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/msp/tlscacerts/tlsca.${DOMAIN_NAME}-cert.pem"

  echo "Generating the admin msp"
  set -x
  fabric-ca-client enroll -u https://ordererAdmin:ordererAdminpw@localhost:${ORDERER_CA_PORT} --caname ca-orderer -M "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/users/Admin@${DOMAIN_NAME}/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/ordererOrg/ca-cert.pem"
  { set +x; } 2>/dev/null

  cp "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/msp/config.yaml" "${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/users/Admin@${DOMAIN_NAME}/msp/config.yaml"
}
  createOrderer
  `, domainName, ordererCaPort)

	// Append the second content block to the script file
	if err := appendToScriptFile(scriptContent0, filePath); err != nil {
		fmt.Println("Error appending to script file:", err)
		return
	}

	for i, org := range keys {
		// Define dynamic values
		orgName := org
		caPort := info.Organisations[i].Ca.Port
		caser := cases.Title(language.English)
		orgCap := caser.String(org)

		// Prepare the script content with variables defined at the beginning
		scriptContent1 := fmt.Sprintf(`
function create%sCertificates(){
	# Define dynamic variables
	export ORG_NAME_DOMAIN="%s.%s"
	export ORG_NAME="%s"
	export CA_PORT=%d

	echo "Enrolling the CA admin"
	mkdir -p organizations/peerOrganizations/${ORG_NAME_DOMAIN}/

	export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/

	set -x
	fabric-ca-client enroll -u https://admin:adminpw@localhost:${CA_PORT} --caname ca-${ORG_NAME} --tls.certfiles "${PWD}/organizations/fabric-ca/${ORG_NAME}/ca-cert.pem"
	{ set +x; } 2>/dev/null

  echo "NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-${CA_PORT}-ca-${ORG_NAME}.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-${CA_PORT}-ca-${ORG_NAME}.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-${CA_PORT}-ca-${ORG_NAME}.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-${CA_PORT}-ca-${ORG_NAME}.pem
    OrganizationalUnitIdentifier: orderer" > "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/msp/config.yaml"

	# Since the CA serves as both the organization CA and TLS CA, copy the org's root cert that was generated by CA startup into the org level ca and tlsca directories

	# Copy ${ORG_NAME}'s CA cert to ${ORG_NAME}'s /msp/tlscacerts directory (for use in the channel MSP definition)
	mkdir -p "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/msp/tlscacerts"
	cp "${PWD}/organizations/fabric-ca/${ORG_NAME}/ca-cert.pem" "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/msp/tlscacerts/ca.crt"

	# Copy ${ORG_NAME}'s CA cert to ${ORG_NAME}'s /tlsca directory (for use by clients)
	mkdir -p "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/tlsca"
	cp "${PWD}/organizations/fabric-ca/${ORG_NAME}/ca-cert.pem" "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/tlsca/tlsca.${ORG_NAME_DOMAIN}-cert.pem"

	# Copy ${ORG_NAME}'s CA cert to ${ORG_NAME}'s /ca directory (for use by clients)
	mkdir -p "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/ca"
	cp "${PWD}/organizations/fabric-ca/${ORG_NAME}/ca-cert.pem" "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/ca/ca.${ORG_NAME_DOMAIN}-cert.pem"

	echo "Registering user"
	set -x
	fabric-ca-client register --caname ca-${ORG_NAME} --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/fabric-ca/${ORG_NAME}/ca-cert.pem"
	{ set +x; } 2>/dev/null

	echo "Registering the org admin"
	set -x
	fabric-ca-client register --caname ca-${ORG_NAME} --id.name ${ORG_NAME}admin --id.secret ${ORG_NAME}adminpw --id.type admin --tls.certfiles "${PWD}/organizations/fabric-ca/${ORG_NAME}/ca-cert.pem"
	{ set +x; } 2>/dev/null

	echo "Generating the user msp"
	set -x
	fabric-ca-client enroll -u https://user1:user1pw@localhost:${CA_PORT} --caname ca-${ORG_NAME} -M "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/users/User1@${ORG_NAME_DOMAIN}/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/${ORG_NAME}/ca-cert.pem"
	{ set +x; } 2>/dev/null

	cp "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/msp/config.yaml" "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/users/User1@${ORG_NAME_DOMAIN}/msp/config.yaml"

	echo "Generating the org admin msp"
	set -x
	fabric-ca-client enroll -u https://${ORG_NAME}admin:${ORG_NAME}adminpw@localhost:${CA_PORT} --caname ca-${ORG_NAME} -M "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/users/Admin@${ORG_NAME_DOMAIN}/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/${ORG_NAME}/ca-cert.pem"
	{ set +x; } 2>/dev/null

	cp "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/msp/config.yaml" "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/users/Admin@${ORG_NAME_DOMAIN}/msp/config.yaml"

`, orgCap, orgName, domainName, orgName, caPort)

		// Append the second content block to the script file
		if err := appendToScriptFile(scriptContent1, filePath); err != nil {
			fmt.Println("Error appending to script file:", err)
			return
		}

		for i := range info.Organisations[i].Peers {

			// Prepare the script content with variables defined at the beginning
			scriptContent2 := fmt.Sprintf(`

	# Define dynamic variables
	export PEER="peer%v"

	echo "Registering ${PEER}"
	set -x
	fabric-ca-client register --caname ca-${ORG_NAME} --id.name ${PEER} --id.secret ${PEER}pw --id.type peer --tls.certfiles "${PWD}/organizations/fabric-ca/${ORG_NAME}/ca-cert.pem"
	{ set +x; } 2>/dev/null

	echo "Generating the ${PEER} msp"
	set -x
	fabric-ca-client enroll -u https://${PEER}:${PEER}pw@localhost:${CA_PORT} --caname ca-${ORG_NAME} -M "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/peers/${PEER}.${ORG_NAME_DOMAIN}/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/${ORG_NAME}/ca-cert.pem"
	{ set +x; } 2>/dev/null

	cp "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/msp/config.yaml" "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/peers/${PEER}.${ORG_NAME_DOMAIN}/msp/config.yaml"

	echo "Generating the ${PEER}-tls certificates, use --csr.hosts to specify Subject Alternative Names"
	set -x
	fabric-ca-client enroll -u https://${PEER}:${PEER}pw@localhost:${CA_PORT} --caname ca-${ORG_NAME} -M "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/peers/${PEER}.${ORG_NAME_DOMAIN}/tls" --enrollment.profile tls --csr.hosts ${PEER}.${ORG_NAME_DOMAIN} --csr.hosts localhost --tls.certfiles "${PWD}/organizations/fabric-ca/${ORG_NAME}/ca-cert.pem"
	{ set +x; } 2>/dev/null

	# Copy the tls CA cert, server cert, server keystore to well known file names in the peer's tls directory that are referenced by peer startup config
	cp "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/peers/${PEER}.${ORG_NAME_DOMAIN}/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/peers/${PEER}.${ORG_NAME_DOMAIN}/tls/ca.crt"
	cp "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/peers/${PEER}.${ORG_NAME_DOMAIN}/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/peers/${PEER}.${ORG_NAME_DOMAIN}/tls/server.crt"
	cp "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/peers/${PEER}.${ORG_NAME_DOMAIN}/tls/keystore/"* "${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/peers/${PEER}.${ORG_NAME_DOMAIN}/tls/server.key"
		
		`, i)

			// Append the second content block to the script file
			if err := appendToScriptFile(scriptContent2, filePath); err != nil {
				fmt.Println("Error appending to script file:", err)
				return
			}
		}
		scriptContent3 := fmt.Sprintf(`
	}
	create%sCertificates`, orgCap)
		// Append the second content block to the script file
		if err := appendToScriptFile(scriptContent3, filePath); err != nil {
			fmt.Println("Error appending to script file:", err)
			return
		}
	}
	fmt.Println("Successfully created registerEnroll.sh")
}

func CreateStartNetwork(domainName string, orgPeers map[string]int) {
	caser := cases.Title(language.English)

	filePath := fmt.Sprintf("./fabrix/%v/Network/startNetwork.sh", domainName)
	scriptContent0 := fmt.Sprintf(`
	#!/bin/bash
	export DOMAIN_NAME="%s"

	echo "------------Register the ca admin for each organization—----------------"

docker compose -f docker/docker-compose-ca.yaml up -d
sleep 3

sudo chmod -R 777 organizations/

echo "------------Register and enroll the users for each organization—-----------"

chmod +x registerEnroll.sh

./registerEnroll.sh
sleep 3

echo "—-------------Build the infrastructure—-----------------"

docker compose -f docker/docker-compose-orgs.yaml up -d
sleep 3

echo "-------------Generate the genesis block—-------------------------------"

export FABRIC_CFG_PATH=${PWD}/config

export CHANNEL_NAME=fabrixchannel

configtxgen -profile ThreeOrgsChannel -outputBlock ${PWD}/channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
sleep 2

echo "------ Create the application channel------"

export ORDERER_CA=${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/msp/tlscacerts/tlsca.${DOMAIN_NAME}-cert.pem

export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/tls/server.crt

export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/tls/server.key

osnadmin channel join --channelID $CHANNEL_NAME --config-block ${PWD}/channel-artifacts/$CHANNEL_NAME.block -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2

osnadmin channel list -o localhost:7053 --ca-file $ORDERER_CA --client-cert $ORDERER_ADMIN_TLS_SIGN_CERT --client-key $ORDERER_ADMIN_TLS_PRIVATE_KEY
sleep 2

export FABRIC_CFG_PATH=${PWD}/peercfg
export CORE_PEER_TLS_ENABLED=true
	`, domainName)

	if err := appendToScriptFile(scriptContent0, filePath); err != nil {
		fmt.Println("Error appending to script file:", err)
		return
	}
	peerStringList := []string{}

	for i, org := range keys {
		org := strings.ToLower(org)
		orgMSP := fmt.Sprintf("%vMSP", caser.String(org))
		upperOrg := strings.ToUpper(org)
		scriptContent1 := fmt.Sprintf(`
	#Define dynamic variables
	export ORG_NAME_DOMAIN="%s.%s"
	export ORG_NAME="%s"
	export ORG_MSP="%s"
	`, org, domainName, org, orgMSP)

		if err :=     appendToScriptFile(scriptContent1, filePath); err != nil {
			fmt.Println("Error appending to script file:", err)
			return
		}

		for j := range info.Organisations[i].Peers {
			// for testing purpose now we are updating only one peer as anchor peer
			// and installing chaincode in that peer
			if j < 1 {
				peerPort := info.Organisations[i].Peers[j].Port
				scriptContent2 := fmt.Sprintf(`export PEER="peer%v" 
	export PEER_PORT=%d
	export ORG_CAP="%s"
	export CORE_PEER_LOCALMSPID=${ORG_MSP} 
	export CORE_PEER_ADDRESS=localhost:${PEER_PORT} 
	export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/peers/${PEER}.${ORG_NAME_DOMAIN}/tls/ca.crt
	export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/users/Admin@${ORG_NAME_DOMAIN}/msp
	export ${ORG_CAP}_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/peers/${PEER}.${ORG_NAME_DOMAIN}/tls/ca.crt

	echo "—---------------Join ${ORG_NAME} peer to the channel—-------------"

	peer channel join -b ${PWD}/channel-artifacts/$CHANNEL_NAME.block
	sleep 1
	peer channel list

	echo "—-------------${ORG_NAME} anchor peer update—-----------"

	peer channel fetch config ${PWD}/channel-artifacts/config_block.pb -o localhost:7050 --ordererTLSHostnameOverride orderer.${DOMAIN_NAME} -c $CHANNEL_NAME --tls --cafile $ORDERER_CA
	sleep 1

	cd channel-artifacts

	configtxlator proto_decode --input config_block.pb --type common.Block --output config_block.json
	jq '.data.data[0].payload.data.config' config_block.json > config.json
	cp config.json config_copy.json

	jq ".channel_group.groups.Application.groups.${ORG_MSP}.values += {AnchorPeers:{mod_policy: \"Admins\",\"value\":{\"anchor_peers\": [{\"host\": \"${PEER}.${ORG_NAME_DOMAIN}\",\"port\": ${PEER_PORT}}]},\"version\": \"0\"}}" config_copy.json > modified_config.json

	configtxlator proto_encode --input config.json --type common.Config --output config.pb
	configtxlator proto_encode --input modified_config.json --type common.Config --output modified_config.pb
	configtxlator compute_update --channel_id $CHANNEL_NAME --original config.pb --updated modified_config.pb --output config_update.pb

	configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate --output config_update.json
	echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . > config_update_in_envelope.json
	configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope --output config_update_in_envelope.pb

	cd ..

	peer channel update -f ${PWD}/channel-artifacts/config_update_in_envelope.pb -c $CHANNEL_NAME -o localhost:7050  --ordererTLSHostnameOverride orderer.${DOMAIN_NAME} --tls --cafile $ORDERER_CA
	sleep 1

	echo "—---------------package chaincode—-------------"

	peer lifecycle chaincode package chaincode.tar.gz --path ${PWD}/../Chaincode/ --lang golang --label chaincode_1.0
	sleep 1

	echo "—---------------install chaincode in ${ORG_NAME} peer—-------------"

	peer lifecycle chaincode install chaincode.tar.gz
	sleep 3

	peer lifecycle chaincode queryinstalled

	export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid chaincode.tar.gz)

	echo "—---------------Approve chaincode in ${ORG_NAME} peer—-------------"

	peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.${DOMAIN_NAME} --channelID $CHANNEL_NAME --name sample-chaincode --version 1.0 --collections-config ../Chaincode/collection.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
	sleep 1

		`, j, peerPort, upperOrg)

				if err := appendToScriptFile(scriptContent2, filePath); err != nil {
					fmt.Println("Error appending to script file:", err)
					return
				}
				// for committing chaincode need peer details
				peerString := fmt.Sprintf("--peerAddresses localhost:%d --tlsRootCertFiles $%s_PEER_TLSROOTCERT", peerPort, upperOrg)
				peerStringList = append(peerStringList, peerString)
			}

		}
	}
	peerStrings := strings.Join(peerStringList, "  ")

	scriptContent3 := fmt.Sprintf(`
	echo "—---------------Commit chaincode —-------------"
	peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name sample-chaincode --version 1.0 --sequence 1 --collections-config ../Chaincode/collection.json --tls --cafile $ORDERER_CA --output json
	peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.${DOMAIN_NAME} --channelID $CHANNEL_NAME --name sample-chaincode --version 1.0 --sequence 1 --collections-config ../Chaincode/collection.json --tls --cafile $ORDERER_CA %v
	sleep 1
	peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name sample-chaincode --cafile $ORDERER_CA

	`, peerStrings)

	if err := appendToScriptFile(scriptContent3, filePath); err != nil {
		fmt.Println("Error appending to script file:", err)
		return
	}
	fmt.Println("Successfully created startNetwork.sh")
}

func appendToScriptFile(content string, filePath string) error {
	// Open the file in append mode. Create the file if it doesn't exist.
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the content to the file
	if _, err := file.WriteString(content); err != nil {
		return err
	}

	// Make the script executable (optional)
	// err = os.Chmod(filePath, 0755)
	// if err != nil {
	// 	fmt.Println("Error setting executable permissions:", err)
	// 	return err
	// }

	return nil
}

func createStopNetwork(domainName string) {
	filePath := fmt.Sprintf("./fabrix/%v/Network/stopNetwork.sh", domainName)
	scriptContent0 :=
		`
	#!/bin/bash

	docker compose -f docker/docker-compose-orgs.yaml down
	sleep 2

	docker compose -f docker/docker-compose-ca.yaml down
	sleep 2

	docker rm -f $(docker ps -a | awk '($2 ~ /dev-peer.*/) {print $1}')

	docker volume rm $(docker volume ls -q)

	rm -rf channel-artifacts/
	rm chaincode.tar.gz
	rm -rf organizations/

	docker ps -a

	docker rm $(docker container ls -q) --force

	yes | docker container prune

	yes | docker system prune

	yes | docker volume prune

	yes | docker network prune

	`

	if err := appendToScriptFile(scriptContent0, filePath); err != nil {
		fmt.Println("Error appending to script file:", err)
		return
	}
	fmt.Println("Successfully created stopNetwork.sh")
}

func createPeercfg(domainName string) {

	filePath := fmt.Sprintf("./fabrix/%v/Network/peercfg/core.yaml", domainName)

	if err := os.MkdirAll(fmt.Sprintf("./fabrix/%v/Network/peercfg", domainName), os.ModePerm); err != nil {
		fmt.Printf("Error creating directories, %s\n", err)
		return
	}

	newFile, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file, %s\n", err)
		return
	}
	defer newFile.Close()

	_, err = newFile.Write(coreYaml)
	if err != nil {
		fmt.Printf("Error writing to file, %s\n", err)
		return
	}

	fmt.Println("successfully created core.yaml")
}

// these are for future works
// func CreateCertificates(orgPeers map[string]int) {
// 	for orgName, peerCount := range orgPeers {
// 		caName := fmt.Sprintf("ca-%s", orgName)
// 		adminName := "admin"
// 		adminPW := "adminpw"
// 		port := 7054
// 		tlsCertPath := "${PWD}/organizations/fabric-ca/" + orgName + "/ca-cert.pem"

// 		fmt.Printf("Processing organization: %s with %d peers\n", orgName, peerCount)

// 		// Commands with dynamic variables
// 		cmds := []string{
// 			fmt.Sprintf("echo 'Enrolling the CA admin for %s'", orgName),
// 			fmt.Sprintf("mkdir -p organizations/peerOrganizations/%s/", orgName),
// 			fmt.Sprintf("export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/%s/", orgName),
// 			fmt.Sprintf("fabric-ca-client enroll -u https://%s:%s@localhost:%d --caname %s --tls.certfiles \"%s\"", adminName, adminPW, port, caName, tlsCertPath),
// 		}

// 		for _, cmdStr := range cmds {
// 			cmd := exec.Command("bash", "-c", cmdStr)
// 			output, err := cmd.CombinedOutput()
// 			if err != nil {
// 				fmt.Printf("Error executing command: %v\n", err)
// 				fmt.Println(string(output))
// 				return
// 			}
// 			fmt.Println(string(output))
// 		}

// 		// Additional commands can be added to handle peers for each organization
// 		for i := 0; i < peerCount; i++ {
// 			peerName := fmt.Sprintf("peer%d.%s", i, orgName)
// 			fmt.Printf("Processing %s\n", peerName)

// 			// Example dynamic command for a peer
// 			peerCmd := fmt.Sprintf("fabric-ca-client register --caname %s --id.name %s --id.secret peer%dPW --id.type peer --tls.certfiles \"%s\"", caName, peerName, i, tlsCertPath)

// 			cmd := exec.Command("bash", "-c", peerCmd)
// 			output, err := cmd.CombinedOutput()
// 			if err != nil {
// 				fmt.Printf("Error executing command: %v\n", err)
// 				fmt.Println(string(output))
// 				return
// 			}
// 			fmt.Println(string(output))
// 		}
// 	}
// }

// func createPeercfg(domainName string) {
// 	filePath := fmt.Sprintf("./fabrix/%v/Network/peercfg/core.yaml", domainName)

// 	// Set the file you want to read
// 	viper.SetConfigName("core")
// 	viper.SetConfigType("yaml")
// 	viper.AddConfigPath("pkg/configs/defaults/peercfg")

// 	// Read the config file
// 	if err := viper.ReadInConfig(); err != nil {
// 		fmt.Printf("Error reading config file, %s", err)
// 		return
// 	}

// 	// Create a new file and write the modified content
// 	newFile, err := os.Create(filePath)
// 	if err != nil {
// 		fmt.Printf("Error creating file, %s", err)
// 		return
// 	}
// 	defer newFile.Close()

// 	// Write the new content
// 	if err := viper.WriteConfigAs(filePath); err != nil {
// 		fmt.Printf("Error writing config file, %s", err)
// 		return
// 	}

// 	fmt.Println("Config file has been successfully modified and written to new_config.yaml")
// }
