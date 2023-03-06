package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed init-templates/schema.graphql
var schemaFileContent string

//go:embed init-templates/gqlgen.yml.gotmpl
var configFileTemplate string

var modregex = regexp.MustCompile(`module ([^\s]*)`)

func getConfigFileContent(pkgName string) string {
	var buf bytes.Buffer
	if err := template.Must(template.New("gqlgen.yml").Parse(configFileTemplate)).Execute(&buf, pkgName); err != nil {
		panic(err)
	}
	return buf.String()
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, fs.ErrNotExist)
}

func findModuleRoot(dir string) (roots string) {
	if dir == "" {
		panic("dir not set")
	}
	dir = filepath.Clean(dir)

	// Look for enclosing go.mod.
	for {
		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
			return dir
		}
		d := filepath.Dir(dir)
		if d == dir { // the parent of the root is itself, so we can go no further
			break
		}
		dir = d
	}
	return ""
}

func initFile(filename, contents string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return fmt.Errorf("unable to create directory for file '%s': %w\n", filename, err)
	}
	if err := os.WriteFile(filename, []byte(contents), 0o644); err != nil {
		return fmt.Errorf("unable to write file '%s': %w\n", filename, err)
	}

	return nil
}

func extractModuleName(modDir string) string {
	content, err := os.ReadFile(filepath.Join(modDir, "go.mod"))
	if err != nil {
		log.Fatal("couldn't read file")
	}

	for {
		advance, tkn, err := bufio.ScanLines(content, false)
		if err != nil {
			panic(fmt.Errorf("error parsing mod file: %w", err))
		}
		if advance == 0 {
			break
		}
		s := strings.Trim(string(tkn), " \t")
		if len(s) != 0 && !strings.HasPrefix(s, "//") {
			break
		}
		if advance <= len(content) {
			content = content[advance:]
		}
	}
	moduleName := string(modregex.FindSubmatch(content)[1])
	return moduleName
}

func main() {

	configFlag := flag.String("config", "gqlgen.yml", "the config filename")
	schemaFlag := flag.String("schema", "graph/schema.graphql", "where to write the schema stub to")
	//serverFlag := flag.String("server", "where to write the server stub to", "server.go")
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		err = fmt.Errorf("unable to determine current directory:%w", err)
		log.Println(err)
	}

	modRoot := findModuleRoot(cwd)
	if modRoot == "" {
		err = fmt.Errorf("go.mod is missing. Please, do 'go mod init' first\n")
		log.Println(err)
	}

	// create config
	fmt.Println("Creating", *configFlag)
	moduleName := extractModuleName(modRoot)
	if err := initFile(*configFlag, getConfigFileContent(moduleName)); err != nil {
		log.Println(err)
	}

	// create schema
	fmt.Println("Creating", *schemaFlag)
	if err := initFile(*schemaFlag, schemaFileContent); err != nil {
		log.Println(err)
	}
}
