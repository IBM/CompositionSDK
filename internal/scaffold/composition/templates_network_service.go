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
	"github.com/operator-framework/operator-sdk/internal/scaffold"
	"github.com/operator-framework/operator-sdk/internal/scaffold/ansible"
	"github.com/operator-framework/operator-sdk/internal/scaffold/input"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"strings"
)

const RolesTemplatesNetworkServiceFile = "templates" + ansible.FilePathSep + "network_service.yaml.j2"

type RolesTemplatesNetworkService struct {
	input.Input
	Resource       scaffold.Resource
	Template       map[string]interface{}
	TemplateString string
}

// GetInput - gets the input
func (r *RolesTemplatesNetworkService) GetInput() (input.Input, error) {
	if r.Path == "" {
		r.Path = filepath.Join(ansible.RolesDir, r.Resource.LowerKind, RolesTemplatesNetworkServiceFile)
	}
	r.modifyTemplate()
	data, _ := yaml.Marshal(&r.Template)
	r.TemplateString = string(data)
	r.TemplateString = strings.ReplaceAll(r.TemplateString, "$LABELS_PLACE_HOLDER", labelsTmpl)
	r.TemplateString = strings.ReplaceAll(r.TemplateString, "$PROPERTIES_PLACE_HOLDER", propertiesTmpl)
	r.TemplateBody = rolesTemplatesNetworkServiceTmpl
	r.Delims = ansible.AnsibleDelims

	return r.Input, nil
}

func (r *RolesTemplatesNetworkService) modifyTemplate() {
	// handle metadata
	metadata := parseAsYaml(r.Template["metadata"])
	metadata["name"] = "{{ meta.name }}"
	metadata["namespace"] = "{{ meta.namespace }}"
	metadata["labels"] = "$LABELS_PLACE_HOLDER"
	r.Template["metadata"] = metadata

	// handle properties
	spec := parseAsYaml(r.Template["spec"])
	spec["properties"] = "$PROPERTIES_PLACE_HOLDER"
	r.Template["spec"] = spec
}

func parseAsYaml(yamlMap interface{}) map[string]interface{} {
	yamlInput, _ := yaml.Marshal(yamlMap)
	var mapping map[string]interface{}
	yaml.Unmarshal(yamlInput, &mapping)
	return mapping
}

const rolesTemplatesNetworkServiceTmpl = `[[.TemplateString]]
`

const labelsTmpl = `
  {% for key in labels %}
  {{ key }}: {{ labels[key] }}
  {% endfor %}
`

const propertiesTmpl = `
  {% for key in properties %}
  {{ key }}: {{ properties[key] }}
  {% endfor %}
`
