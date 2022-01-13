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
func checkIfFileExists(filePath string) (bool, error) {
	_, error := os.Stat(filePath)

	if error == nil {
		return true, nil
	}

	if strings.Contains(filePath, WEBSOCKET_CONFIG_FILE_PATH) && errors.Is(error, os.ErrNotExist) {
		return false, errors.New("Error checkIfFileExists: Project does not exist. Try running websocket-generator init")
	}

	if errors.Is(error, os.ErrNotExist) {
		return false, nil
	}

	return false, error
}

func ReadFile(fileName string) (map[string]interface{}, error) {
	jsonFile, error := os.Open(fileName)
	if error != nil {
		return nil, fmt.Errorf("Error ReadFile: Open command failed with the following error", error)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal(byteValue, &result)

	return result, nil
}

func WriteJsonFile(fileName string, data map[string]interface{}) error {
	converted := make(map[string]interface{}, len(data))
	for key, value := range data {
		key = strcase.ToLowerCamel(key)
		converted[key] = value
	}

	file, _ := json.MarshalIndent(converted, "", "  ")
	error := ioutil.WriteFile(fileName, file, 0644)

	if error != nil {
		return fmt.Errorf("Error WriteJsonFile: Writing to JSON file failed", error)
	}

	return nil
}

func WriteFile(filePath string, data []byte) error {
	error := ioutil.WriteFile(filePath, data, 0644)

	if error != nil {
		return fmt.Errorf("Error WriteFile: Writing to file failed", error)
	}

	return nil
}

func CopyAndMoveFile(sourceFilePath, destinationFilePath string) error {
	var error error
	var sourceFile *os.File
	var destinationFile *os.File
	var sourceFileInfo os.FileInfo

	sourceFile, error = os.Open(sourceFilePath)

	if error != nil {
		return fmt.Errorf("Error CopyAndMoveFile: Open command failed", error)
	}

	defer sourceFile.Close()

	destinationFile, error = os.Create(destinationFilePath)

	if error != nil {
		return fmt.Errorf("Error CopyAndMoveFile: Create command failed", error)
	}

	defer destinationFile.Close()

	_, error = io.Copy(destinationFile, sourceFile);

	if error != nil {
		return fmt.Errorf("Error CopyAndMoveFile: Copy command failed", error)
	}

	sourceFileInfo, error = os.Stat(sourceFilePath)

	if error != nil {
		return fmt.Errorf("Error CopyAndMoveFile: Stat command failed", error)
	}

	return os.Chmod(destinationFilePath, sourceFileInfo.Mode())
}

func CopyAndMoveFolder(sourceFilePath string, destinationFilePath string) error {
	var error error
	var fileDirectories []os.FileInfo
	var sourceFileInfo os.FileInfo

	sourceFileInfo, error = os.Stat(sourceFilePath)

	if error != nil {
		return fmt.Errorf("Error CopyAndMoveFolder: Stat command failed", error)
	}

	error = os.MkdirAll(destinationFilePath, sourceFileInfo.Mode())

	if error != nil {
		return fmt.Errorf("Error CopyAndMoveFile: MkdirAll command failed", error)
	}

	fileDirectories, error = ioutil.ReadDir(sourceFilePath)

	if error != nil {
		return fmt.Errorf("Error CopyAndMoveFile: ReadDir command failed", error)
	}

	for _, fileDirectory := range fileDirectories {
		fullSourceFilePath := path.Join(sourceFilePath, fileDirectory.Name())
		destinationFilePath := path.Join(destinationFilePath, fileDirectory.Name())

		if fileDirectory.IsDir() {
			error = CopyAndMoveFolder(fullSourceFilePath, destinationFilePath)
			if error != nil {
				return fmt.Errorf("Error CopyAndMoveFolder: CopyAndMoveFolder command failed", error)
			}
		} else {
			error = CopyAndMoveFile(fullSourceFilePath, destinationFilePath)
			if error != nil {
				return fmt.Errorf("Error CopyAndMoveFile: CopyAndMoveFolder command failed", error)
			}
		}
	}

	return nil
}
