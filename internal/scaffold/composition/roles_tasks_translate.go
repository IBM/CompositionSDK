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

	"github.com/operator-framework/operator-sdk/internal/scaffold"
	"github.com/operator-framework/operator-sdk/internal/scaffold/input"
	"strings"
)

const RolesTasksTranslateFile = "tasks" + ansible.FilePathSep + "translate.yaml"

type RolesTasksTranslate struct {
	input.Input
	Resource           scaffold.Resource
	LowerCaseFullGroup string
}

// GetInput - gets the input
func (r *RolesTasksTranslate) GetInput() (input.Input, error) {
	if r.Path == "" {
		r.Path = filepath.Join(ansible.RolesDir, r.Resource.LowerKind, RolesTasksTranslateFile)
	}
	r.LowerCaseFullGroup = strings.ToLower(strings.ReplaceAll(r.Resource.FullGroup, ".", "_"))
	r.TemplateBody = rolesTasksTranslateAnsibleTmpl
	r.Delims = ansible.AnsibleDelims

	return r.Input, nil
}

const rolesTasksTranslateAnsibleTmpl = `---
- name: create labels dictionary
  set_fact:
    labels: {}

- name: merge default labels
  set_fact:
    labels: "{{ labels | combine( defaults.labels ) }}"
  when: defaults is defined and defaults.labels is defined

- name: merge cr labels
  set_fact:
    labels: "{{ labels | combine( _[[.LowerCaseFullGroup]]_[[.Resource.LowerKind]].metadata.labels ) }}"
  when: _[[.LowerCaseFullGroup]]_[[.Resource.LowerKind]].metadata.labels is defined

- name: create properties dictionary
  set_fact:
    properties: {}

- name: merge default properties
  set_fact:
    properties: "{{ properties | combine( defaults.properties ) }}"
  when: defaults is defined and defaults.properties is defined

- name: merge cr properties
  set_fact:
    properties: "{{ properties | combine( _[[.LowerCaseFullGroup]]_[[.Resource.LowerKind]]_spec ) }}"
  when: _[[.LowerCaseFullGroup]]_[[.Resource.LowerKind]]_spec is defined

- name: generate network service from the template
  set_fact:
    nsvc_cr: "{{ lookup('template', 'network_service.yaml.j2') }}"

- name: k8s apply network service CR list
  k8s:
    definition: "{{ nsvc_cr }}"
    apply: yes
`
