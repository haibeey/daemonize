package main

import (
	
	"os"
	"log"
	"fmt"
	"flag"
	"os/exec"
	"strings"
)

var (
	workDir   = flag.String("f", ".", "Working directory to be used during runtime")
	createService= flag.Bool("s",false,"set to True if you want to run the app as service else false")
	program = flag.String("b","","Program to start in host OS binary format for instance python")
	programArgs = flag.String("a","","space separated arguement to be pass to the program for instance python filename.py")
)

func main(){
	flag.Parse();
	
	if err:=os.Chdir(*workDir); err!=nil{
		log.Fatal(err)
	}
	
	if *createService{
		//TODO(WIP)
	}else{
		if *program==""{
			log.Fatal("You must provide a program to run in this mode")
		}
		if *programArgs==""{
			fmt.Fprintln(os.Stderr,"Runing program without arguments")
		}

		programArgs:=strings.Map(reWithSpace,*programArgs)
		argsList:=strings.Split(programArgs," ")

		args:=[]string{}
		args=append(args,argsList...)

		cmd := exec.Command(*program,args...)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(cmd.Process.Pid)
		//orphan child to make it a deamon(unix style)
		os.Exit(0)
	}
}

//reWithSpace subtitutes all white space characters to space
func reWithSpace(r rune) rune {
    switch {
    case r=='\t':
        return ' '
    case r =='\n':
		return ' '
    }
    return r
}