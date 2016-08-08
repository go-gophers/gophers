package cmd

import (
	"go/build"
	"go/types"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/go/loader"
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

var mainTemplate = template.Must(template.New("main").Parse(strings.TrimSpace(`
// +build ignore

package main

// generated with https://github.com/go-gophers/gophers

import (
    "flag"
    "os"
    "time"

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
`)))

func importPackage(path string, race bool) *importData {
	// check go env
	if debugF {
		cmd := exec.Command(GoBin, "env")
		b, err := cmd.CombinedOutput()
		log.Printf(strings.Join(cmd.Args, " "))
		log.Printf("%s", b)
		if err != nil {
			log.Fatal(err)
		}
	}

	// install package
	args := []string{"install", "-v"}
	if race {
		args = append(args, "-race")
	}
	args = append(args, path)
	cmd := exec.Command(GoBin, args...)
	log.Printf("Running %s", strings.Join(cmd.Args, " "))
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("%s", b)
		log.Fatal(err)
	}

	// import package
	buildPack, err := build.Import(path, WD, 0)
	if err != nil {
		log.Fatalf("build.Import(%q, %q, 0): %s", path, WD, err)
	}
	var conf loader.Config
	conf.Import(path)
	prog, err := conf.Load()
	if err != nil {
		log.Fatalf("conf.Load(): %s", err)
	}
	pack := prog.Imported[path].Pkg
	if buildPack.Name != pack.Name() || buildPack.ImportPath != pack.Path() {
		log.Fatalf("failed to locate package:\n%#v\n%#v", buildPack, pack)
	}

	// get test functions
	var tests []string
	for _, name := range pack.Scope().Names() {
		if !strings.HasPrefix(name, "Test") {
			continue
		}

		if f, ok := pack.Scope().Lookup(name).(*types.Func); ok {
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

	return &importData{
		PackageName:       pack.Name(),
		PackageImportPath: path,
		PackageDir:        buildPack.Dir,
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
