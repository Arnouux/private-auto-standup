package main

import (
    "fmt"
    "os/exec"
    "bufio"
    "strings"
)

func main() {
    app := "git"

    arg0 := "show"

    cmd := exec.Command(app, arg0)
    stdout, err := cmd.Output()

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // Print the output
    scanner := bufio.NewScanner(strings.NewReader(string(stdout)))
    dataSaved := ""
    savingData := false
    for scanner.Scan() {
        line := scanner.Text()
        if(strings.HasPrefix(line, "Date")) {
            scanner.Scan()
            scanner.Scan()
            line = scanner.Text()
            dataSaved += "Objective of this commit" + line + "\n"
            savingData = true
        }

        if (savingData && len(line) > 1 && (line[0] == '+' || line[0] == '-')) {
            dataSaved += line + "\n"
        }
    }

    fmt.Println(dataSaved)
    
    // git diff .
}