package command

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

const ARTIFACTS_FILE_NAME = "artifacts.json"

//go:embed schemas/artifact.json
var artifactSchema string

var validateArtifactsCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate different types of files",
	Long:  "Validate different types of files (artifacts, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		validateType, _ := cmd.Flags().GetString("type")
		if len(args) == 0 {
			fmt.Println("Error: Path argument is required")
			os.Exit(1)
		}

		path := args[0]
		if validateType == "artifacts" {
			if err := validateArtifacts(path); err != nil {
				fmt.Printf("Error validating artifacts: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Artifacts validation completed successfully")
		} else {
			fmt.Printf("Error: Unsupported validation type: %s\n", validateType)
			fmt.Println("Supported types: artifacts")
			os.Exit(1)
		}
	},
}

func validateArtifacts(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("invalid path: %v", err)
	}

	var hasErrors bool
	if fileInfo.IsDir() {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && info.Name() == ARTIFACTS_FILE_NAME {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return fmt.Errorf("error getting absolute path: %v", err)
				}
				if err := validateArtifactFile(path); err != nil {
					fmt.Printf("❌ %s - %v\n", absPath, err)
					hasErrors = true
				} else {
					fmt.Printf("✅ %s\n", absPath)
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		if hasErrors {
			return fmt.Errorf("some files failed validation")
		}
	} else {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("error getting absolute path: %v", err)
		}
		if err := validateArtifactFile(path); err != nil {
			fmt.Printf("❌ %s - %v\n", absPath, err)
			return err
		}
		fmt.Printf("✅ %s\n", absPath)
	}

	return nil
}

func validateArtifactFile(path string) error {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	var jsonData interface{}
	if err := json.Unmarshal(fileContent, &jsonData); err != nil {
		return fmt.Errorf("invalid JSON format: %v", err)
	}

	schemaLoader := gojsonschema.NewStringLoader(artifactSchema)
	documentLoader := gojsonschema.NewBytesLoader(fileContent)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	if !result.Valid() {
		var errors []string
		for _, desc := range result.Errors() {
			errors = append(errors, desc.String())
		}
		return fmt.Errorf("validation failed: %s", strings.Join(errors, ", "))
	}

	return nil
}

func init() {
	rootCmd.AddCommand(validateArtifactsCmd)
	validateArtifactsCmd.Flags().String("type", "", "Type of validation to perform (artifact)")
}
