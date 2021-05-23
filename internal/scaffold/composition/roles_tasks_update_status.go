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
)

const RolesTasksUpdateStatusFile = "tasks" + ansible.FilePathSep + "update_status.yaml"

type RolesTasksUpdateStatus struct {
	input.Input
	Resource scaffold.Resource
}

// GetInput - gets the input
func (r *RolesTasksUpdateStatus) GetInput() (input.Input, error) {
	if r.Path == "" {
		r.Path = filepath.Join(ansible.RolesDir, r.Resource.LowerKind, RolesTasksUpdateStatusFile)
	}
	r.TemplateBody = rolesTasksUpdateStatusAnsibleTmpl
	r.Delims = ansible.AnsibleDelims

	return r.Input, nil
}

const rolesTasksUpdateStatusAnsibleTmpl = `---
- name: Fetch a network service by name
  set_fact:
    nsvc: "{{ lookup('k8s', api_version=nsvc_api_version, kind=nsvc_kind, namespace=meta.namespace, resource_name=meta.name ) }}"

- name: update the resource status
  k8s_status:
    api_version: "{{ api_version }}"
    kind: "{{ kind }}"
    name: "{{ meta.name }}"
    namespace: "{{ meta.namespace }}"
    status: "{{ nsvc.status }}"
  when: nsvc is defined and nsvc.status is defined
`
