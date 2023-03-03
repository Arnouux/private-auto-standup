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
        if(strings.HasPrefix(scanner.Text(), "Date")) {
            savingData = true
            fmt.Println(scanner.Text())
        }
        if (savingData) {
            dataSaved += scanner.Text()
        }
    }

    fmt.Println(dataSaved)
    
    // git diff .
}