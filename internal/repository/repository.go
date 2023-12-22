/*
Copyright © 2023 Tommy Breslein <github.com/tbreslein>
*/
package repository

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/deckarep/golang-set/v2"
	"github.com/pelletier/go-toml/v2"
)

type ParseError string

const (
	decodeError       ParseError = "DECODE_ERROR"
	missingFieldError ParseError = "MISSING_FIELD_ERROR"
	unknownError      ParseError = "UNKNOWN_PARSE_ERROR"
)

type targetType int

const (
	build targetType = iota
	// test
	// format
	// lint
	// run
)

const (
	projectManifestName = "frankenfest.toml"
	configFile          = "frankenrepo.toml"
)

type frankenfest struct {
	Version int
	pkgs    []frankenPkg
}

type frankenPkg struct {
	name         string
	path         string
	externalDeps []string
	internalDeps []string
	build        []string
	// test         []string
	// format       []string
	// lint         []string
	// run map[CommandType][]command
}

type Repo struct {
	externalDeps mapset.Set[string]
	targetList   []target
	// targetGraph  TargetGraph
}

type targetGraph struct{}

type target struct {
	Name       string
	path       string
	cmds       []exec.Cmd
	t          targetType
	dependsOn  []*target
	dependedBy []*target
	// BuildCmds []exec.Cmd
	// TestCmds   []exec.Cmd
	// FormatCmds []exec.Cmd
	// LintCmds   []exec.Cmd
	// RunCmds   []exec.Cmd
}

func InitRepo(frankenfestDir string, args []string) Repo {
	repo := initFrankenfest(&frankenfestDir).toRepo()
	// targetStrings := processPkgsArgs(args)
	return repo
}

func (ff frankenfest) toRepo() Repo {
	repo := Repo{}
	for _, pkg := range ff.pkgs {
		for _, dep := range pkg.externalDeps {
			repo.externalDeps.Add(dep)
		}
	}
	return repo
}

func initFrankenfest(frankenfestDir *string) frankenfest {
	frankenfestPath := filepath.Join(*frankenfestDir, projectManifestName)
	bytes, err := os.ReadFile(frankenfestPath)
	if err != nil {
		log.Fatalln(err)
	}

	frankenfest := frankenfest{}
	err = toml.
		NewDecoder(strings.NewReader(string(bytes))).
		DisallowUnknownFields().
		Decode(&frankenfest)

	if err != nil {
		var decode_err *toml.DecodeError
		var smissing_err *toml.StrictMissingError

		if errors.As(err, &decode_err) {
			row, col := decode_err.Position()
			log.Fatal(
				decodeError+":: ", frankenfestPath, ":", row, ":", col, "\n",
				decode_err.String(),
			)
		} else if errors.As(err, &smissing_err) {
			log.Fatal(
				missingFieldError+":: ", frankenfestPath, "\n",
				smissing_err.String(),
			)
		} else {
			log.Fatal(
				unknownError+":: ", frankenfestPath, "\n",
				err,
			)
		}
	}

	return frankenfest
}

// func processPkgsArgs(args []string) []string {
// 	var packageList []string
// 	if len(args) == 0 {
// 		packageList = append(packageList, "all")
// 	} else {
// 		packageList = args
// 	}
// 	fmt.Print("Running frankenrepo on these packages: ")
// 	for _, a := range packageList {
// 		fmt.Print(a)
// 	}
// 	fmt.Println("")
// 	return packageList
// }
