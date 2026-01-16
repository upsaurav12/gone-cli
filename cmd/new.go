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
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
	"github.com/upsaurav12/bootstrap/pkg/addons"
	"github.com/upsaurav12/bootstrap/pkg/framework"
	"github.com/upsaurav12/bootstrap/pkg/parser"
	"github.com/upsaurav12/bootstrap/templates"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

type ProjectInput struct {
	Name     string
	Type     string
	Router   string
	Port     string
	DB       string
	Entities []string
}

var asciiStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(primaryBlue)).
	Bold(true)

const (
	stepName = iota
	stepType
	stepRouter
	stepPort
	stepDB
	stepConfirm
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(softBlue))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(softBlue))

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(mutedGray))

	boxStyle = lipgloss.NewStyle().
			Padding(3, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderGray).
			Width(60)
)

func renderStep(m wizardModel, title, label, body, hint string) string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleStyle.Render(title),
		"",
		labelStyle.Render(label),
		"",
		body,
		"",
		hintStyle.Render(hint),
	)

	box := boxStyle.
		Width(m.width - 4).
		Height(m.height - 2).
		Render(content)

	return box
}

type wizardModel struct {
	step  int
	input ProjectInput

	text textinput.Model
	list list.Model
	quit bool

	width  int
	height int
}

type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

func initialWizardModel() wizardModel {
	return wizardModel{
		step: stepName,
		text: newTextInput(""),
	}
}

const (
	primaryBlue = lipgloss.Color("33")  // bright blue
	softBlue    = lipgloss.Color("39")  // lighter blue
	mutedGray   = lipgloss.Color("241") // hints
	borderGray  = lipgloss.Color("238") // borders
)

func newTextInput(placeholder string) textinput.Model {
	ti := textinput.New()
	ti.Prompt = "› "
	ti.Placeholder = placeholder
	ti.SetValue("") // ← critical
	ti.Focus()
	return ti
}

func (m wizardModel) Init() tea.Cmd {
	return nil
}

func renderHeader() string {
	fig := figure.NewFigure("Bootstrap  CLI", "slant", true)

	ascii := strings.Trim(fig.String(), "\n")

	return asciiStyle.Render(ascii)
}

func (m wizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// ✅ HANDLE WINDOW SIZE FIRST
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	// ✅ HANDLE KEYS
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			m.quit = true
			return m, tea.Quit

		case "enter":
			switch m.step {

			case stepName:
				m.input.Name = m.text.Value()
				m.step = stepType
				m.list = newList("Project type", []string{"rest"})
				return m, nil

			case stepType:
				m.input.Type = m.list.SelectedItem().(item).Title()
				m.step = stepRouter
				m.list = newList("Router", []string{"gin", "chi", "echo"})
				return m, nil

			case stepRouter:
				m.input.Router = m.list.SelectedItem().(item).Title()
				m.step = stepPort
				m.text = newTextInput("")
				return m, nil

			case stepPort:
				port := m.text.Value()
				if port == "" {
					port = "8080"
				}
				if _, err := strconv.Atoi(port); err != nil {
					return m, nil
				}
				m.input.Port = port
				m.step = stepDB
				m.list = newList("Database", []string{"postgres", "mysql", "mongo"})
				return m, nil

			case stepDB:
				m.input.DB = m.list.SelectedItem().(item).Title()
				m.step = stepConfirm
				return m, nil

			case stepConfirm:
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	if m.step == stepName || m.step == stepPort {
		m.text, cmd = m.text.Update(msg)
		return m, cmd
	}

	if m.step == stepType || m.step == stepRouter || m.step == stepDB {
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m wizardModel) View() string {
	if m.quit {
		return ""
	}

	switch m.step {

	case stepName:
		return renderStep(
			m,
			renderHeader()+"\nCreate New Project",
			"Project name",
			m.text.View(),
			"Enter to continue • Esc to quit",
		)

	case stepType:
		return renderStep(
			m,
			"Project Type",
			"Select project type",
			m.list.View(),
			"↑↓ navigate • Enter select • Esc quit",
		)

	case stepRouter:
		return renderStep(
			m,
			"Router",
			"Select router",
			m.list.View(),
			"↑↓ navigate • Enter select • Esc quit",
		)

	case stepPort:
		return renderStep(
			m,
			"Application Port",
			"Port (default: 8080)",
			m.text.View(),
			"Enter to continue • Esc quit",
		)

	case stepDB:
		return renderStep(
			m,
			"Database",
			"Select database",
			m.list.View(),
			"↑↓ navigate • Enter select • Esc quit",
		)

	case stepConfirm:
		summary := fmt.Sprintf(
			"Project:  %s\nType:     %s\nRouter:   %s\nPort:     %s\nDatabase: %s",
			m.input.Name,
			m.input.Type,
			m.input.Router,
			m.input.Port,
			m.input.DB,
		)

		return renderStep(
			m,
			"Confirm Configuration",
			"Review your selections",
			summary,
			"Enter to generate • Esc to cancel",
		)
	}

	return ""
}

func newList(title string, values []string) list.Model {
	items := make([]list.Item, len(values))
	for i, v := range values {
		items[i] = item(v)
	}

	l := list.New(items, list.NewDefaultDelegate(), 20, 10)
	l.Title = title
	return l
}

func copyProjectYAML(srcPath, destDir string) error {
	if srcPath == "" {
		return nil // nothing to copy
	}

	content, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	destPath := filepath.Join(destDir, "project.yaml")

	return os.WriteFile(destPath, content, 0644)
}

func RunInteractiveWizard() (*ProjectInput, error) {
	p := tea.NewProgram(initialWizardModel())
	model, err := p.Run()
	if err != nil {
		return nil, err
	}

	m := model.(wizardModel)
	if m.quit {
		return nil, fmt.Errorf("aborted")
	}

	return &m.input, nil
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "command for creating a new project.",
	Long:  `command for creating a new project.`,
	Run: func(cmd *cobra.Command, args []string) {
		interactive, _ := cmd.Flags().GetBool("interactive")

		if interactive || (len(args) == 0 && YAMLPath == "") {
			input, err := RunInteractiveWizard()
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}

			projectType = input.Type
			projectRouter = input.Router
			projectPort = input.Port
			DBType = input.DB
			Entities = input.Entities

			createNewProject(input.Name, input.Router, input.Type, cmd.OutOrStdout())
			return
		}

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
	newCmd.Flags().Bool("interactive", false, "run interactive project setup")

}

func buildTemplateData(projectName string,
	projectPort string,
	frameworkConfig framework.FrameworkConfig,
	dbConfig *addons.DbAddOneConfig, // optional
	yamlConfig *parser.Config, uppercase []string) TemplateData {

	data := TemplateData{
		Name:          frameworkConfig.Name,
		ModuleName:    projectName,
		PortName:      projectPort,
		DBType:        DBType,
		Imports:       frameworkConfig.Imports,
		Start:         frameworkConfig.Start,
		ContextName:   frameworkConfig.ContextName,
		ContextType:   frameworkConfig.ContextType,
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
			Name:          frameworkConfig.Name,
			ModuleName:    projectName,
			PortName:      projectPort,
			DBType:        yamlConfig.Project.Database,
			Imports:       frameworkConfig.Imports,
			Start:         frameworkConfig.Start,
			ContextName:   frameworkConfig.ContextName,
			ContextType:   frameworkConfig.ContextType,
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

	// ✅ COPY project.yaml if provided
	if err := copyProjectYAML(YAMLPath, projectName); err != nil {
		fmt.Fprintf(out, "warning: could not copy project.yaml: %v\n", err)
	}

	// Prepare configs

	var frameworkConfig framework.FrameworkConfig
	frameworkConfig = framework.FrameworkRegistory[projectRouter]

	var dbConfig *addons.DbAddOneConfig

	jobs := []TemplateJob{
		{"common", projectName},
		{"rest/clean", projectName},
	}

	var data TemplateData

	var yamlConfig *parser.Config

	if YAMLPath != "" {
		var err error
		yamlConfig, err = parser.ReadYAML(YAMLPath)
		if err != nil {
			fmt.Printf("error while reading yaml file: %s", err)
			return
		}
	}

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
	newFile := strings.Replace(
		fileName,
		"example",
		"user",
		1, // only first replacement
	)

	entityData := data
	entityData.Entity = strings.Title("user")
	entityData.LowerEntity = strings.ToLower("user")
	targetPath := filepath.Join(destinationPath, newFile)

	tmpl, err := template.New(filepath.Base(tmpltPath)).Parse(string(content))
	if err != nil {
		return err
	}

	outFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return tmpl.Execute(outFile, entityData)
}
