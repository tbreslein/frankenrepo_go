/*
Copyright © 2023 Tommy Breslein <github.com/tbreslein>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type CommandType int

const (
	Build CommandType = iota
	Test
	Format
	Lint
	Run
)

const (
	ProjectManifestName = "frankenfest.toml"
	ConfigFile          = "frankenrepo.toml"
)

type Command struct {
	T         CommandType
	CmdString string
}

type Manifest struct {
	Version int
	Pkgs    []Pkg
}

type Pkg struct {
	Name string
	// cmds map[CommandType][]command
}

func ParseManifest(manifestDir string) Manifest {
	manifestPath := filepath.Join(manifestDir, ProjectManifestName)
	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		log.Fatalln(err)
	}

	manifest := Manifest{}
	err = toml.Unmarshal(manifestBytes, &manifest)
	var derr *toml.DecodeError
	if errors.As(err, &derr) {
		row, col := derr.Position()
		log.Fatal(
			"PARSING ERROR:: ",
			manifestPath, ":", row, ":", col,
			"\n", derr.String(),
		)
	}
	return manifest
}

func processPkgsArgs(args []string) []string {
	var packageList []string
	if len(args) == 0 {
		packageList = append(packageList, "all")
	} else {
		packageList = args
	}
	fmt.Print("Running frankenrepo on these packages: ")
	for _, a := range packageList {
		fmt.Print(a)
	}
	fmt.Println("")
	return packageList
}

// rootCmd represents the base command when called without any subcommands
var (
	cfgFile    string
	workingDir string
	rootCmd    = &cobra.Command{
		Use:   "frankenrepo",
		Short: "Tool to manage multi-language monorepos",
		Long: `Frankenrepo is a proof-of-concept tool to show that monorepo build tools can be
    language agnostic.`,

		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) {},
	}

	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "run build on package(s)",
		Run: func(cmd *cobra.Command, args []string) {
			packageList := processPkgsArgs(args)
			repo := ParseManifest(workingDir)

			fmt.Println(packageList)
			fmt.Println(repo)
		},
	}

	testCmd = &cobra.Command{
		Use:   "test",
		Short: "run test on package(s)",
		Run: func(cmd *cobra.Command, args []string) {
			packageList := processPkgsArgs(args)
			repo := ParseManifest(workingDir)

			fmt.Println(packageList)
			fmt.Println(repo)
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file"+
			" (default at $HOME/{,.config,.config/frankenrepo}/"+ConfigFile+")",
	)
	rootCmd.PersistentFlags().StringVarP(
		&workingDir,
		"working-directory",
		"C",
		".",
		"runs frankenrepo with this path as the CWD",
	)

	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(testCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(filepath.Join(home, ".config"))
		viper.AddConfigPath(filepath.Join(home, ".config", "frankenrepo"))
		viper.SetConfigType("toml")
		viper.SetConfigName("frankenrepo")
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	}
}
