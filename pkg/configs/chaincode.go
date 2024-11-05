package configs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func CreateDeployChaincode(domainName string, orgPeers map[string]int, channelName string) {
	info.ChannelName = channelName

	caser := cases.Title(language.English)

	filePath := fmt.Sprintf("./fabrix/%v/Network/deployChaincode.sh", domainName)
	scriptContent0 := fmt.Sprintf(`
	#!/bin/bash
	export DOMAIN_NAME=%s

export CHANNEL_NAME=%s

export ORDERER_CA=${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/msp/tlscacerts/tlsca.${DOMAIN_NAME}-cert.pem

export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/tls/server.crt

export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/tls/server.key

export FABRIC_CFG_PATH=${PWD}/peercfg
export CORE_PEER_TLS_ENABLED=true
	`, domainName, channelName)

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

		if err := appendToScriptFile(scriptContent1, filePath); err != nil {
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

func InstallChaincode(ccPath string, ccLang string) {
	domainName := "auto.com"
	path := filepath.Join("fabrix", domainName, "Network")

	// Create the initial content of the script
	initialContent := fmt.Sprintf(`#!/bin/bash
echo 'Packaging Chaincode'
export DOMAIN_NAME=%s
export CHANNEL_NAME=%s
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/${DOMAIN_NAME}/orderers/orderer.${DOMAIN_NAME}/msp/tlscacerts/tlsca.${DOMAIN_NAME}-cert.pem
export FABRIC_CFG_PATH=${PWD}/peercfg
export CORE_PEER_TLS_ENABLED=true
peer lifecycle chaincode package chaincode.tar.gz --path ${PWD}/%s --lang %s --label chaincode_1.0
`, domainName, "mychannel", ccPath, ccLang)

	// Create a temporary file for the script
	tmpFile, err := os.CreateTemp("", "install_chaincode_*.sh")
	if err != nil {
		fmt.Printf("Error creating temporary script file: %s\n", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Write the initial content to the temp file
	_, err = tmpFile.Write([]byte(initialContent))
	if err != nil {
		fmt.Printf("Error writing initial content to temporary script file: %s\n", err)
		return
	}

	// Close and reopen the temp file in append mode
	tmpFile.Close()
	tmpFile, err = os.OpenFile(tmpFile.Name(), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error reopening temporary script file for appending: %s\n", err)
		return
	}
	defer tmpFile.Close()

	// Append additional commands as needed
	extraCommands := `
# Add more commands here as needed
echo "Additional configuration"
peer lifecycle chaincode install chaincode.tar.gz
`

	_, err = tmpFile.WriteString(extraCommands)
	if err != nil {
		fmt.Printf("Error writing additional commands to temporary script file: %s\n", err)
		return
	}

	fmt.Println("Script file created and updated successfully!")
	
	// Display file contents
	// fmt.Println("Script file contents:")
	// content, _ := os.ReadFile(tmpFile.Name())
	// fmt.Println(string(content))


	// Make the script executable
	err = os.Chmod(tmpFile.Name(), 0755)
	if err != nil {
		fmt.Printf("Error making script executable: %s\n", err)
		return
	}

	// Run the script
	cmd := exec.Command("bash", tmpFile.Name())
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error executing script: %v\n", err)
		return
	}

	fmt.Println("Script executed successfully!")

	// caser := cases.Title(language.English)
	// orgNum := viper.GetInt("numberOfOrganisations")
	// // Additional commands can be added to handle peers for each organization
	// for i := 0; i < orgNum; i++ {

	// 	peerName := fmt.Sprintf("peer%d.%s", i, orgName)
	// 	fmt.Printf("Processing %s\n", peerName)

	// }

	// peerStringList := []string{}

	// for i, org := range keys {
	// 	org := strings.ToLower(org)
	// 	orgMSP := fmt.Sprintf("%vMSP", caser.String(org))
	// 	upperOrg := strings.ToUpper(org)
	// 	scriptContent1 := fmt.Sprintf(`
	// #Define dynamic variables
	// export ORG_NAME_DOMAIN="%s.%s"
	// export ORG_NAME="%s"
	// export ORG_MSP="%s"
	// `, org, domainName, org, orgMSP)

	// 	if err := appendToScriptFile(scriptContent1, filePath); err != nil {
	// 		fmt.Println("Error appending to script file:", err)
	// 		return
	// 	}

	// 	for j := range info.Organisations[i].Peers {
	// 		// for testing purpose now we are updating only one peer as anchor peer
	// 		// and installing chaincode in that peer
	// 		if j < 1 {
	// 			peerPort := info.Organisations[i].Peers[j].Port
	// 			scriptContent2 := fmt.Sprintf(`export PEER="peer%v"
	// export PEER_PORT=%d
	// export ORG_CAP="%s"
	// export CORE_PEER_LOCALMSPID=${ORG_MSP}
	// export CORE_PEER_ADDRESS=localhost:${PEER_PORT}
	// export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/peers/${PEER}.${ORG_NAME_DOMAIN}/tls/ca.crt
	// export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/users/Admin@${ORG_NAME_DOMAIN}/msp
	// export ${ORG_CAP}_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/${ORG_NAME_DOMAIN}/peers/${PEER}.${ORG_NAME_DOMAIN}/tls/ca.crt

	// echo "—---------------package chaincode—-------------"

	// peer lifecycle chaincode package chaincode.tar.gz --path ${PWD}/../Chaincode/ --lang golang --label chaincode_1.0
	// sleep 1

	// echo "—---------------install chaincode in ${ORG_NAME} peer—-------------"

	// peer lifecycle chaincode install chaincode.tar.gz
	// sleep 3

	// peer lifecycle chaincode queryinstalled

	// export CC_PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid chaincode.tar.gz)

	// echo "—---------------Approve chaincode in ${ORG_NAME} peer—-------------"

	// peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.${DOMAIN_NAME} --channelID $CHANNEL_NAME --name sample-chaincode --version 1.0 --collections-config ../Chaincode/collection.json --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile $ORDERER_CA --waitForEvent
	// sleep 1

	// 	`, j, peerPort, upperOrg)

	// 			if err := appendToScriptFile(scriptContent2, filePath); err != nil {
	// 				fmt.Println("Error appending to script file:", err)
	// 				return
	// 			}
	// 			// for committing chaincode, need peer details
	// 			peerString := fmt.Sprintf("--peerAddresses localhost:%d --tlsRootCertFiles $%s_PEER_TLSROOTCERT", peerPort, upperOrg)
	// 			peerStringList = append(peerStringList, peerString)
	// 		}
	// 	}
	// }
	// peerStrings := strings.Join(peerStringList, "  ")

	// scriptContent3 := fmt.Sprintf(`
	// echo "—---------------Commit chaincode —-------------"
	// peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name sample-chaincode --version 1.0 --sequence 1 --collections-config ../Chaincode/collection.json --tls --cafile $ORDERER_CA --output json
	// peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.${DOMAIN_NAME} --channelID $CHANNEL_NAME --name sample-chaincode --version 1.0 --sequence 1 --collections-config ../Chaincode/collection.json --tls --cafile $ORDERER_CA %v
	// sleep 1
	// peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name sample-chaincode --cafile $ORDERER_CA

	// `, peerStrings)

	// if err := appendToScriptFile(scriptContent3, filePath); err != nil {
	// 	fmt.Println("Error appending to script file:", err)
	// 	return
	// }
	// fmt.Println("Successfully created startNetwork.sh")

}
