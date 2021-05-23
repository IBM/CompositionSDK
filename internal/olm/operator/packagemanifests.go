// Copyright 2019 The Operator-SDK Authors
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

package olm

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/pflag"
)

// PackageManifestsCmd configures deployment and teardown of an operator
// managed in a package manifests format via OLM.
type PackageManifestsCmd struct {
	OperatorCmd

	// ManifestsDir is a directory containing 1..N package directories and
	// a package manifest.
	// OperatorVersion can be set to the version of the desired operator package
	// and Run()/Cleanup() will deploy that operator version.
	ManifestsDir string
	// OperatorVersion is the version of the operator to deploy. It must be
	// a semantic version, ex. 0.0.1.
	OperatorVersion string
	// IncludePaths are path to manifests of Kubernetes resources that either
	// supplement or override defaults generated by methods of PackageManifestsCmd. These
	// manifests can be but are not limited to: RBAC, Subscriptions,
	// CatalogSources, OperatorGroups.
	//
	// Kinds that are overridden if supplied:
	// - CatalogSource
	// - Subscription
	// - OperatorGroup
	IncludePaths []string
}

func (c *PackageManifestsCmd) AddToFlagSet(fs *pflag.FlagSet) {
	c.OperatorCmd.AddToFlagSet(fs)

	fs.StringVar(&c.OperatorVersion, "operator-version", "",
		"Version of operator to deploy")
	fs.StringSliceVar(&c.IncludePaths, "include", nil,
		"Path to Kubernetes resource manifests, ex. Role, Subscription. "+
			"These supplement or override defaults generated by run/cleanup")
}

func (c *PackageManifestsCmd) validate() error {
	if c.ManifestsDir == "" {
		return errors.New("manifests dir must be set")
	}
	if c.OperatorVersion == "" {
		return errors.New("operator version must be set")
	}
	return c.OperatorCmd.validate()
}

func (c *PackageManifestsCmd) initialize() {
	c.OperatorCmd.initialize()
}

func (c *PackageManifestsCmd) Run() error {
	c.initialize()
	if err := c.validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	m, err := c.newManager()
	if err != nil {
		return fmt.Errorf("error initializing operator manager: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	return m.run(ctx)
}

func (c *PackageManifestsCmd) Cleanup() (err error) {
	c.initialize()
	if err := c.validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	m, err := c.newManager()
	if err != nil {
		return fmt.Errorf("error initializing operator manager: %w", err)
	}
	// Cleanups should clean up all resources, which includes the registry.
	m.forceRegistry = true
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	return m.cleanup(ctx)
}
