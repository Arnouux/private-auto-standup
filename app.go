package main

import (
    "os"
    "fmt"
    "os/exec"
    "bufio"
    "strings"
    "time"
    "strconv"
)

func main() {
    app := "git"

    // git log --after="date today-1"
    untilYesterday := time.Now().Add(time.Duration(-1) * time.Hour * 24)
    arg0 := "log"
    arg1 := "--after=" + fmt.Sprintf("%d-%d-%d", untilYesterday.Year(), untilYesterday.Month(), untilYesterday.Day())
    cmd := exec.Command(app, arg0, arg1)
    stdout, err := cmd.Output()
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    commits := make([]string, 0)
    scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
    for scanner.Scan() {
        line := scanner.Text()
        if(strings.HasPrefix(line, "commit")) {
            commits = append(commits, line[7:13])
        }
    }

    dataSaved := ""

    arg0 = "show"
    for i, commit := range commits {
        arg1 = commit
        cmd = exec.Command(app, arg0)
        stdout, err = cmd.Output()
        if err != nil {
            fmt.Println(err.Error())
            return
        }
    
        scanner = bufio.NewScanner(strings.NewReader(string(stdout)))
        savingData := false
        for scanner.Scan() {
            line := scanner.Text()
            if(strings.HasPrefix(line, "Date")) {
                scanner.Scan()
                scanner.Scan()
                line = scanner.Text()
                dataSaved += "Objective of commit " + strconv.Itoa(i) + ": " + line + "\n"
                savingData = true
            }
    
            if (savingData && len(line) > 1 && (line[0] == '+' || line[0] == '-')) {
                dataSaved += line + "\n"
            }
        }
    
    }

    dataSaved += "\nAbove this are the commits of yesterday\nBelow is the work in progress:\n"
    
    // git diff .
    arg0 = "diff"
    arg1 = "."
    cmd = exec.Command(app, arg0, arg1)
    stdout, err = cmd.Output()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    scanner = bufio.NewScanner(strings.NewReader(string(stdout)))
    for scanner.Scan() {
        line := scanner.Text()

        if (len(line) > 1 && (strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-"))) {
            dataSaved += line + "\n"
        }
    }

    dataSaved += "Explain how the commits of yesterday have changed the code and then explain the current work in progress (refer to commits and current work separately only). Only focus on lines starting with '+' (for added code) and '-' (for removed code)"

    fmt.Println(dataSaved)

    key, err := os.ReadFile("api.key")
    if err != nil {
        errorShown := fmt.Sprintf("Could not retrieve API key from file: %s", err)
        panic(errorShown)
    }

    // todo send to gpt api
    fmt.Println(string(key))
}