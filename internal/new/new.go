/*******************************************************************************
 * Copyright 2017.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package new

import (
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/winc-link/hummingbird-cli/config"
	"os"
	"os/exec"
	"path/filepath"
)

type Project struct {
	ProjectName string `survey:"name"`
}

var CmdNew = &cobra.Command{
	Use:     "new",
	Example: "hb new demo-driver",
	Short:   "create a new driver layout.",
	Long:    `create a new driver layout.`,
	Run:     run,
}
var (
	repoURL     string
	ProjectName string
)

func init() {
	//CmdNew.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
	CmdNew.Flags().StringVarP(&ProjectName, "p", "p", ProjectName, "project name")

}
func NewProject() *Project {
	return &Project{}
}

func run(cmd *cobra.Command, args []string) {
	p := NewProject()
	if len(args) == 0 {
		err := survey.AskOne(&survey.Input{
			Message: "What is your project name?",
			Help:    "project name.",
			Suggest: nil,
		}, &p.ProjectName, survey.WithValidator(survey.Required))
		if err != nil {
			return
		}
	} else {
		p.ProjectName = args[0]
	}

	// clone repo
	yes, err := p.cloneTemplate()
	if err != nil || !yes {
		return
	}

	err = p.replacePackageName()
	if err != nil || !yes {
		return
	}

	err = p.modTidy()
	if err != nil || !yes {
		return
	}
	p.rmGit()
	fmt.Printf(config.LogoContent + "\n")
	fmt.Printf("🎉 Project \u001B[36m%s\u001B[0m created successfully!\n\n", p.ProjectName)
}

func (p *Project) cloneTemplate() (bool, error) {
	stat, _ := os.Stat(p.ProjectName)
	if stat != nil {
		var overwrite = false

		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Folder %s already exists, do you want to overwrite it?", p.ProjectName),
			Help:    "Remove old project and create new project.",
		}
		err := survey.AskOne(prompt, &overwrite)
		if err != nil {
			return false, err
		}
		if !overwrite {
			return false, nil
		}
		err = os.RemoveAll(p.ProjectName)
		if err != nil {
			fmt.Println("remove old project error: ", err)
			return false, err
		}
	}
	repo := ""

	if repoURL == "" {
		protocol := ""
		prompt := &survey.Select{
			Message: "Please select a protocol:",
			Options: []string{
				"MQTT",
				"TCP",
				"UDP",
				"CoAP",
				"HTTP",
				"WebSocket",
				"Modbus-TCP",
			},
			//Description: func(value string, index int) string {
			//	if index == 1 {
			//		return "A basic project structure"
			//	}
			//	return "It has rich functions such as db, jwt, cron, migration, test, etc"
			//},
		}
		err := survey.AskOne(prompt, &protocol)
		if err != nil {
			return false, err
		}
		registry := ""
		prompt2 := &survey.Select{
			Message: "Please select a registry:",
			Options: []string{
				"Github",
				"Gitee",
			},
		}
		err = survey.AskOne(prompt2, &registry)
		if err != nil {
			return false, err
		}

		switch registry {
		case "Github":
			repo = config.GitHubAddr
		case "Gitee":
			repo = config.GiteeAddr
		}

		switch protocol {
		case "MQTT":
			repo = repo + config.MqttProtocolDriver
		case "TCP":
			repo = repo + config.TcpProtocolDriver
		case "UDP":
			repo = repo + config.UdpProtocolDriver
		case "CoAP":
			repo = repo + config.CoapProtocolDriver
		case "HTTP":
			repo = repo + config.HttpProtocolDriver
		case "WebSocket":
			repo = repo + config.WebSocketProtocolDriver
		case "Modbus-TCP":
			repo = repo + config.ModbusProtocolDriver
		}
		err = os.RemoveAll(p.ProjectName)
		if err != nil {
			fmt.Println("remove old project error: ", err)
			return false, err
		}
	} else {
		repo = repoURL
	}
	fmt.Printf("git clone %s\n", repo)
	cmd := exec.Command("git", "clone", repo, p.ProjectName)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("git clone %s error: %s\n", repo, err)
		return false, err
	}
	return true, nil
}

func (p *Project) replacePackageName() error {
	packageName := getProjectName(p.ProjectName)

	err := p.replaceFiles(packageName)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "edit", "-module", p.ProjectName)
	cmd.Dir = p.ProjectName
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("go mod edit error: ", err)
		return err
	}
	return nil
}
func (p *Project) modTidy() error {
	fmt.Println("go mod tidy")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = p.ProjectName
	if err := cmd.Run(); err != nil {
		fmt.Println("go mod tidy error: ", err)
		return err
	}
	return nil
}
func (p *Project) rmGit() {
	os.RemoveAll(p.ProjectName + "/.git")
}

func (p *Project) replaceFiles(packageName string) error {
	err := filepath.Walk(p.ProjectName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		newData := bytes.ReplaceAll(data, []byte(packageName), []byte(p.ProjectName))
		if err := os.WriteFile(path, newData, 0644); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("walk file error: ", err)
		return err
	}
	return nil
}

func getProjectName(dir string) string {
	modFile, err := os.Open(dir + "/go.mod")
	if err != nil {
		fmt.Println("go.mod does not exist", err)
		return ""
	}
	defer modFile.Close()

	var moduleName string
	_, err = fmt.Fscanf(modFile, "module %s", &moduleName)
	if err != nil {
		fmt.Println("read go mod error: ", err)
		return ""
	}
	return moduleName
}
