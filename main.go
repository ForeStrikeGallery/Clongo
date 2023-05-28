package main

import (
	"bufio"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"encoding/json"
	"io/ioutil"
	"github.com/google/uuid"
	"os"
	"strings"
) 

func getUserInput() (map[string]interface{}, error) {
	var jsonMap map[string]interface{}

	fmt.Print("clongo> ")
	reader := bufio.NewReader(os.Stdin) 
	input, _ := reader.ReadString('\n')
	
	if !hasValidFunctionKeyword(input) {
		return jsonMap, fmt.Errorf("Invalid Function")
	}

	parts := strings.Split(input, " ") 
	jsonInput:= strings.Join(parts[1:], " ")	

	err := json.Unmarshal([]byte(jsonInput), &jsonMap)

	if err != nil {
		return jsonMap, fmt.Errorf("Invalid Json")
	}
	 	
	return jsonMap, nil	
}

func hasValidFunctionKeyword(input string) bool {	

	functionKeywords := []string {"insert", "find"}
	for _, keyword := range functionKeywords {
		if strings.HasPrefix(input, keyword) {
			return true
		}	
	}	
	
	return false 
}

func main() {


	for {
		userInput, err := getUserInput() 
		
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

	 	document := bson.M(userInput)	
		data, err := bson.Marshal(document)
	
		if err != nil {
			fmt.Println("Error marshalling BSON document: ", err)
			continue
		}

		filepath := "/Users/mmohan/Projects/Clongo/Data/"
		filename := uuid.New().String() + ".bson"

		err = ioutil.WriteFile(filepath + filename, data, 0644)

		if err != nil {	
			fmt.Println("Error saving BSON document: ", err) 
			continue
		} 
	}
} 
