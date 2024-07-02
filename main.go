package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	program     = flag.String("program", "", "Program to start in host OS binary format for instance python")
	programArgs = flag.String("args", "", "Space separated arguement to be pass to the program for instance python filename.py")
	name        = flag.String("name", "", "Name of the program")
	show        = flag.Bool("show", false, "Show list of process")
	kill        = flag.Bool("kill", false, "Kill process by name")
)

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

	if *show {
		showEntry()
		os.Exit(0)
	}

	if *kill {
		if len(*name) <= 0 {
			fmt.Println(len(*name), *name)
			log.Fatal("name must be passed")
		}

		err := killEntry(*name)
		if err != nil {
			log.Fatalf("Error killing program %s", err.Error())
		}

		os.Exit(0)
	}

	if *program == "" {
		log.Fatal("You must provide a program to run in this mode")
	}
	if *programArgs == "" {
		fmt.Fprintln(os.Stderr, "Runing program without arguments")
	}

	if *name == "" {
		*name = filepath.Base(*program)
	}

	var (
		std *os.File
		err error
	)

	outFilename := fmt.Sprintf("%s/.daemonize/%s.in", getHomeDir(), *name)
	std, err = os.OpenFile(outFilename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Not able to create output file %s", err.Error())
	}

	defer std.Close()

	programArgs := strings.Map(reWithSpace, *programArgs)
	argsList := strings.Split(programArgs, " ")

	for i := 0; i < len(argsList); i++ {
		argsList[i] = strings.TrimSpace(argsList[i])
	}

	argsList = append(argsList, []string{">", outFilename}...)

	cmd := exec.Command(*program, argsList...)

	if std != nil {
		cmd.Stderr = std
		cmd.Stdout = std
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Command woundn't start due : %s", err.Error())
	}

	log.Println(cmd.Process.Pid)
	err = addEntry(cmd.Process.Pid, *name)
	if err != nil {
		log.Fatalf("Entry error %s", err.Error())
		os.Exit(0)
	}
	os.Exit(0)

}

// reWithSpace subtitutes all white space characters to space
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

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if len(usr.HomeDir) <= 0 {
		log.Fatal(errors.New("current user have no home directory"))
	}

	return usr.HomeDir

}

func createDaemonizeDir() error {
	directoryname := fmt.Sprintf("%s/.daemonize", getHomeDir())
	if _, err := os.Stat(directoryname); os.IsNotExist(err) {
		err := os.Mkdir(directoryname, 0755) //create a directory and give it required permissions
		if err != nil {
			return err
		}

	} else {
		return nil
	}

	return nil

}

func getDaemonizeEntryList() (EntryList, error) {
	el := EntryList{}
	daemonizeFile := fmt.Sprintf("%s/.daemonize/daemonizeentry.in", getHomeDir())

	fileInfo, err := os.Stat(daemonizeFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			_, err := os.Create(daemonizeFile)
			return el, err
		} else {
			return el, err
		}
	}

	daemonize, err := os.OpenFile(daemonizeFile, os.O_RDWR, 0755)
	if err != nil {
		return el, err
	}

	defer daemonize.Close()

	if fileInfo.Size() <= 0 {
		return el, nil
	}

	var buf bytes.Buffer

	_, err = io.Copy(&buf, daemonize)
	if err != nil {
		return el, err
	}

	enc := gob.NewDecoder(&buf)
	if err := enc.Decode(&el); err != nil {
		return el, err
	}

	return el, err

}

func addEntry(pid int, programName string) error {

	e := Entry{pid, programName}

	entryList, err := getDaemonizeEntryList()
	if err != nil {
		return err
	}

	entryList = append(entryList, e)

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(entryList); err != nil {
		return err
	}

	daemonizeFile := fmt.Sprintf("%s/.daemonize/daemonizeentry.in", getHomeDir())
	daemonize, err := os.OpenFile(daemonizeFile, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer daemonize.Close()

	_, err = daemonize.Write(buf.Bytes())
	return err
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func showEntry() {

	entryList, err := getDaemonizeEntryList()
	if err != nil {
		log.Fatalf("Error showing entry %s", err.Error())
	}

	data := [][]string{}

	for index, entry := range entryList {
		data = append(data, []string{fmt.Sprintf("%d", entry.Pid), entry.Name})
		proc, err := process.NewProcess(int32(entry.Pid))
		if err != nil {
			log.Fatalf("Error getting process: %s", err.Error())
		}

		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			data[index] = append(data[index], "unknown")
		} else {
			data[index] = append(data[index], fmt.Sprintf("%.1f%s", cpuPercent, "%"))
		}

		memInfo, err := proc.MemoryInfo()
		if err != nil {
			data[index] = append(data[index], "unknown")
		} else {
			data[index] = append(data[index], formatBytes(memInfo.RSS))
		}

		ioCounters, err := proc.IOCounters()
		if err != nil {
			data[index] = append(data[index], "unknown")
		} else {
			data[index] = append(data[index], formatBytes(ioCounters.ReadBytes))
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"PID", "Name", "cpu", "memory", "disk"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()

}

func killEntry(name string) error {

	entryList, err := getDaemonizeEntryList()
	if err != nil {
		log.Fatalf("Error getting entry %s", err.Error())
	}

	el := EntryList{}
	for _, entry := range entryList {
		if entry.Name == name {
			proc, err := os.FindProcess(entry.Pid)
			if err != nil {
				log.Println(err)
			}
			proc.Kill()

		} else {
			el = append(el, entry)
		}
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(el); err != nil {
		return err
	}

	daemonizeFile := fmt.Sprintf("%s/.daemonize/daemonizeentry.in", getHomeDir())
	stat, err := os.Stat(daemonizeFile)
	if err != nil {
		return err
	}

	daemonize, err := os.OpenFile(daemonizeFile, os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	err = daemonize.Truncate(stat.Size())
	if err != nil {
		return err
	}

	daemonize.Close()

	daemonize, err = os.OpenFile(daemonizeFile, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer daemonize.Close()

	_, err = daemonize.Write(buf.Bytes())
	return err
}
