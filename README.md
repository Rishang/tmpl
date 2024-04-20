# tmpl
A CLI tool for replacing content in file or directory based on Jinja2 templates.

## Table of Contents
- [tmpl](#tmpl)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Usage](#usage)
    - [How to use the CLI tool](#how-to-use-the-cli-tool)

## Installation

Download the latest release from the [releases page](https://github.com/Rishang/tmpl/releases) and extract the binary to a directory in your Linux PATH.


## Usage

### How to use the CLI tool

Here is the help message for the CLI tool:

```shell
tmpl --help

tmpl is a CLI tool for replacing content in files based on Jinja2 templates.

Options:
  -c string
        Specify the path to the JSON configuration file. Default is 'tmpl.json'. (default "tmpl.json")
  -f string
        Specify the path to a single file to update.
  -p string
        Specify the path to the directory where files will be updated.

Examples:
  tmpl -p /path/to/directory -c /path/to/config.json
  tmpl -f /path/to/file.txt -c /path/to/config.json

The config file should contain JSON formatted like below example:
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
  File.txt:

  Hello {{ appName }}! Your secret is {{ secret }} at {{ cwd }}.
  ----
  Output:
  Hello Gjinja! Your secret is mysecret as /home/user.
```
