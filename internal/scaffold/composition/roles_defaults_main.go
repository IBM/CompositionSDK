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
	"github.com/operator-framework/operator-sdk/internal/scaffold/ansible"
	"path/filepath"
	"strings"

	"github.com/operator-framework/operator-sdk/internal/scaffold"
	"github.com/operator-framework/operator-sdk/internal/scaffold/input"

	"gopkg.in/yaml.v2"
)

const RolesDefaultsMainFile = "defaults" + ansible.FilePathSep + "main.yaml"

type RolesDefaultsMain struct {
	input.Input
	Resource         scaffold.Resource
	Labels           interface{}
	Properties       interface{}
	LabelsString     string
	LabelsExist      bool
	PropertiesString string
	PropertiesExist  bool
}

// GetInput - gets the input
func (r *RolesDefaultsMain) GetInput() (input.Input, error) {
	if r.Path == "" {
		r.Path = filepath.Join(ansible.RolesDir, r.Resource.LowerKind, RolesDefaultsMainFile)
	}
	r.TemplateBody = rolesDefaultsMainAnsibleTmpl
	r.Delims = ansible.AnsibleDelims

	labelsYaml := convertToYaml(r.Labels)
	r.LabelsExist = len(labelsYaml) > 0
	r.LabelsString = getYamlString(labelsYaml)

	propertiesYaml := convertToYaml(r.Properties)
	r.PropertiesExist = len(propertiesYaml) > 0
	r.PropertiesString = getYamlString(propertiesYaml)

	return r.Input, nil
}

func convertToYaml(mapping interface{}) map[interface{}]interface{} {
	yamlInput, _ := yaml.Marshal(mapping)
	var output map[interface{}]interface{}
	_ = yaml.Unmarshal(yamlInput, &output)
	return output
}

func getYamlString(yamlMap map[interface{}]interface{}) string {
	data, _ := yaml.Marshal(&yamlMap)
	str := string(data)
	str = strings.ReplaceAll(str, "\n", "\n    ")
	return str
}

const rolesDefaultsMainAnsibleTmpl = `---
# defaults file for [[.Resource.LowerKind]]
defaults:
  [[- if .LabelsExist ]]
  labels:
    [[.LabelsString]]
  [[- end ]]
  [[- if .PropertiesExist ]]
  properties:
    [[.PropertiesString]]
  [[- end ]]
`
