//go:binary-only-package
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func execCommand(commandName string, params []string) bool {
	cmd := exec.Command(commandName, params...)

	fmt.Println(commandName, params)

	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()

	if err != nil {
		panic(err)
	}

	_ = cmd.Start()

	readerStd := bufio.NewReader(stdout)
	readerErr := bufio.NewReader(stderr)

	for {
		line, err2 := readerStd.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Print(line)
	}

	for {
		line, err2 := readerErr.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Print(line)
	}

	_ = cmd.Wait()
	return true
}

// Contains ...
func Contains(array []string, key string) (index int) {
	for index, val := range array {
		if strings.Contains(val, key) {
			return index
		}
	}
	return -1
}

func main() {
	var cmd []string
	cmd = append(cmd, os.Args[0])
	cmd = append(cmd, "go")
	cmd = append(cmd, "get")
	cmd = append(cmd, os.Args[1:]...)
	originCmd := make([]string, len(cmd))
	copy(originCmd, cmd)
	urlIndex := Contains(cmd, "golang.org/x/")
	if Contains(cmd, "get") != -1 && urlIndex != -1 {
		url := cmd[urlIndex]
		cmd[urlIndex] = strings.Replace(url, "golang.org/x/", "github.com/golang/", -1)
		execCommand(cmd[1], cmd[2:])
		moveToGoLangDir(cmd[urlIndex])
		execCommand(originCmd[1], originCmd[2:])
	} else {
		fmt.Println("不支持的命令，或错误的地址")
	}

}

func move(srcProjectPath, destProjectPath string) {
	err := os.MkdirAll(destProjectPath, os.ModePerm)
	check(err)
	err = os.RemoveAll(destProjectPath)
	check(err)
	err = os.Rename(srcProjectPath, destProjectPath)
	check(err)
}

func moveToGoLangDir(url string) {
	goPath, exist := os.LookupEnv("GOPATH")
	if exist {
		projectName := strings.Split(strings.Replace(url, "github.com/golang/", "", -1), "/")[0]
		var destProjectPathArr []string
		destProjectPathArr = append(destProjectPathArr, goPath)
		destProjectPathArr = append(destProjectPathArr, "src")
		destProjectPathArr = append(destProjectPathArr, "golang.org")
		destProjectPathArr = append(destProjectPathArr, "x")
		destProjectPathArr = append(destProjectPathArr, projectName)
		destProjectPath := path.Join(destProjectPathArr...)

		var srcProjectPathArr []string
		srcProjectPathArr = append(srcProjectPathArr, goPath)
		srcProjectPathArr = append(srcProjectPathArr, "src")
		srcProjectPathArr = append(srcProjectPathArr, "github.com")
		srcProjectPathArr = append(srcProjectPathArr, "golang")
		srcProjectPathArr = append(srcProjectPathArr, projectName)
		srcProjectPath := path.Join(srcProjectPathArr...)
		move(srcProjectPath, destProjectPath)
	}
}
