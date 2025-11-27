/*

Copyright © 2025 Saurav Upadhyay sauravup041103@gmail.com

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

	"github.com/spf13/cobra"
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

func init() {
	// Add the new command to the rootCmd
	rootCmd.AddCommand(newCmd)

	// Define the --template flag for this command
	newCmd.Flags().StringVar(&projectType, "type", "", "type of the project")
	newCmd.Flags().StringVar(&projectPort, "port", "", "port of the project")
	newCmd.Flags().StringVar(&projectRouter, "router", "", "router of the project")
	newCmd.Flags().StringVar(&DBType, "db", "", "data type of the project")
}

func createNewProject(projectName string, projectRouter string, template string, out io.Writer) {
	err := os.Mkdir(projectName, 0755)
	if err != nil {
		fmt.Fprintf(out, "Error creating directory: %v\n", err)
		return
	}

	renderTemplateDir("common", projectName, TemplateData{
		ModuleName: projectName,
		PortName:   projectPort,
	})

	renderTemplateDir("rest"+"/"+projectRouter, projectName, TemplateData{
		ModuleName: projectName,
		PortName:   projectPort,
		DBType:     DBType,
	})

	if err != nil {
		fmt.Fprintf(out, "Error rendering templates: %v\n", err)
		return
	}

	if DBType != "" {
		dbTemplatePath := "db/" + DBType
		err := renderTemplateDir(dbTemplatePath, filepath.Join(projectName, "internal", "db"), TemplateData{
			ModuleName: projectName,
			PortName:   projectPort,
			DBType:     DBType,
		})
		if err != nil {

			fmt.Fprintf(out, "Error rendering  DB templates: %v\n", err)
			return
		}

		fmt.Fprintf(out, "✓ Added database support for '%s'\n", DBType)
	}

	fmt.Fprintf(out, "✓ Created '%s' successfully\n", projectName)
}

type TemplateData struct {
	ModuleName string
	PortName   string
	DBType     string
}

func renderTemplateDir(templatePath, destinationPath string, data TemplateData) error {
	return fs.WalkDir(templates.FS, templatePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Compute relative path (remove the base templatePath)
		relPath, _ := filepath.Rel(templatePath, path)
		targetPath := filepath.Join(destinationPath, strings.TrimSuffix(relPath, ".tmpl"))

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// ✅ Important: use full `path` for ReadFile
		content, err := templates.FS.ReadFile(path)
		if err != nil {
			return err
		}

		// Parse template
		tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
		if err != nil {
			return err
		}

		fmt.Println("path for tmpl: ", filepath.Base(path))

		// Write file
		outFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer outFile.Close()

		return tmpl.Execute(outFile, data)
	})
}
