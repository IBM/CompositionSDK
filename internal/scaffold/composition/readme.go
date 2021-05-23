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
	"github.com/operator-framework/operator-sdk/internal/scaffold/input"
)

type Readme struct {
	input.Input
}

// GetInput - gets the input
func (r *Readme) GetInput() (input.Input, error) {
	if r.Path == "" {
		r.Path = "README.md"
	}
	r.TemplateBody = readmeTmpl
	r.Delims = ansible.AnsibleDelims

	return r.Input, nil
}

const readmeTmpl = `# [[.ProjectName]]

## Install of the operator

1.  Create a namespace in Kubernetes

1.  Set the ` + "`NAMESPACE`" + ` environment variable to hold the name of your namespace.   

    ` + "```" + `
    $ export NAMESPACE=<YOUR NAMESPACE>
` + "    ```" + `

1.  Set the ` + "`REGISTRY`" + ` environment variable to hold the name of your docker registry.   

    ` + "```" + `
    $ export REGISTRY=<YOUR REGISTRY>
` + "    ```" + `

1.  Set the ` + "`IMAGE`" + ` environment variable to hold the image of the operator.

    ` + "```" + `
    $ export IMAGE=$REGISTRY/$(basename $(pwd)):v0.0.1
` + "    ```" + `

1.  Run make to push the image to your docker registry:
    ` + "```" + `
    $ make docker-push "IMAGE=$IMAGE"
` + "    ```" + `

1.  Run make to install the operator:

    ` + "```" + `
    $ make install "IMAGE=$IMAGE" "NAMESPACE=$NAMESPACE"
` + "    ```" + `

## Troubleshooting

Check the logs of the operator:

` + "```" + `
$ kubectl logs -l name=[[.ProjectName]] -c operator --tail=1000 -n $NAMESPACE
` + "```" + `

More information is in the logs of Ansible:

` + "```" + `
$ kubectl logs -l name=[[.ProjectName]] -c ansible --tail=1000 -n $NAMESPACE
` + "```" + `

### Cleanup

` + "```" + `
$ make clean "NAMESPACE=$NAMESPACE"
` + "```" + `
`
