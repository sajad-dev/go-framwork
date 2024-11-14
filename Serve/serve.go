package serve

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

var mu sync.Mutex

var addr string

func Serve() {

	getAddress()

	color.Blue("Adress Local : http://127.0.0.1:8000")

	Run()
	runTest()

	Compile()

}

func runTest() {
	command := exec.Command("go", "test", "-v", addr+"/"+"Test")
	var out, stderr strings.Builder
	command.Stdout = &out
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		color.Red(stderr.String())
	}
	if strings.Contains(out.String(), "FAIL") {
		color.Yellow(out.String())
	} else {
		color.Cyan(out.String())
	}
}

func getAddress() {
	ou := exec.Command("pwd")
	adderss, _ := ou.Output()
	addr = strings.TrimSpace(string(adderss))
}
func Compile() {

	for {
		if len(os.Args) > 1 {
			os.Exit(1)
		}
		Commands()
		time.Sleep(1 * time.Second)
	}
}

func build() {
	comand := exec.Command("go", "build", addr)
	var out, stderr strings.Builder
	comand.Stdout = &out
	comand.Stderr = &stderr

	err := comand.Run()

	if err != nil {
		color.Red("Error In Compile: " + stderr.String())
	}

}

func Commands() {
	comand := exec.Command("git", "status", addr)
	ou, err := comand.Output()
	if err != nil {
		return
	}

	output := strings.Split(string(ou), "Changes not staged for commit:")
	if len(output) > 1 && strings.Contains(output[1], "modified") {

		build()

		killProcess()
		comandGit := exec.Command("git", "add", ".", addr)
		if err := comandGit.Run(); err != nil {
		}

		go Run()

	}
}

func Run() {
	currentTime := time.Now()

	fmt.Print("Compile in : ")
	color.Red(currentTime.Format("2006/01/02 15:04:05"))
	killProcess()
	comandRun := exec.Command("")
	if len(os.Args) > 1 {
		comandRun = exec.Command(addr+"/"+"go-framwork", os.Args[1:]...)

	} else {
		comandRun = exec.Command(addr + "/" + "go-framwork")

	}
	stdout, err := comandRun.StdoutPipe()
	if err != nil {
		color.Red(err.Error())
	}
	stderr, err := comandRun.StderrPipe()
	if err != nil {
		color.Red(err.Error())
	}

	err = comandRun.Start()
	if err != nil {
		color.Blue("Error: " + err.Error())

	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			color.Green(scanner.Text())
		}
		scanner = bufio.NewScanner(stderr)
		for scanner.Scan() {
			color.Red(scanner.Text())
		}
	}()
}

func killProcess() {
	commandPr := exec.Command("lsof", "-t", "-i", ":8000")
	ou, err := commandPr.Output()
	if err != nil {
		return
	}
	pid := strings.TrimSpace(string(ou))
	pid = strings.Split(pid, "\n")[0]
	if pid != "" {
		commandKill := exec.Command("kill", pid)
		var out, stderr strings.Builder
		commandKill.Stdout = &out
		commandKill.Stderr = &stderr
		err := commandKill.Run()
		if err != nil {
			color.Blue(stderr.String() + " " + pid + "=> 2000")
		}
	}

	commandPr = exec.Command("lsof", "-t", "-i", ":3000")
	ou, err = commandPr.Output()
	if err != nil {
		return
	}
	pid = strings.TrimSpace(string(ou))
	pid = strings.Split(pid, "\n")[0]
	if pid != "" {
		commandKill := exec.Command("kill", pid)
		var out, stderr strings.Builder
		commandKill.Stdout = &out
		commandKill.Stderr = &stderr
		err := commandKill.Run()
		if err != nil {
			color.Blue(stderr.String() + " " + pid + "=> 3000")
		}
	}
}
