package cmd

import (
	"fmt"
	"os"
	"encoding/json"
	"io"
	"io/ioutil"
	"path"
	"errors"
	"strings"

	"github.com/iancoleman/strcase"
)
func checkIfFileExists(filePath string) bool {
	_, err := os.Stat(filePath)

	if err == nil {
  	return true
	}

	if strings.Contains(filePath, WEBSOCKET_CONFIG_FILE_PATH) && errors.Is(err, os.ErrNotExist) {
		fmt.Println("Project does not exist. Try running websocket-generator init")
	}

	return false
}

func ReadFile(fileName string) map[string]interface{} {
	jsonFile, errorJsonFile := os.Open(fileName)
	if errorJsonFile != nil {
		fmt.Println("Open command in ReadFile failed with the following error", errorJsonFile)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal(byteValue, &result)

	return result
}

func WriteJsonFile(fileName string, data map[string]interface{})  {
	converted := make(map[string]interface{}, len(data))
	for key, value := range data {
		key = strcase.ToLowerCamel(key)
		converted[key] = value
	}

	file, _ := json.MarshalIndent(converted, "", "  ")
	error := ioutil.WriteFile(fileName, file, 0644)

	if error != nil {
		fmt.Println("Write to file failed Error:", error)
	}
}

func WriteFile(filePath string, data []byte) {
	error := ioutil.WriteFile(filePath, data, 0644)

	if error != nil {
		fmt.Println("Write to file failed Error:", error)
	}
}

func CopyAndMoveFile(sourceFilePath, destinationFilePath string) error {
	var error error
	var sourceFile *os.File
	var destinationFile *os.File
	var sourceFileInfo os.FileInfo

	sourceFile, error = os.Open(sourceFilePath)

	if error != nil {
		fmt.Println("Open command in CopyAndMoveFile failed")
		return error
	}

	defer sourceFile.Close()

	destinationFile, error = os.Create(destinationFilePath)

	if error != nil {
		fmt.Println("Create command in CopyAndMoveFile failed")
		return error
	}

	defer destinationFile.Close()

	_, error = io.Copy(destinationFile, sourceFile);

	if error != nil {
		fmt.Println("Copy command in CopyAndMoveFile failed")
		return error
	}

	sourceFileInfo, error = os.Stat(sourceFilePath)

	if error != nil {
		fmt.Println("Stat command in CopyAndMoveFile failed")
		return error
	}

	return os.Chmod(destinationFilePath, sourceFileInfo.Mode())
}

func CopyAndMoveFolder(sourceFilePath string, destinationFilePath string) error {
	var error error
	var fileDirectories []os.FileInfo
	var sourceFileInfo os.FileInfo

	sourceFileInfo, error = os.Stat(sourceFilePath)

	if error != nil {
		fmt.Println("Stat command in CopyAndMoveFolder failed")
		return error
	}

	error = os.MkdirAll(destinationFilePath, sourceFileInfo.Mode())

	if error != nil {
		fmt.Println("MkdirAll command in CopyAndMoveFolder failed")
		return error
	}

	fileDirectories, error = ioutil.ReadDir(sourceFilePath)

	if error != nil {
		fmt.Println("ReadDir command in CopyAndMoveFolder failed")
		return error
	}

	for _, fileDirectory := range fileDirectories {
		fullSourceFilePath := path.Join(sourceFilePath, fileDirectory.Name())
		destinationFilePath := path.Join(destinationFilePath, fileDirectory.Name())

		if fileDirectory.IsDir() {
			error = CopyAndMoveFolder(fullSourceFilePath, destinationFilePath)
			if error != nil {
				fmt.Println("CopyAndMoveFolder failed")
				return error
			}
		} else {
			error = CopyAndMoveFile(fullSourceFilePath, destinationFilePath)
			if error != nil {
				fmt.Println("CopyAndMoveFile failed")
				return error
			}
		}
	}

	return nil
}
