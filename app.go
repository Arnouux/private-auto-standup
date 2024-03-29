package main

import (
    "os"
    "fmt"
    "bytes"
    "net/http"
    "encoding/json"
    "os/exec"
    "bufio"
    "strings"
    "time"
    "strconv"
    "io"
)

type Message struct {
    Role string         `json:"role"`
    Content string      `json:"content"`
}

type PostData struct {
    Model string        `json:"model"`
    Messages []Message  `json:"messages"`
}

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
        errorMsg := fmt.Sprintf("Could not retrieve API key from file: %s", err)
        panic(errorMsg)
    }
    keyString := "Bearer " + string(key)

    // what if > 3000 chars
    fmt.Println(len(dataSaved))

    message := Message {
        Role: "user",
        Content: dataSaved,
    }

    postData := PostData {
        Model: "gpt-3.5-turbo",
        Messages: []Message{message},
    }

    body, _ := json.Marshal(postData)

    request, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Authorization", keyString)

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        errorMsg := fmt.Sprintf("Could not send request to openai: %s", err)
        panic(errorMsg)
    }
    defer response.Body.Close()

    result, err := io.ReadAll(response.Body)
    if err != nil {
        errorMsg := fmt.Sprintf("Could not read response from openai: %s", err)
        panic(errorMsg)
    }

    fmt.Println(string(result))

}