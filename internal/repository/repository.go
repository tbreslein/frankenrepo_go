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
	"slices"
	"sort"
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

type targetPrefix string

const (
	buildNamePrefix targetPrefix = "__BUILD__"
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
	buildTargets []target // slice into targetList of all build targets
}

// type targetGraph struct{}

type target struct {
	Name      string
	path      string
	cmds      []exec.Cmd
	t         targetType
	dependsOn []*target
	depNames  []string
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
		if len(pkg.build) > 0 {
			repo.targetList = append(
				repo.targetList,
				target{
					Name:     string(buildNamePrefix) + pkg.name,
					path:     pkg.path,
					t:        build,
					depNames: pkg.internalDeps,
				},
			)
		}
	}
	// fill in dependsOn and dependedBy.
	// add the 'all' targets, which are the special targets of each targetType
	// other than 'run' that depend on every other target of their group.

    slices.SortFunc[](repo.targetList)
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
