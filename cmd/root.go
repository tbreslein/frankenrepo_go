/*
Copyright © 2023 Tommy Breslein <github.com/tbreslein>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type CommandType int

const (
	Build CommandType = iota
	Test
)

type Command struct {
	T         CommandType
	CmdString string
}

type Repo struct {
	Pkgs []Pkg
}

type Pkg struct {
	Name string
	// cmds map[CommandType][]command
}

func ParseManifest(manifestDir string) Repo {
	manifestBytes, err := os.ReadFile(filepath.Join(manifestDir, "frankenfest.yaml"))
	if err != nil {
		log.Fatalln(err)
	}

	repo := Repo{}
	err = yaml.Unmarshal(manifestBytes, &repo)
	if err != nil {
		log.Fatalln(err)
	}
	return repo
}

func processPkgsArgs(args []string) []string {
	var packageList []string
	if len(args) == 0 {
		packageList = append(packageList, "all")
	} else {
		packageList = args
	}
	for _, a := range packageList {
		fmt.Println(a)
	}
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
		"config file (default at $HOME/{,.config,.config/frankenrepo}.frankenrepo.yaml)",
	)
	rootCmd.PersistentFlags().StringVarP(
		&workingDir,
		"working directory",
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
		viper.SetConfigType("yaml")
		viper.SetConfigName("frankenrepo")
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	} else {
		println("no config file found")
	}
}
