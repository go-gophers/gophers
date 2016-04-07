package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/mattn/go-shellwords"
)

var (
	systemF = flag.String("system", "", "Command to start system under test")
	testF   = flag.String("test", "", "Flags and arguments for test binaries")
	delayF  = flag.Duration("delay", time.Second, "Delay before running first test binary")
)

func run() int {
	testArgs, err := shellwords.Parse(*testF)
	if err != nil {
		log.Fatalf("Failed to parse %q: %s", *testF, err)
	}

	args, err := shellwords.Parse(*systemF)
	if err != nil || len(args) == 0 {
		log.Fatalf("Failed to parse %q: %s", *systemF, err)
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		log.Fatal(err)
	}
	args = args[1:]

	sutL := log.New(os.Stderr, path+": ", 0)
	sutL.Printf("Starting %s %v", path, args)
	sut := exec.Command(path, args...)
	sut.Stdout = os.Stderr
	sut.Stderr = os.Stderr
	err = sut.Start()
	if err != nil {
		sutL.Fatal(err)
	}
	sutL.Printf("PID %d", sut.Process.Pid)

	defer func() {
		err = sut.Process.Kill()
		if err != nil {
			sutL.Printf("Failed to stop: %s", err)
		} else {
			sutL.Printf("Stopped")
		}
	}()

	fis, err := ioutil.ReadDir(".")
	if err != nil {
		log.Print(err)
		return -1
	}

	time.Sleep(*delayF)

	for _, fi := range fis {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".test") {
			testL := log.New(os.Stderr, fi.Name()+": ", 0)
			testL.Printf("Running %s %v", "./"+fi.Name(), testArgs)
			test := exec.Command("./"+fi.Name(), testArgs...)
			test.Stdout = os.Stdout
			test.Stderr = os.Stdout
			err = test.Start()
			if err != nil {
				testL.Print(err)
				return -1
			}
			testL.Printf("PID %d", test.Process.Pid)

			err = test.Wait()
			if err != nil {
				testL.Print(err)
				if e, ok := err.(*exec.ExitError); ok {
					return e.Sys().(syscall.WaitStatus).ExitStatus()
				}
			}
		}
	}

	return 0
}

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		log.Printf("%s is a tool to start system under test and run tests for it.\n\n", os.Args[0])
		log.Printf("It starts system by running command specified by -system, then runs all *.test files")
		log.Printf("(generated with `go test -c`) in current directory with -test flags and arguments,")
		log.Printf("then stops system under test.")
		log.Printf("\nFlags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *systemF == "" {
		log.Printf("-system is required.\n\n")
		flag.Usage()
		os.Exit(2)
	}

	log.SetPrefix("gophers: ")
	os.Exit(run())
}
