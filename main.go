package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"errors"
	"os/exec"
	"os/user"
	"strings"
)

var (
	workDir     = flag.String("f", ".", "Working directory to be used during runtime")
	program     = flag.String("b", "", "Program to start in host OS binary format for instance python")
	programArgs = flag.String("a", "", "space separated arguement to be pass to the program for instance python filename.py")
	stdout = flag.String("o","","file path to show stdout and stderror result. default to stdout.in home directory")
)


func main() {
	flag.Parse()

	if err := os.Chdir(*workDir); err != nil {
		log.Fatal(err)
	}

	if *program == "" {
		log.Fatal("You must provide a program to run in this mode")
	}
	if *programArgs == "" {
		fmt.Fprintln(os.Stderr, "Runing program without arguments")
	}
	
	var (
		std *os.File 
		err error
	)

	if *stdout==""{
		outfilename:="daemonizeout.in"
		stdout=&outfilename
	}

	std,err=os.OpenFile(getHomeDir()+"/"+*stdout,os.O_RDWR|os.O_APPEND|os.O_CREATE,0777)
	if err!=nil{
		fmt.Fprintln(os.Stderr, "Not able to create output file %T",err)
	}

	programArgs := strings.Map(reWithSpace, *programArgs)
	argsList := strings.Split(programArgs, " ")

	//strip out space
	for i := 0; i < len(argsList); i++ {
		argsList[i] = strings.TrimSpace(argsList[i])
	}

	cmd := exec.Command(*program, argsList...)
	
	if std!=nil{
		cmd.Stderr=std
		cmd.Stdout=std
	}
	
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cmd.Process.Pid)
	//orphan child to make it a deamon(unix style)
	os.Exit(0)

}

//reWithSpace subtitutes all white space characters to space
func reWithSpace(r rune) rune {
	switch {
	case r == '\t':
		return ' '
	case r == '\n':
		return ' '
	case r == '\r':
		return ' '
	}

	return r
}

func getHomeDir()string{
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}

	if len(usr.HomeDir)<=0{
		log.Fatal(errors.New("current user have no home directory"))
	}

	return usr.HomeDir

}
