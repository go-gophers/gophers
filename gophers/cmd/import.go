package cmd

import (
	"go/build"
	"go/types"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/go/loader"

	"github.com/go-gophers/gophers/utils/log"
)

type importData struct {
	PackageName       string
	PackageImportPath string
	PackageDir        string
	Tests             []string

	Test         bool
	Load         bool
	LoadWeighted bool
}

var mainTemplate = template.Must(template.New("main").Parse(`// +build ignore

package main

// generated with https://github.com/go-gophers/gophers

import (
    "flag"
    "os"

    "github.com/go-gophers/gophers/gophers/runner"
    "github.com/go-gophers/gophers/utils/log"

    "{{ .PackageImportPath }}"
)

func main() {
    flag.Parse()

    r := runner.New(log.New(os.Stderr, "", 0), "127.0.0.1:10311")
{{- range .Tests }}
    r.Add("{{ . }}", {{ $.PackageName }}.{{ . }}, 1)
{{- end }}

{{ if .Test }}
    r.Test(nil)
{{ else }}
    l, err := runner.NewStepLoader(5, 10, 1, 1 * time.Second)
    if err != nil {
        panic(err)
    }

	{{ if .Load }}
		r.Load(nil, l)
	{{ else }}
		r.LoadWeighted(l)
	{{ end -}}
{{ end -}}
}
`))

func installPackage() {
	// check go env
	if debugF {
		cmd := exec.Command(GoBin, "env")
		b, err := cmd.CombinedOutput()
		log.Printf(strings.Join(cmd.Args, " "))
		log.Printf("\n%s", b)
		if err != nil {
			log.Fatal(err)
		}
	}

	// install package
	args := []string{"install", "-v"}
	cmd := exec.Command(GoBin, args...)
	log.Printf("Running %s", strings.Join(cmd.Args, " "))
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%s", b)
		log.Fatal(err)
	}
}

func extractTestFunctions(scope *types.Scope) []string {
	var tests []string
	for _, name := range scope.Names() {
		if !strings.HasPrefix(name, "Test") {
			continue
		}

		if f, ok := scope.Lookup(name).(*types.Func); ok {
			sig := f.Type().(*types.Signature)

			// basic signature checks
			if sig.Recv() != nil {
				log.Printf("Skipping %q - test function should not be a method", f.String())
				continue
			}
			if sig.Variadic() {
				log.Printf("Skipping %q - test function should not be variadic", f.String())
				continue
			}
			if sig.Results() != nil {
				log.Printf("Skipping %q - test function should not return result", f.String())
				continue
			}

			// check params
			params := sig.Params()
			if params != nil || params.Len() == 1 {
				if named, ok := params.At(0).Type().(*types.Named); ok {
					if named.Obj().Name() == "TestingT" {
						tests = append(tests, f.Name())
						continue
					}
				}
			}

			log.Printf("Skipping %q - test function should have one parameter of type gophers.TestingT", f.String())
		}
	}

	return tests
}

func importPackage(dir string) *importData {
	installPackage()

	buildPkg, err := build.ImportDir(dir, 0)
	if err != nil {
		log.Fatalf(`build.ImportDir(%q, 0): %s`, dir, err)
	}

	// load package
	var conf loader.Config
	conf.Import(buildPkg.ImportPath)
	prog, err := conf.Load()
	if err != nil {
		log.Fatalf("conf.Load(): %s", err)
	}

	// get our single package
	packages := make([]string, 0, len(prog.Imported))
	for p := range prog.Imported {
		packages = append(packages, p)
	}
	if len(packages) != 1 {
		log.Fatalf("expected 1 package, got %d: %v", len(packages), packages)
	}
	pack := prog.Imported[packages[0]].Pkg

	// TODO compare pack and buildPkg

	tests := extractTestFunctions(pack.Scope())

	return &importData{
		PackageName:       pack.Name(),
		PackageImportPath: buildPkg.ImportPath,
		PackageDir:        buildPkg.Dir,
		Tests:             tests,
	}
}

func renderTemplate(data *importData, filename string) {
	path := filepath.Join(data.PackageDir, filename)
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = mainTemplate.Execute(f, data)
	if err != nil {
		log.Fatal(err)
	}

	args := []string{"-s", "-w", path}
	cmd := exec.Command(GoFmtBin, args...)
	log.Printf("Running %s", strings.Join(cmd.Args, " "))
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%s", b)
		log.Fatal(err)
	}

	log.Printf("%s created", path)
}
