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

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package vsphere

import (
	"github.com/elastic/beats/v7/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "vsphere", asset.ModuleFieldsPri, AssetVsphere); err != nil {
		panic(err)
	}
}

// AssetVsphere returns asset data.
// This is the base64 encoded zlib format compressed contents of module/vsphere.
func AssetVsphere() string {
	return "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"
}
