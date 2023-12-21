/*
Copyright © 2023 Tommy Breslein <github.com/tbreslein>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tbreslein/frankenrepo/internal/repository"
)

type FrankenError string

const (
	CommandError FrankenError = "COMMAND_ERROR"
	UnknownError FrankenError = "UNKNOWN_ERROR"
)

var (
	cfgFile    string
	workingDir string
	rootCmd    = &cobra.Command{
		Use:   "frankenrepo",
		Short: "Tool to manage multi-language monorepos",
		Long: "Frankenrepo is a proof-of-concept tool to show that monorepo" +
			"build tools can be language agnostic.",
	}

	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "run build on package(s)",
		Run: func(cmd *cobra.Command, args []string) {
			// packageList := processPkgsArgs(args)
			repo := repository.InitRepo(args, workingDir)

			// fmt.Println(packageList)
			// fmt.Println(repo)
		},
	}

	// testCmd = &cobra.Command{
	// 	Use:   "test",
	// 	Short: "run test on package(s)",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		// TODO: extract processPkgsArgs and ParseFrankenfest into one func
	// 		packageList := processPkgsArgs(args)
	// 		repo := ParseFrankenfest(workingDir)
	// 		runTest(packageList, &repo, &workingDir)
	// 	},
	// }
)

// func runTest(packageList []string, repo *Frankenfest, workingDir *string) {
// 	// TODO: check deps
//
// 	if len(packageList) == 1 && packageList[0] == "all" {
// 		// TODO: dependency check ++ build those dependencies!
// 		for _, r := range repo.pkgs {
// 			if len(r.internalDeps) > 0 {
// 				for _, dep := range r.internalDeps {
// 					// TODO: this needs to run recursively
// 					dep_idx := 0
// 					for i, r2 := range repo.pkgs {
// 						if r2.name == dep {
// 							dep_idx = i
// 						}
// 					}
// 					cmd_string := strings.Split(repo.pkgs[dep_idx].build[0], " ")
// 					cmd := exec.Command(cmd_string[0], cmd_string[1:]...)
// 					cmd.Dir = filepath.Join(*workingDir, repo.pkgs[dep_idx].path)
// 					cmd.Stdout = os.Stdout
// 					cmd.Stderr = os.Stderr
// 					if err := cmd.Run(); err != nil {
// 						log.Fatal(CommandError+"::\n", err)
// 					}
// 				}
// 			}
// 			for _, t := range r.test {
// 				cmd_string := strings.Split(t, " ")
// 				cmd := exec.Command(cmd_string[0], cmd_string[1:]...)
// 				cmd.Dir = filepath.Join(*workingDir, r.path)
// 				cmd.Stdout = os.Stdout
// 				// capture this in a variable instead, so we can log it
// 				// properly
// 				cmd.Stderr = os.Stderr
// 				if err := cmd.Run(); err != nil {
// 					log.Fatal(CommandError+"::\n", err)
// 				}
// 			}
// 		}
// 	}
//
// 	// fmt.Println(packageList)
// 	// fmt.Println(repo)
// }

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
