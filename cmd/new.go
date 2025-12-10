/*

Copyright © 2025 Saurav Upadhyay sauravup041103@gmail.com

*/

package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/upsaurav12/bootstrap/pkg/addons"
	"github.com/upsaurav12/bootstrap/pkg/framework"
	"github.com/upsaurav12/bootstrap/pkg/parser"
	"github.com/upsaurav12/bootstrap/templates"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "command for creating a new project.",
	Long:  `command for creating a new project.`,
	Run: func(cmd *cobra.Command, args []string) {

		var dirName string
		if len(args) < 1 && YAMLPath != "" {
			yamlConfig, err := parser.ReadYAML(YAMLPath)
			if err != nil {
				fmt.Println("error while creating project using yaml: ", err)
				return
			}

			dirName = yamlConfig.Project.Name
		} else {
			dirName = args[0]
		}

		// Check if the project name is provided
		if len(args) < 1 && YAMLPath == "" {
			fmt.Fprintln(cmd.OutOrStdout(), "Error: project name is required")
			return
		}
		// Get the template flag value from the command context
		tmpl, _ := cmd.Flags().GetString("type")

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
var Entities []string
var yamlFile string

type TemplateData struct {
	Name          string
	ModuleName    string
	PortName      string
	DBType        string
	Imports       string
	Start         string
	Entity        string
	Entities      []string
	ContextName   string
	ContextType   string
	Router        string
	Bind          string
	JSON          string
	LowerEntity   string
	OtherImports  string
	UpperEntity   []string
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
	Returnable    string
	ReturnKeyword string
	HTTPHandler   string
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
	newCmd.Flags().StringSliceVar(&Entities, "entities", nil, "different entities")
}

func buildTemplateData(projectName string,
	projectPort string,
	frameworkConfig framework.FrameworkConfig,
	dbConfig *addons.DbAddOneConfig, // optional
	yamlConfig *parser.Config, uppercase []string) TemplateData {

	data := TemplateData{
		Name:        frameworkConfig.Name,
		ModuleName:  projectName,
		PortName:    projectPort,
		DBType:      DBType,
		Imports:     frameworkConfig.Imports,
		Start:       frameworkConfig.Start,
		ContextName: frameworkConfig.ContextName,
		ContextType: frameworkConfig.ContextType,
		// Entity:        uppercase,
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
		Returnable:    frameworkConfig.Returnable,
		ReturnKeyword: frameworkConfig.ReturnKeyword,
		HTTPHandler:   frameworkConfig.HTTPHandler,
		Entities:      Entities,
		// UpperEntity: ,
	}

	if yamlConfig != nil {
		data = TemplateData{
			Name:        frameworkConfig.Name,
			ModuleName:  projectName,
			PortName:    projectPort,
			DBType:      yamlConfig.Project.Database,
			Imports:     frameworkConfig.Imports,
			Start:       frameworkConfig.Start,
			ContextName: frameworkConfig.ContextName,
			ContextType: frameworkConfig.ContextType,
			// Entity:        uppercase,
			Router:        frameworkConfig.Router,
			Bind:          frameworkConfig.Bind,
			JSON:          frameworkConfig.JSON,
			LowerEntity:   Entitys,
			UpperEntity:   nil,
			OtherImports:  frameworkConfig.OtherImports,
			ApiGroup:      frameworkConfig.ApiGroup,
			Get:           frameworkConfig.Get,
			FullContext:   frameworkConfig.FullContext,
			ToTheClient:   frameworkConfig.ToTheClient,
			Response:      frameworkConfig.Response,
			ImportHandler: frameworkConfig.ImportHandler,
			ImportRouter:  frameworkConfig.ImportRouter,
			Returnable:    frameworkConfig.Returnable,
			ReturnKeyword: frameworkConfig.ReturnKeyword,
			HTTPHandler:   frameworkConfig.HTTPHandler,
			Entities:      yamlConfig.Entities,
		}
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

func returnUppercase(entity string) string {
	if entity == "" {
		return ""
	}

	runes := []rune(entity)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func createNewProject(projectName, projectRouter, template string, out io.Writer) {
	err := os.Mkdir(projectName, 0755)
	if err != nil {
		log.Fatalf("error occured while creating a new project %s: ", err)
	}
	err = os.MkdirAll(filepath.Join(projectName, "internal"), 0755)

	if err != nil {
		log.Fatalf("error occured while creating a new project %s: ", projectName)
	}

	// Prepare configs

	var frameworkConfig framework.FrameworkConfig
	frameworkConfig = framework.FrameworkRegistory[projectRouter]

	var dbConfig *addons.DbAddOneConfig

	jobs := []TemplateJob{
		{"common", projectName},
		{"rest/clean", projectName},
	}

	// if DBType != "" {
	// 	cfg := addons.DbRegistory[DBType]
	// 	dbConfig = &cfg
	// }

	var data TemplateData

	yamlConfig, err := parser.ReadYAML(YAMLPath)
	if err != nil {
		fmt.Printf("error while reading yaml file: %s", err)
		return
	}

	// this is subjected to change as we would have more things to add in
	// in this project.
	if DBType != "" {
		cfg := addons.DbRegistory[DBType]
		dbConfig = &cfg
	}
	if DBType == "" && yamlConfig != nil {
		DBType = yamlConfig.Project.Database
		cfg := addons.DbRegistory[DBType]
		dbConfig = &cfg
	}

	if yamlConfig != nil {

		frameworkConfig = framework.FrameworkRegistory[yamlConfig.Project.Router]
		Entities = yamlConfig.Entities

		if yamlConfig.Project.Database != "" {

			cfg := addons.DbRegistory[yamlConfig.Project.Database]
			dbConfig = &cfg

			DBType = yamlConfig.Project.Database

		}
	}

	var uppercase []string

	for _, entity := range Entities {
		u := returnUppercase(entity)

		uppercase = append(uppercase, u)
	}

	// data.UpperEntity = uppercase

	// till this line of code

	data = buildTemplateData(
		projectName,
		projectPort,
		frameworkConfig,
		dbConfig,
		yamlConfig, uppercase)

	data.UpperEntity = uppercase

	// Render templates
	_ = renderTemplateDir("common", projectName, data)
	_ = renderTemplateDir("rest/clean", projectName, data)

	if DBType != "" {
		jobs = append(jobs,
			TemplateJob{"db/" + DBType, projectName},
			TemplateJob{"db/database", filepath.Join(projectName, "internal", "db")},
		)

		// fmt.Fprintf(out, "✓ Added database support for '%s'\n", DBType)
	}

	for _, job := range jobs {
		if err := renderTemplateDir(job.TemplateDir, job.DestDir, data); err != nil {
			fmt.Fprintf(out, "Error rendering template %s → %s: %v\n",
				job.TemplateDir, job.DestDir, err)
			return
		}
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

		content, err := templates.FS.ReadFile(path)
		if err != nil {
			return err
		}

		fileName := strings.TrimSuffix(relPath, ".tmpl")

		if templatePath == "common" {
			base := filepath.Base(fileName)
			if base == "env" || base == "golang-ci.yml" {
				fileName = "." + fileName
			}
		}

		if len(data.Entities) == 0 {
			return writeSingle(data, fileName, path, content, destinationPath)
		}

		for _, entity := range data.Entities {

			// correct file renaming
			newFile := strings.Replace(
				fileName,
				"example",
				strings.ToLower(entity),
				1, // only first replacement
			)

			entityData := data
			entityData.Entity = strings.Title(entity)
			entityData.LowerEntity = strings.ToLower(entity)
			// capture errors!!
			if err := writeSingle(entityData, newFile, path, content, destinationPath); err != nil {
				return err
			}
		}

		return nil
	})
}

func writeSingle(data TemplateData, fileName string, tmpltPath string, content []byte, destinationPath string) error {
	targetPath := filepath.Join(destinationPath, fileName)

	tmpl, err := template.New(filepath.Base(tmpltPath)).Parse(string(content))
	if err != nil {
		return err
	}

	outFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return tmpl.Execute(outFile, data)
}
