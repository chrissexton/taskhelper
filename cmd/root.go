package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// AddAction Flag to determine whether the user is adding a task with this template
var AddAction bool

var rootCmd = &cobra.Command{
	Use:   "taskhelper",
	Short: "taskwarrior helper",
	Long:  "This is a helper function for taskwarrior that allows you to set up task templates and generate reports",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var template string = args[0]

		if viper.IsSet(template) {
			runTemplate(template, shift(args))
			return
		}

		proxy(template)
	},
}

func runTemplate(template string, args []string) {
	if len(args) == 0 {
		report(template)
	}
	add(template, args)

}

func report(template string) {
	var reportName = template
	configName := template + ".report"
	if viper.IsSet(configName) {
		fmt.Println("Using the template as the report")
		reportName = viper.GetString(configName)
	}

	parameters := []string{reportName}
	execute("task", parameters)
}

func proxy(command string) {
	execute("task", []string{command})
}

func add(template string, task []string) {
	configName := template + ".add"
	if !viper.IsSet(configName) {
		err := fmt.Errorf("A Template wasn't found for `" + template + "`")
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Println("Using the template as the report")
	var addTemplate = viper.GetString(configName)

	parameters := append([]string{"add", addTemplate}, task...)

	fmt.Println("Adding a task using \"" + template + "\"")
	err := execute("task", parameters)
	fmt.Printf("Command finished with error: %v", err)
	return
}

func execute(command string, parameters []string) error {
	fmt.Println("Running `" + command + "` with \"" + strings.Join(parameters, " ") + "\"")
	cmd := getCmd(command, parameters)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	cmd.Wait()
	return err
}

func print(readCloser io.ReadCloser) {
	r := bufio.NewReader(readCloser)
	line, _, err := r.ReadLine()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println(string(line))
}

func getCmd(command string, parameters []string) *exec.Cmd {
	if len(parameters) == 1 {
		var parameter = parameters[0]
		return exec.Command(command, parameter)
	}
	return exec.Command(command, parameters...)
}

// Shift an element
func shift(s []string) []string {
	var shifted = s[1:]
	return shifted
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&AddAction, "add", "a", false, "Add a new task with the template")

	viper.SetConfigName("taskhelper")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/taskhelper")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

// Configuration Initialization
func initConfig() {
	viper.SetConfigName(".taskhelper")
	viper.AddConfigPath("$HOME/.config/.taskhelper")
	viper.AutomaticEnv()

	fmt.Println("Using config file:", viper.ConfigFileUsed())
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Execute the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
