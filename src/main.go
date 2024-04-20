package main

import (
    "encoding/base64"
    "encoding/json"
    "flag"
    "fmt"
    "github.com/noirbizarre/gonja"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
)

type ConfigItem struct {
    Type  string `json:"type"`
    Key   string `json:"key"`
    Value string `json:"value"`
}

func usage() {
    fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
    fmt.Println("tmpl is a CLI tool for replacing content in files based on Jinja2 templates.")
    fmt.Println("\nOptions:")
    flag.PrintDefaults()
    fmt.Println("\nExamples:")
    fmt.Println("  tmpl -p /path/to/directory -c /path/to/config.json")
    fmt.Println("  tmpl -f /path/to/file.txt -c /path/to/config.json")
    fmt.Println("\nThe config file should contain JSON formatted like below example:")
    fmt.Println(`
    [
      {
        "type": "string",
        "key": "appName",
        "value": "Gjinja"
      },
      {
        "type": "base64",
        "key": "secret",
        "value": "bXlzZWNyZXQ="
      },
      {
        "type": "command",
        "key": "cwd",
        "value": "pwd"
      }
    ]
    ----
    Example:
    tmpl -f /path/to/file.txt -c /path/to/config.json
    
    file.txt:
  
    Hello {{ appName }}! Your secret is {{ secret }} at {{ cwd }}.
    ----
    Output saved to file.txt:
    Hello Gjinja! Your secret is mysecret as /home/user.
  `)
}

func readConfig(path string) ([]ConfigItem, error) {
    var items []ConfigItem
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(data, &items)
    return items, err
}

func processConfigItem(item ConfigItem) (interface{}, error) {
    switch item.Type {
    case "base64":
        decodedBytes, err := base64.StdEncoding.DecodeString(item.Value)
        if err != nil {
            return nil, err
        }
        return string(decodedBytes), nil
    case "command":
        return executeCommand(item.Value)
    default:
        return item.Value, nil
    }
}

func executeCommand(cmd string) (string, error) {
    out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("command execution failed: %v", err)
    }
    return string(out), nil
}

func updateFile(filePath string, context map[string]interface{}) error {
    content, err := ioutil.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("failed to read file: %v", err)
    }

    tmpl, err := gonja.FromString(string(content))
    if err != nil {
        return fmt.Errorf("failed to parse template: %v", err)
    }

    rendered, err := tmpl.Execute(gonja.Context(context))
    if err != nil {
        return fmt.Errorf("failed to execute template: %v", err)
    }

    return ioutil.WriteFile(filePath, []byte(rendered), 0644)
}


func main() {
    path := flag.String("p", "", "Specify the path to the directory where files will be updated.")
    configFile := flag.String("c", "tmpl.json", "Specify the path to the JSON configuration file. Default is 'tmpl.json'.")
    file := flag.String("f", "", "Specify the path to a single file to update.")
    flag.Usage = usage
    flag.Parse()

    if *file != "" && *path != "" {
        fmt.Println("Error: -f and -p flags are mutually exclusive.")
        return
    }

    configItems, err := readConfig(*configFile)
    if err != nil {
        fmt.Println("Error reading config, for cli reference use --help flag:", err)
        return
    }

    context := make(map[string]interface{})
    for _, item := range configItems {
        value, err := processConfigItem(item)
        if err != nil {
            fmt.Printf("Error processing config item '%s': %v\n", item.Key, err)
            continue
        }
        context[item.Key] = value
    }

    if *file != "" {
        fmt.Println("Updating:", *file)
        if err := updateFile(*file, context); err != nil {
            fmt.Println("Error updating file:", err)
        }
    } else if *path != "" {
        err = filepath.Walk(*path, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            if !info.IsDir() {
                fmt.Println("Updating:", info.Name())
                err := updateFile(path, context)
                if err != nil {
                    fmt.Println("Error updating file:", err)
                }
            }
            return nil
        })

        if err != nil {
            fmt.Println("Error processing files:", err)
        }
    } else {
        fmt.Println("Error: Either -f or -p must be specified.")
        flag.Usage()
    }
}
