// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package mongodb

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/elastic/beats/v7/metricbeat/mb"

	"github.com/elastic/elastic-agent-libs/transport/tlscommon"
)

func init() {
	// Register the ModuleFactory function for the "mongodb" module.
	if err := mb.Registry.AddModule("mongodb", NewModule); err != nil {
		panic(err)
	}
}

// ModuleConfig contains the common configuration for this module
type ModuleConfig struct {
	Hosts []string          `config:"hosts"    validate:"nonzero,required"`
	TLS   *tlscommon.Config `config:"ssl"`

	Database string `config:"database"`

	Username string `config:"username"`
	Password string `config:"password"`

	Credentials struct {
		AuthMechanism           string            `config:"auth_mechanism"`
		AuthMechanismProperties map[string]string `config:"auth_mechanism_properties"`
		AuthSource              string            `config:"auth_source"`
		PasswordSet             bool              `config:"password_set"`
	} `config:"credentials"`
}

type MetricSet struct {
	mb.BaseMetricSet
	Config        ModuleConfig
	ClientOptions *options.ClientOptions
}

type module struct {
	mb.BaseModule
}

// NewModule creates a new mb.Module instance and validates that at least one host has been
// specified
func NewModule(base mb.BaseModule) (mb.Module, error) {
	return &module{base}, nil
}

func NewMetricSet(base mb.BaseMetricSet) (*MetricSet, error) {
	// Validate that at least one host has been specified.
	config := ModuleConfig{}
	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, fmt.Errorf("could not read config: %w", err)
	}

	clientOptions := options.Client()

	// options.Credentials must be nil for the driver to work properly if no auth is provided.
	// Zero values breaks the connnection.
	if config.Username != "" && config.Password != "" {
		clientOptions.Auth = &options.Credential{
			AuthMechanism: config.Credentials.AuthMechanism,
			AuthSource:    config.Credentials.AuthSource,
			Username:      config.Username,
			Password:      config.Password,
			PasswordSet:   config.Credentials.PasswordSet,
		}

		// clientOptions.Auth.AuthMechanismProperties is the only field here that might be nil, be empty or filled.
		if config.Credentials.AuthMechanismProperties != nil {
			clientOptions.Auth.AuthMechanismProperties = config.Credentials.AuthMechanismProperties
		}
	}

	if config.TLS.IsEnabled() {
		tlsConfig, err := tlscommon.LoadTLSConfig(config.TLS)
		if err != nil {
			return nil, fmt.Errorf("could not load provided TLS configuration: %w", err)
		}

		clientOptions.SetTLSConfig(tlsConfig.ToConfig())
	}

	host := base.Host()
	if parts := strings.SplitN(host, "://", 2); len(parts) != 2 {
		host = "mongodb://" + host // add scheme
	}
	if err := clientOptions.ApplyURI(host).Validate(); err != nil {
		return nil, fmt.Errorf("URL validation error: %w", err)
	}

	readPreference, err := readpref.New(readpref.PrimaryMode)
	if err != nil {
		return nil, err
	}

	clientOptions.SetReadPreference(readPreference)

	// TODO(SS): Decide a better timeout? Currently, mb.DefaultModuleConfig().Timeout
	// will make the client use no timeout. If we do not call this method, then by default,
	// it is set to 30s.
	clientOptions.SetConnectTimeout(mb.DefaultModuleConfig().Timeout)

	if err := clientOptions.Validate(); err != nil {
		return nil, fmt.Errorf("client options validation failed: %w", err)
	}

	return &MetricSet{Config: config, ClientOptions: clientOptions, BaseMetricSet: base}, nil
}

func NewClient(opts *options.ClientOptions) (*mongo.Client, error) {
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("could not create mongodb client: %w", err)
	}

	if err = client.Connect(context.Background()); err != nil {
		return nil, fmt.Errorf("could not connect to mongodb: %w", err)
	}

	return client, nil
}
