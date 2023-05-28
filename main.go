package main

import (
	"bufio"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"encoding/json"
	"io/ioutil"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"strings"
) 

func validateUserInput(operation string, jsonInput string) (error) {

	if !hasValidFunctionKeyword(operation) {
		return fmt.Errorf("Invalid Operation")
	}
		
	var jsonMap map[string]interface{} 
	err := json.Unmarshal([]byte(jsonInput), &jsonMap)

	if err != nil {
		return fmt.Errorf("Invalid Json")
	}
	 	
	return nil	
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

func searchDocuments(rootDir string, searchTerm map[string]interface{}) ([]string, error) {
	var matchingDocuments []string 
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo,  err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		} 

		if filepath.Ext(path) == ".bson" {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
		
			var document map[string]interface{}

			err = bson.Unmarshal(data, &document) 
			if err != nil {
				return err
			}

			matching := true
			for key, value := range searchTerm {
				val, ok := document[key]  
				if !ok || val != value {
					matching = false
					break 
				}
			}

			if matching {
				jsonData, err := json.Marshal(document)
				if err != nil {
					return err 
				}
				
				jsonString := string(jsonData)	
				matchingDocuments = append(matchingDocuments, jsonString)	
			}	
		}
		
		return nil
	}) 

	if err != nil {
		return nil, err 
	}

	return matchingDocuments, nil
}

func insert(userInput map[string]interface{}) error {
	document := bson.M(userInput)	
	data, err := bson.Marshal(document)

	if err != nil {
		return fmt.Errorf("Error marshalling BSON document: ", err)
	}

	filepath := "/Users/mmohan/Projects/Clongo/Data/"
	filename := uuid.New().String() + ".bson"

	err = ioutil.WriteFile(filepath + filename, data, 0644)

	if err != nil {	
		return fmt.Errorf("Error saving BSON document: ", err) 
	}

	return nil 
} 


func main() {
	for {
		fmt.Print("clongo> ")
		reader := bufio.NewReader(os.Stdin) 
		input, _ := reader.ReadString('\n')
	
		parts := strings.Split(input, " ")
		operation := parts[0]
    	jsonInput:= strings.Join(parts[1:], " ")
			
		err := validateUserInput(operation, jsonInput) 

		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

	    var jsonMap map[string]interface{} 
		err = json.Unmarshal([]byte(jsonInput), &jsonMap)
	
		if operation == "insert" {
			err := insert(jsonMap)
				
			if err != nil {
				fmt.Println("Error while inserting: ", err) 
			} 
		}

		if operation == "find" {
			matchingDocuments, err := searchDocuments("/Users/mmohan/Projects/Clongo/Data/", jsonMap) 
			if err != nil {
				fmt.Println("Error while searching for documents: ", err) 
				continue
			}

			for _, document := range matchingDocuments {
				fmt.Println(document + "\n")
			}	
		}
	}
} 
