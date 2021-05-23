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
	"path/filepath"
)

const BuildInstallFile = "install.sh"

type BuildInstallfile struct {
	input.Input
	CrdPath string
	// Resource defines the inputs for the new custom resource
	Resource scaffold.Resource
}

// GetInput - gets the input
func (b *BuildInstallfile) GetInput() (input.Input, error) {
	if b.Path == "" {
		b.Path = filepath.Join(scaffold.BuildDir, BuildInstallFile)
	}
	b.IsExec = true
	b.TemplateBody = buildInstallCompositionTmpl
	b.Delims = ansible.AnsibleDelims
	b.CrdPath = crdPathForResource(scaffold.CRDsDir, &b.Resource)
	return b.Input, nil
}

const buildInstallCompositionTmpl = `#! /usr/bin/env bash

NAMESPACE=${1:-default}

kubectl apply -f [[.CrdPath]]

kubectl apply -f deploy/service_account.yaml -n $NAMESPACE
kubectl apply -f deploy/role.yaml -n $NAMESPACE
kubectl apply -f deploy/role_binding.yaml -n $NAMESPACE
envsubst < deploy/operator.yaml.template | kubectl apply -n ${NAMESPACE} -f -
`
