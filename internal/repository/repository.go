/*
Copyright © 2023 Tommy Breslein <github.com/tbreslein>
*/
package repository

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type ParseError string

const (
	DecodeError       ParseError = "DECODE_ERROR"
	MissingFieldError ParseError = "MISSING_FIELD_ERROR"
	UnknownError      ParseError = "UNKNOWN_PARSE_ERROR"
)

type CommandType int

const (
	Build CommandType = iota
	Test
	Format
	Lint
	Run
)

type Command struct {
	T         CommandType
	CmdString string
}

const (
	ProjectManifestName = "frankenfest.toml"
	ConfigFile          = "frankenrepo.toml"
)

type Frankenfest struct {
	Version int
	pkgs    []FrankenPkg
}

type FrankenPkg struct {
	name         string
	path         string
	externalDeps []string
	internalDeps []string
	build        []string
	test         []string
	// format       []string
	// lint         []string
	// run map[CommandType][]command
}

type Repo struct {
	ExternalDeps []string
	Pkgs         PkgGraph
}

type PkgGraph struct{}

type PkgGraphNode struct {
	Name       string
	Path       string
	BuildCmds  []exec.Cmd
	TestCmds   []exec.Cmd
	FormatCmds []exec.Cmd
	LintCmds   []exec.Cmd
	DependsOn  []*PkgGraphNode
	DependedBy []*PkgGraphNode
}

func InitRepo(frankenfestDir string, args []string) Frankenfest {
	frankenfest := initFrankenfest(frankenfestDir)
	// targets := processPkgsArgs(args)
	return frankenfest
}

func initFrankenfest(frankenfestDir string) Frankenfest {
	frankenfestPath := filepath.Join(frankenfestDir, ProjectManifestName)
	bytes, err := os.ReadFile(frankenfestPath)
	if err != nil {
		log.Fatalln(err)
	}

	frankenfest := Frankenfest{}
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
				DecodeError+":: ", frankenfestPath, ":", row, ":", col, "\n",
				decode_err.String(),
			)
		} else if errors.As(err, &smissing_err) {
			log.Fatal(
				MissingFieldError+":: ", frankenfestPath, "\n",
				smissing_err.String(),
			)
		} else {
			log.Fatal(
				UnknownError+":: ", frankenfestPath, "\n",
				err,
			)
		}
	}

	return frankenfest
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
