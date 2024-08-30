package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	program     = flag.String("program", "", "Program to start in host OS binary format for instance python")
	programArgs = flag.String("args", "", "Space separated argument to be pass to the program for instance python filename.py")
	name        = flag.String("name", "", "Name of the program")
	show        = flag.Bool("show", false, "Show the list of all processes")
	kill        = flag.Bool("kill", false, "Kill process by name or pid")
)

type Config struct {
	Program     string
	Args        []string
	Name        string
	Show        bool
	Kill        bool
	stdFileName string
	executor    Executor
}

type Executor interface {
	Command(program, stdFileName string, args ...string) int
}

func NewConfig(program, name string, args []string) *Config {
	return &Config{
		Program:  program,
		Name:     name,
		Args:     args,
		executor: SysExecutor{},
	}
}

func (c *Config) Run() {
	if c.Show {
		showEntry()
		os.Exit(0)
	}

	if c.Kill {
		if len(c.Name) <= 0 {
			log.Fatal("name or pid must be passed")
		}

		err := killEntry(c.Name)
		if err != nil {
			log.Fatalf("Error killing program %s", err.Error())
		}

		os.Exit(0)
	}

	if c.Program == "" {
		log.Fatal("You must provide a program to run in this mode")
	}
	if len(c.Args) <= 0 {
		fmt.Fprintln(os.Stderr, "Runing program without arguments")
	}

	if c.Name == "" {
		c.Name = filepath.Base(c.Program)
	}

	pid := c.executor.Command(c.Program, c.stdFileName, c.Args...)
	err := addEntry(pid, c.Name)
	if err != nil {
		log.Fatalf("Entry error %s", err.Error())
		os.Exit(0)
	}
	os.Exit(0)
}

type SysExecutor struct{}

func (s SysExecutor) Command(program, stdFileName string, args ...string) int {

	cmd := exec.Command(program, args...)
	std, err := os.OpenFile(stdFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalf("Not able to create output file %s", err.Error())
	}

	if std != nil {
		cmd.Stderr = std
		cmd.Stdout = std
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Command woundn't start due : %s", err.Error())
	}

	defer std.Close()

	return cmd.Process.Pid

}

type Entry struct {
	Pid  int
	Name string
}

type EntryList []Entry

func main() {
	flag.Parse()

	if err := createDaemonizeDir(); err != nil {
		log.Fatalf("Couldn't create dir due to %s", err.Error())
	}

	programArgs := strings.Map(reWithSpace, *programArgs)
	argsList := strings.Split(programArgs, " ")

	for i := 0; i < len(argsList); i++ {
		argsList[i] = strings.TrimSpace(argsList[i])
	}

	config := NewConfig(*program, *name, argsList)
	config.Show = *show
	config.Kill = *kill

	config.stdFileName = fmt.Sprintf("%s/.daemonize/%s.in", getHomeDir(), config.Name)

	config.Run()

}
