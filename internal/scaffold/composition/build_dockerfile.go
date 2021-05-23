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
	"github.com/operator-framework/operator-sdk/version"
)

const BuildDockerfileFile = "Dockerfile"

type BuildDockerfile struct {
	input.Input
	RolesDir         string
	ImageTag         string
	GeneratePlaybook bool
}

// GetInput - gets the input
func (b *BuildDockerfile) GetInput() (input.Input, error) {
	if b.Path == "" {
		b.Path = filepath.Join(scaffold.BuildDir, BuildDockerfileFile)
	}
	b.TemplateBody = buildDockerfileAnsibleTmpl
	b.Delims = ansible.AnsibleDelims
	b.RolesDir = ansible.RolesDir
	b.ImageTag = strings.TrimSuffix(version.Version, "+git")
	return b.Input, nil
}

const buildDockerfileAnsibleTmpl = `FROM quay.io/operator-framework/ansible-operator:v0.16.0

COPY requirements.yaml ${HOME}/requirements.yaml
RUN ansible-galaxy collection install -r ${HOME}/requirements.yaml \
 && chmod -R ug+rwx ${HOME}/.ansible

COPY watches.yaml ${HOME}/watches.yaml

COPY [[.RolesDir]]/ ${HOME}/[[.RolesDir]]/
[[- if .GeneratePlaybook ]]
COPY playbook.yaml ${HOME}/playbook.yaml
[[- end ]]
COPY ansible.cfg ${HOME}/ansible.cfg
ENV ANSIBLE_CONFIG=${HOME}/ansible.cfg
`
