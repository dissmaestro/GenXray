package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Client struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Flow  string `json:"flow,omitempty"` // Optional field
}

type VlessInboundConfig struct {
	Clients []Client `json:"clients"`
	// Other vless settings...
}

type Config struct {
	Log       map[string]interface{} `json:"log"`
	Routing   map[string]interface{} `json:"routing"`
	Inbounds  []interface{}          `json:"inbounds"`
	Outbounds []interface{}          `json:"outbounds"`
}

func main() {
	// Getting name of user from cmd
	args := os.Args
	fmt.Println("You should write name of your key \n exmple = dyadaya@vasya \n symbol \"@\"is nessesary:", args)
	if len(args) == 2 {
		fmt.Println("You are my sweatty, DEAR ", args[1])
		if !strings.Contains(args[1], "@") {
			fmt.Println("\n ERROR: argument must contains symbol \"@\" \n ")
			return
		}
	} else {
		fmt.Println("\n ERROR: You should give only 1 argument \n ")
		return
	}

	// Open the JSON file
	configFile := "config.json"
	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	// Unmarshall the JSON data into a struct
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println("Cannot serialize data(unmarshal)")
		return
	}

	vlessInboundIndex := -1
	for i, inbound := range config.Inbounds {
		if inboundMap, ok := inbound.(map[string]interface{}); ok {
			if protocol, ok := inboundMap["protocol"].(string); ok && protocol == "vless" {
				vlessInboundIndex = i
				break
			}
		}
	}

	if vlessInboundIndex == -1 {
		fmt.Println("Vless inbound configuration not found")
		return
	}

	vlessInbound := config.Inbounds[vlessInboundIndex].(map[string]interface{})
	vlessSettings := vlessInbound["settings"].(map[string]interface{})
	clientList, ok := vlessSettings["clients"].([]interface{})
	if !ok {
		fmt.Println("Error retrieving client list")
		return
	}

	// GEnerate uuid

	cmd := exec.Command("/opt/xray/xray uuid")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("ERROR: do you have Xray?", err)
		return
	}
	// Create a new client
	newClient := Client{
		Id:    string(output),     // Replace with a unique ID
		Email: args[1],            // Replace with the client's email
		Flow:  "xtls-rprx-vision", // Optional flow value
	}

	// Append the new client to the list
	clientList = append(clientList, newClient)
	vlessSettings["clients"] = clientList

	// Update the vless inbound config in the main config
	config.Inbounds[vlessInboundIndex] = vlessInbound

	// Marshal the updated config back to JSON
	newData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Write the updated config to the file
	err = os.WriteFile(configFile, newData, 0644)
	if err != nil {
		fmt.Println("Error writing config file:", err)
		return
	}

	fmt.Println("New client added successfully! ")
	fmt.Printf("New client added successfully! \n uuid = %s \n ShortIds = bba4b98aea9b4c44 \n  ", string(output))
}
