/*

Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/

package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/upsaurav12/bootstrap/pkg/addons"
	"github.com/upsaurav12/bootstrap/pkg/framework"
	"github.com/upsaurav12/bootstrap/templates"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "command for creating a new project.",
	Long:  `command for creating a new project.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the project name is provided
		if len(args) < 1 {
			fmt.Fprintln(cmd.OutOrStdout(), "Error: project name is required")
			return
		}

		// Get the template flag value from the command context
		tmpl, _ := cmd.Flags().GetString("type")

		// Get the project name (first argument)
		dirName := args[0]

		// Create the new project
		createNewProject(dirName, projectRouter, tmpl, cmd.OutOrStdout())
	},
}

var projectType string
var projectPort string
var projectRouter string
var DBType string
var YAMLPath string
var Entitys string

type TemplateData struct {
	ModuleName    string
	PortName      string
	DBType        string
	Imports       string
	Start         string
	Entity        string
	ContextName   string
	ContextType   string
	Router        string
	Bind          string
	JSON          string
	LowerEntity   string
	OtherImports  string
	ApiGroup      func(entity string, get string, lowerentity string) string
	Get           string
	FullContext   string
	ToTheClient   string
	Response      string
	ImportHandler string
	ImportRouter  string
	ServiceName   string
	Image         string
	Environment   string
	Port          string
	Volume        string
	VolumeName    string
	DBName        string
	DBEnvPrefix   string
	Import        string
	Driver        string
	DSN           string
}

type TemplateJob struct {
	TemplateDir string
	DestDir     string
}

func init() {
	// Add the new command to the rootCmd
	rootCmd.AddCommand(newCmd)

	// Define the --template flag for this command
	newCmd.Flags().StringVar(&projectType, "type", "", "type of the project")
	newCmd.Flags().StringVar(&projectPort, "port", "", "port of the project")
	newCmd.Flags().StringVar(&projectRouter, "router", "", "router of the project")
	newCmd.Flags().StringVar(&DBType, "db", "", "data type of the project")
	newCmd.Flags().StringVar(&YAMLPath, "yaml", "", "yaml file path")
	newCmd.Flags().StringVar(&Entitys, "entity", "", "entity")
}

func buildTemplateData(projectName string,
	projectPort string,
	frameworkConfig framework.FrameworkConfig,
	dbConfig *addons.DbAddOneConfig, // optional
	entity string, uppercase string) TemplateData {

	data := TemplateData{
		ModuleName:    projectName,
		PortName:      projectPort,
		DBType:        DBType,
		Imports:       frameworkConfig.Imports,
		Start:         frameworkConfig.Start,
		ContextName:   frameworkConfig.ContextName,
		ContextType:   frameworkConfig.ContextType,
		Entity:        uppercase,
		Router:        frameworkConfig.Router,
		Bind:          frameworkConfig.Bind,
		JSON:          frameworkConfig.JSON,
		LowerEntity:   Entitys,
		OtherImports:  frameworkConfig.OtherImports,
		ApiGroup:      frameworkConfig.ApiGroup,
		Get:           frameworkConfig.Get,
		FullContext:   frameworkConfig.FullContext,
		ToTheClient:   frameworkConfig.ToTheClient,
		Response:      frameworkConfig.Response,
		ImportHandler: frameworkConfig.ImportHandler,
		ImportRouter:  frameworkConfig.ImportRouter,
	}

	if dbConfig != nil {
		data.ServiceName = dbConfig.ServiceName
		data.DBName = dbConfig.DBName
		data.DBEnvPrefix = dbConfig.DBEnvPrefix
		data.Port = dbConfig.Port
		data.DSN = dbConfig.DSN
		data.Driver = dbConfig.Driver
		data.Import = dbConfig.Import
		data.Image = dbConfig.Image
		data.Environment = dbConfig.Environment
		data.Volume = dbConfig.Volume
		data.VolumeName = dbConfig.VolumeName
	}

	return data
}

func returnUppercase() string {
	var uppercase string
	if len(Entitys) > 0 {
		runes := []rune(Entitys)
		runes[0] = unicode.ToUpper(runes[0])

		uppercase = string(runes)
	}
	return uppercase
}

func createNewProject(projectName, projectRouter, template string, out io.Writer) {
	os.Mkdir(projectName, 0755)
	os.MkdirAll(filepath.Join(projectName, "internal"), 0755)

	// Prepare configs
	frameworkConfig := framework.FrameworkRegistory[projectRouter]
	var dbConfig *addons.DbAddOneConfig

	if DBType != "" {
		cfg := addons.DbRegistory[DBType]
		dbConfig = &cfg
	}

	uppercase := returnUppercase()

	// Build TemplateData once
	data := buildTemplateData(
		projectName,
		projectPort,
		frameworkConfig,
		dbConfig,
		Entitys, uppercase)

	// Render templates
	renderTemplateDir("common", projectName, data)
	renderTemplateDir("rest/clean", projectName, data)

	if DBType != "" {
		renderTemplateDir("db/"+DBType, projectName, data)
		renderTemplateDir("db/database", filepath.Join(projectName, "internal/db"), data)

		fmt.Fprintf(out, "✓ Added database support for '%s'\n", DBType)
	}

	fmt.Fprintf(out, "✓ Created '%s' successfully\n", projectName)
}

func IsHidden(path string) (bool, error) {
	// Unix hidden check
	name := filepath.Base(path)
	return strings.HasPrefix(name, "."), nil
}

func renderTemplateDir(templatePath, destinationPath string, data TemplateData) error {
	return fs.WalkDir(templates.FS, templatePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(templatePath, path)

		if d.IsDir() {
			return os.MkdirAll(filepath.Join(destinationPath, relPath), 0755)
		}

		// Read template file
		content, err := templates.FS.ReadFile(path)
		if err != nil {
			return err
		}

		// Compute dynamic output filename
		fileName := strings.TrimSuffix(relPath, ".tmpl")

		if templatePath == "common" {
			base := filepath.Base(fileName)
			if base == "env" || base == "golang-ci.yml" {
				fileName = "." + fileName
			}
		}
		// fmt.Println("patg: ", path)

		// 🔥 dynamic entity replacement
		fileName = strings.ReplaceAll(fileName, "example", strings.ToLower(data.Entity))
		// fileName = strings.ReplaceAll(fileName, "Example", upperFirst(data.Entity))

		targetPath := filepath.Join(destinationPath, fileName)

		// Parse template
		tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
		if err != nil {
			return err
		}

		// Write output file
		outFile, err := os.Create(targetPath)
		// fmt.Println("targetPath : ", targetPath, data)
		if err != nil {
			return err
		}
		defer outFile.Close()

		return tmpl.Execute(outFile, data)
	})
}
