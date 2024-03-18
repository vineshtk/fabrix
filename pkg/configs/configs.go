package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func CreateConfigs(orgPeers map[string]int) {
	// ReadConfig()
	// WriteDockerCa()
	CreateDockerComposeCA(orgPeers)
	// CreateCA(orgPeers)
}

func WriteDockerCa() {
	viper.Set("someKey", "newValue")

	// Marshal the configuration back to YAML
	configContent, err := yaml.Marshal(viper.AllSettings())
	if err != nil {
		fmt.Println("Error marshaling config to YAML", err)
		return
	}

	// Define the new file path
	newPath := "pkg/configs/generated/docker/docker-compose-ca.yaml"

	// Write the modified configuration to a new file
	err = os.WriteFile(newPath, configContent, 0644) // Adjust permissions as needed
	if err != nil {
		fmt.Println("Error writing config to new file", err)
		return
	}

	fmt.Println("Configuration written to new file successfully.")
}

func ReadConfig() {
	// Set the file name of the configurations file
	viper.SetConfigName("docker-compose-ca-default") // name of config file (without extension)

	// Set the type of the configuration file
	viper.SetConfigType("yaml")

	// Set the path to look for the configurations file
	viper.AddConfigPath("pkg/configs/defaults/docker") // path to look for the config file in

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		log.Fatalf("Error while reading config file %s", err)
	}
}

func CreateCA(orgPeers map[string]int) {

	// Reading from default configs

	// Setting up some configurations
	viper.Set("version", "3.7")
	viper.Set("networks.test.network", "fabric_test")
	viper.Set("website", "golang.org")
	viper.Set("services.ca_org1", "org1")

	// Creating and updating a configuration file
	viper.SetConfigName("docker-compose-ca")            // name of config file (without extension)
	viper.SetConfigType("yaml")                         // specifying the config type
	viper.AddConfigPath("pkg/configs/generated/docker") // path to look for the config file in

	// viper.SetConfigFile("pkg/configs/generated/docker/docker-compose-ca.yaml")
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

	fmt.Println("Configuration file created/updated successfully!")
}

func CreateDockerComposeCA(orgPeers map[string]int) {

	viper.SetConfigName("docker-compose-ca")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("pkg/configs/generated/docker")

	// set values for each fields in docker compose file
	viper.Set("version", "3.7")
	viper.Set("networks.test.network", "fabric_test")

	i := 0
	ports := []int{7054, 17054}

	for org := range orgPeers {
		viper.Set(fmt.Sprintf("services.ca_%v.image", org), "hyperledger/fabric-ca:1.5.7")
		viper.Set(fmt.Sprintf("services.ca_%v.labels.service", org), "hyperledger-fabric")

		// environment := map[string]interface{}{
		// 	"FABRIC_CA_HOME":                            "/etc/hyperledger/fabric-ca-server",
		// 	"FABRIC_CA_SERVER_CA_NAME":                  fmt.Sprintf("ca-%v", org),
		// 	"FABRIC_CA_SERVER_TLS_ENABLED":              "true",
		// 	"FABRIC_CA_SERVER_PORT":                     fmt.Sprint(ports[0] + i*1000),
		// 	"FABRIC_CA_SERVER_OPERATIONS_LISTENADDRESS": fmt.Sprintf("0.0.0.0:%v", ports[1]+i*1000),
		// }

		envSlice := []string{
			"FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server",
			fmt.Sprintf("FABRIC_CA_SERVER_CA_NAME=ca-%v", org),
			"FABRIC_CA_SERVER_TLS_ENABLED=true",
			fmt.Sprintf("FABRIC_CA_SERVER_PORT=%v", ports[0]+i*1000),
			fmt.Sprintf("FABRIC_CA_SERVER_OPERATIONS_LISTENADDRESS=0.0.0.0:%v", ports[1]+i*1000),
		}

		// Set environment variables as a slice of strings in the format "KEY=value"
		// var envSlice []string
		// for key, value := range environment {
		// 	envSlice = append(envSlice, key+"="+value.(string))
		// }

		viper.Set(fmt.Sprintf("services.ca_%v.environment", org), envSlice)

		portSlice := []string{
			fmt.Sprintf("%v:%v", ports[0]+i*1000, ports[0]+i*1000),
			fmt.Sprintf("%v:%v", ports[1]+i*1000, ports[1]+i*1000),
		}

		viper.Set(fmt.Sprintf("services.ca_%v.ports", org), portSlice)
		volumeSlice := [1]string{
			fmt.Sprintf("../organizations/fabric-ca/%v:/etc/hyperledger/fabric-ca-server", org),
		}
		viper.Set(fmt.Sprintf("services.ca_%v.command", org), "sh -c 'fabric-ca-server start -b admin:adminpw -d'")
		viper.Set(fmt.Sprintf("services.ca_%v.volumes", org), fmt.Sprintf("../organizations/fabric-ca/%v:/etc/hyperledger/fabric-ca-server", org))
		viper.Set(fmt.Sprintf("services.ca_%v.volumes", org), volumeSlice)

		viper.Set(fmt.Sprintf("services.ca_%v.container_name", org), fmt.Sprintf("ca_%v", org))
		viper.Set(fmt.Sprintf("services.ca_%v.networks", org), "test")

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
	fmt.Println("Configuration file created/updated successfully!")
}
