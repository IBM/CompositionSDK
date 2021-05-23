// Copyright 2018 The Operator-SDK Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package composition

import (
	"fmt"
	"github.com/operator-framework/operator-sdk/internal/scaffold/composition"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/operator-framework/operator-sdk/cmd/operator-sdk/internal/genutil"
	"github.com/operator-framework/operator-sdk/internal/flags/apiflags"
	"github.com/operator-framework/operator-sdk/internal/scaffold"
	"github.com/operator-framework/operator-sdk/internal/scaffold/ansible"
	"github.com/operator-framework/operator-sdk/internal/scaffold/input"
	"github.com/operator-framework/operator-sdk/internal/util/projutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func NewCmd() *cobra.Command { //nolint:golint
	/*
		The nolint here is used to hide the warning
		"func name will be used as new.NewCmd by other packages,
		and that stutters; consider calling this Cmd"
		which is a false positive.
	*/
	newCmd := &cobra.Command{
		Use:   "composition <project-name>",
		Short: "Creates a new composition based operator application",
		Long: `The operator-sdk new command creates a new operator application and
generates a default directory layout based on the input <project-name>.

<project-name> is the project name of the new operator. (e.g app-operator)
`,
		Example: `  # Create a new project directory
  $ mkdir $HOME/projects/example.com/
  $ cd $HOME/projects/example.com/


  # Ansible project
  $ operator-sdk composition test-operator \
    --api-version=pingpong.example.com/v1alpha1 \
    --kind=Pingpong \
    --generate-playbook \
    --nsvc-template=path/to/template
`,
		RunE: newFunc,
	}
	newCmd.Flags().BoolVar(&gitInit, "git-init", false,
		"Initialize the project directory as a git repository (default false)")
	newCmd.Flags().BoolVar(&generatePlaybook, "generate-playbook", false,
		"Generate a playbook skeleton. (Only used for --type ansible)")

	newCmd.Flags().StringVar(&nsvcTemplateFile, "nsvc-template", "",
		"Network service yaml template to be used.")

	// Initialize flagSet struct with common flags
	apiFlags.AddTo(newCmd.Flags())

	return newCmd
}

var (
	apiFlags         apiflags.APIFlags
	projectName      string
	gitInit          bool
	generatePlaybook bool
	nsvcTemplateFile string
)

func newFunc(cmd *cobra.Command, args []string) error {
	if err := parse(cmd, args); err != nil {
		return err
	}
	mustBeNewProject()
	if err := verifyFlags(); err != nil {
		return err
	}

	log.Infof("Creating new %s operator '%s'.", strings.Title(projutil.OperatorTypeAnsible), projectName)
	if err := doAnsibleScaffold(); err != nil {
		log.Fatal(err)
	}
	if gitInit {
		if err := initGit(); err != nil {
			log.Fatal(err)
		}
	}
	log.Info("Project creation complete.")
	return nil
}

func parse(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("command %s requires exactly one argument", cmd.CommandPath())
	}
	projectName = args[0]
	if len(projectName) == 0 {
		return fmt.Errorf("project name must not be empty")
	}
	return nil
}

// mustBeNewProject checks if the given project exists under the current diretory.
// it exits with error when the project exists.
func mustBeNewProject() {
	fp := filepath.Join(projutil.MustGetwd(), projectName)
	stat, err := os.Stat(fp)
	if err != nil && os.IsNotExist(err) {
		return
	}
	if err != nil {
		log.Fatalf("Failed to determine if project (%v) exists", projectName)
	}
	if stat.IsDir() {
		log.Fatalf("Project (%v) in (%v) path already exists. Please use a different project name or delete "+
			"the existing one", projectName, fp)
	}
}

func doAnsibleScaffold() error {
	cfg := &input.Config{
		AbsProjectPath: filepath.Join(projutil.MustGetwd(), projectName),
		ProjectName:    projectName,
	}

	resource, err := scaffold.NewResource(apiFlags.APIVersion, apiFlags.Kind)
	if err != nil {
		return fmt.Errorf("invalid apiVersion and kind: %v", err)
	}

	nsvcTemplate, err := nsvcTemplate(nsvcTemplateFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	roleFiles := ansible.RolesFiles{Resource: *resource}
	roleTemplates := ansible.RolesTemplates{Resource: *resource}

	s := &scaffold.Scaffold{}
	err = s.Execute(cfg,
		// build dir
		&composition.BuildDockerfile{GeneratePlaybook: generatePlaybook},
		&composition.BuildInstallfile{Resource: *resource},
		&composition.BuildCleanfile{Resource: *resource},
		// deploy dir
		&scaffold.CR{Resource: resource},
		&scaffold.ServiceAccount{},
		&composition.Role{},
		&scaffold.RoleBinding{},
		&composition.DeployOperator{},
		// roles dir
		&composition.RolesDefaultsMain{Resource: *resource, Labels: get(get(nsvcTemplate, "metadata"), "labels"), Properties: get(get(nsvcTemplate, "spec"), "properties")},
		&roleFiles,
		&ansible.RolesHandlersMain{Resource: *resource},
		&ansible.RolesMetaMain{Resource: *resource},
		&composition.RolesTasksMain{Resource: *resource},
		&composition.RolesTasksTranslate{Resource: *resource},
		&composition.RolesTasksUpdateStatus{Resource: *resource},
		&roleTemplates,
		&composition.RolesTemplatesNetworkService{Resource: *resource, Template: nsvcTemplate},
		&composition.RolesVarsMain{Resource: *resource},
		&ansible.RolesReadme{Resource: *resource},
		// root dir
		&composition.AnsibleCfg{},
		&composition.Makefile{},
		&composition.Readme{},
		&ansible.RequirementsYml{},
		&composition.Watches{
			GeneratePlaybook: generatePlaybook,
			Resource:         *resource,
		},
	)
	if err != nil {
		return fmt.Errorf("new ansible scaffold failed: %v", err)
	}

	if err = genutil.GenerateCRDNonGo(projectName, *resource, apiFlags.CrdVersion); err != nil {
		return err
	}

	// Remove placeholders from empty directories
	err = os.Remove(filepath.Join(s.AbsProjectPath, roleFiles.Path))
	if err != nil {
		return fmt.Errorf("new ansible scaffold failed: %v", err)
	}
	err = os.Remove(filepath.Join(s.AbsProjectPath, roleTemplates.Path))
	if err != nil {
		return fmt.Errorf("new ansible scaffold failed: %v", err)
	}

	// Decide on playbook.
	if generatePlaybook {
		log.Infof("Generating %s playbook.", strings.Title(projutil.OperatorTypeAnsible))

		err := s.Execute(cfg,
			&ansible.Playbook{Resource: *resource},
		)
		if err != nil {
			return fmt.Errorf("new ansible playbook scaffold failed: %v", err)
		}
	}

	// update deploy/role.yaml for the given resource r.
	if err := scaffold.UpdateRoleForResource(resource, cfg.AbsProjectPath); err != nil {
		return fmt.Errorf("failed to update the RBAC manifest for the resource (%v, %v): %v",
			resource.APIVersion, resource.Kind, err)
	}
	return nil
}

func verifyFlags() error {
	if err := apiFlags.VerifyCommonFlags(projutil.OperatorTypeAnsible); err != nil {
		return err
	}

	return nil
}

func execProjCmd(cmd string, args ...string) error {
	dc := exec.Command(cmd, args...)
	dc.Dir = filepath.Join(projutil.MustGetwd(), projectName)
	return projutil.ExecCmd(dc)
}

func initGit() error {
	log.Info("Running git init")
	if err := execProjCmd("git", "init"); err != nil {
		return fmt.Errorf("failed to run git init: %v", err)
	}
	log.Info("Run git init done")
	return nil
}

func nsvcTemplate(templateFile string) (map[string]interface{}, error) {
	yamlFile, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}
	var template map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &template)
	if err != nil {
		return nil, err
	}
	return template, nil
}
func get(this interface{}, key string) interface{} {
	yamlInput, err := yaml.Marshal(this)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	var mapping map[string]interface{}
	err = yaml.Unmarshal(yamlInput, &mapping)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if val, ok := mapping[key]; ok {
		return val
	}
	return map[string]interface{}{}
}
