// Copyright 2020 The Operator-SDK Authors
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

const DeployOperatorFile = "operator.yaml.template"

type DeployOperator struct {
	input.Input
}

// GetInput - gets the input
func (d *DeployOperator) GetInput() (input.Input, error) {
	if d.Path == "" {
		d.Path = filepath.Join(scaffold.DeployDir, DeployOperatorFile)
	}
	d.TemplateBody = deployOperatorAnsibleTmpl
	d.Delims = ansible.AnsibleDelims

	return d.Input, nil
}

const deployOperatorAnsibleTmpl = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: [[.ProjectName]]
spec:
  replicas: 1
  selector:
    matchLabels:
      name: [[.ProjectName]]
  template:
    metadata:
      labels:
        name: [[.ProjectName]]
    spec:
      serviceAccountName: [[.ProjectName]]
      containers:
        - name: ansible
          command:
          - /usr/local/bin/ao-logs
          - /tmp/ansible-operator/runner
          - stdout
          # Replace this with the built image name
          image: "$IMAGE"
          imagePullPolicy: Always
          volumeMounts:
          - mountPath: /tmp/ansible-operator/runner
            name: runner
            readOnly: true
        - name: operator
          # Replace this with the built image name
          image: "$IMAGE"
          imagePullPolicy: Always
          args:
            - "--max-workers"
            - "10"
          volumeMounts:
          - mountPath: /tmp/ansible-operator/runner
            name: runner
          env:
            - name: WATCH_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: OPERATOR_NAME
              value: "[[.ProjectName]]"
            - name: ANSIBLE_GATHERING
              value: explicit
      volumes:
        - name: runner
          emptyDir: {}
`
