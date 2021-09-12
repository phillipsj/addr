/*
Copyright © 2021 Jamie Phillips

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const adrTemplate = `
# {{.Number}}. {{.Title}}
======
Date: {{.Date}}
## Status
======
{{.Status}}
## Context
======
## Decision
======
## Consequences
======

`

type Adr struct {
	Number int
	Title  string
	Date   string
	Status AdrStatus
}

type AdrStatus string

const (
	PROPOSED   AdrStatus = "Proposed"
	ACCEPTED   AdrStatus = "Accepted"
	DEPRECATED AdrStatus = "Deprecated"
	SUPERSEDED AdrStatus = "Superseded"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new ADR file in the location specified in config file.",
	Long: `Creates a new ADR file in the location specified in config file.. For example:

addr new "my decision record"`,
	Run: func(cmd *cobra.Command, args []string) {
		a := Adr{
			Title:  strings.Join(args, " "),
			Date:   time.Now().Format("2006-02-01 15:04:05"),
			Status: PROPOSED,
		}
		adr, err := parseTemplate(a)
		if err != nil {
			panic(err)
		}
		file := strconv.Itoa(a.Number) + "-" + strings.Join(strings.Split(strings.Trim(a.Title, "\n \t"), " "), "-") + ".md"
		path := filepath.Join("", file)
		writeFile(path, adr)
		fmt.Println("new called")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func parseTemplate(config interface{}) (string, error) {
	out := new(bytes.Buffer)
	t := template.Must(template.New("compiled_template").Parse(adrTemplate))
	if err := t.Execute(out, config); err != nil {
		return "", err
	}
	return out.String(), nil
}

func writeFile(name string, content string) error {
	os.MkdirAll(filepath.Dir(name), 0755)
	err := ioutil.WriteFile(name, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}
