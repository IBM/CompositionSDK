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

type Makefile struct {
	ansible.StaticInput
}

// GetInput - gets the input
func (m *Makefile) GetInput() (input.Input, error) {
	if m.Path == "" {
		m.Path = "Makefile"
	}
	m.TemplateBody = makefileTmpl
	return m.Input, nil
}

const makefileTmpl = `.PHONY: help				##show this help message
help:
	@echo "usage: make [target]\n"; echo "options:"; \fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//' | sed 's/.PHONY:*//' | sed -e 's/^/  /'; echo "";

.PHONY: docker-push			##pushes the operator to docker registry. in order to run it need to add <IMAGE> var as an argument.
docker-push:
	@operator-sdk build ${IMAGE}
	@docker push ${IMAGE}

.PHONY: install			##install the oeprator in the specified <NAMESPACE>.
install:
	@build/install.sh ${NAMESPACE}

.PHONY: clean			##cleans all objects from kubernetes in the specified <NAMESPACE>.
clean:
	@build/clean.sh ${NAMESPACE}

`
