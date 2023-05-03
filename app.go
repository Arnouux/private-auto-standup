package main

import (
    "fmt"
    "os/exec"
    "bufio"
    "strings"
    "time"
)

func main() {
    app := "git"

    arg0 := "log"
    arg1 := "--date=format:'%Y-%m-%d'"
    cmd := exec.Command(app, arg0, arg1)
    stdout, err := cmd.Output()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
    commit, _, _ := "", "", ""
    for scanner.Scan() {
        line := scanner.Text()
        if(strings.HasPrefix(line, "commit")) {
            commit = line[7:13]
        }
        if(strings.HasPrefix(line, "Date")) {
            date, _ := time.Parse("yyyy-mm-dd", line[8:])
            fmt.Println(date)
            fmt.Println(line[8:])
        }
    }
    fmt.Println(commit)


    arg0 = "show"
    cmd = exec.Command(app, arg0)
    stdout, err = cmd.Output()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // Print the output
    scanner = bufio.NewScanner(strings.NewReader(string(stdout)))
    dataSaved := ""
    savingData := false
    for scanner.Scan() {
        line := scanner.Text()
        if(strings.HasPrefix(line, "Date")) {
            scanner.Scan()
            scanner.Scan()
            line = scanner.Text()
            dataSaved += "Objective of this commit:" + line + "\n"
            savingData = true
        }

        if (savingData && len(line) > 1 && (line[0] == '+' || line[0] == '-')) {
            dataSaved += line + "\n"
        }
    }

    dataSaved += "Above this are the commits of yesterday\nBelow is the work in progress:\n"
    
    // git diff .
    arg0 = "diff"
    arg1 = "."
    cmd = exec.Command(app, arg0, arg1)
    stdout, err = cmd.Output()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // Print the output
    scanner = bufio.NewScanner(strings.NewReader(string(stdout)))
    savingData = false
    for scanner.Scan() {
        line := scanner.Text()

        if (len(line) > 1 && (strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-"))) {
            dataSaved += line + "\n"
        }
    }

    fmt.Println(dataSaved)
}