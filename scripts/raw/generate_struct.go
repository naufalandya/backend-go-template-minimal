package raw

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Convert mock types to Go types
func mockTypeToGo(t string) string {
	switch strings.ToLower(t) { // case-insensitive check
	case "string":
		return "string"
	case "int32", "int64", "int":
		return "int32"
	case "bool":
		return "bool"
	case "float", "float32", "float64":
		return "float64"
	case "interface", "struct", "any":
		return "interface{}" // Handle multiple possible names
	default:
		return t // assume it's a custom type
	}
}

func addValidationTags(mockFieldType, fieldName string) string {
	jsonName := toSnakeCase(fieldName)
	switch mockFieldType {
	case "string":
		return fmt.Sprintf("  %s %s `json:\"%s\" validate:\"required,min=1,max=255\"`", fieldName, mockFieldType, jsonName)
	case "int32", "int64", "int":
		return fmt.Sprintf("  %s %s `json:\"%s\" validate:\"required,min=1\"`", fieldName, mockFieldType, jsonName)
	case "bool":
		return fmt.Sprintf("  %s %s `json:\"%s\" validate:\"required\"`", fieldName, mockFieldType, jsonName)
	case "float64":
		return fmt.Sprintf("  %s %s `json:\"%s\" validate:\"required,gte=0\"`", fieldName, mockFieldType, jsonName)
	case "interface{}":
		return fmt.Sprintf("  %s %s `json:\"%s\"`", fieldName, mockFieldType, jsonName)
	default:
		return fmt.Sprintf("  %s %s `json:\"%s\"`", fieldName, mockFieldType, jsonName)
	}
}

func generateStructsFrommock(mockFilePath string) ([]string, error) {
	file, err := os.Open(mockFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var structs []string
	var structLines []string
	var currentStruct string
	inMessage := false
	// Updated regex to handle your .txt format (adjust as needed)
	reField := regexp.MustCompile(`(\w+)\s+(\w+)\s*`) // Simpler pattern for .txt files

	// Parse the .txt file
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "message ") || strings.HasPrefix(line, "struct ") {
			inMessage = true
			structLines = nil

			parts := strings.Fields(line)
			if len(parts) >= 2 {
				currentStruct = parts[1]
			}
			continue
		}

		if inMessage {
			if line == "}" {
				inMessage = false
				// Format the struct and add to structs list
				structs = append(structs, fmt.Sprintf("type %s struct {\n%s\n}", cases.Title(language.English).String(currentStruct), strings.Join(structLines, "\n")))
				continue
			}
			match := reField.FindStringSubmatch(line)
			if len(match) >= 3 { // At least type and name
				fieldType := mockTypeToGo(match[1])
				fieldName := toCamelCase(match[2])
				structLines = append(structLines, addValidationTags(fieldType, fieldName))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return structs, nil
}

// Convert CamelCase or PascalCase to snake_case
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, rune(strings.ToLower(string(r))[0]))
	}
	return string(result)
}

// Process all .txt files in the given folder and generate Go files with structs
func processmockFiles(mockFolderPath string) error {
	// Walk through the mock folder to find all .txt files
	err := filepath.Walk(mockFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Process only .txt files
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			// Generate structs from the .txt file
			structs, err := generateStructsFrommock(path)
			if err != nil {
				return err
			}

			// Create the model folder if it doesn't exist
			modelFolderPath := filepath.Join(filepath.Dir(path), "../models")
			if err := os.MkdirAll(modelFolderPath, os.ModePerm); err != nil {
				return err
			}

			// Create Go file for this .txt file in the model folder
			goFilePath := filepath.Join(modelFolderPath, strings.TrimSuffix(info.Name(), ".txt")+".go")
			goFile, err := os.Create(goFilePath)
			if err != nil {
				return err
			}
			defer goFile.Close()

			// Write the structs to the Go file
			goFileContent := fmt.Sprintf("package model\n\n// Generated structs from %s\n\n", path)
			goFileContent += strings.Join(structs, "\n\n")
			_, err = goFile.WriteString(goFileContent)
			if err != nil {
				return err
			}

			fmt.Printf("Generated: %s\n", goFilePath)
		}
		return nil
	})

	return err
}

func toCamelCase(str string) string {
	parts := strings.Split(str, "_")
	for i := range parts {
		parts[i] = cases.Title(language.English).String(parts[i])
	}
	return strings.Join(parts, "")
}

func main() {
	// Get the current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	// Define the txt folder path relative to the current working directory
	mockFolderPath := filepath.Join(workingDir, "mock")

	// Process all .txt files and generate Go files with structs
	err = processmockFiles(mockFolderPath)
	if err != nil {
		fmt.Println("Error processing .txt files:", err)
		return
	}

	fmt.Println("Generation complete!")
}
