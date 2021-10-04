// +build integ
// Copyright Istio Authors
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

package systemca

import (
	"testing"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/istio"

	// "istio.io/istio/tests/integration/pilot/common"
	"istio.io/istio/pkg/test/framework/components/prometheus"
	"istio.io/istio/pkg/test/framework/label"
	"istio.io/istio/pkg/test/framework/resource"
)

var prom prometheus.Instance

var i istio.Instance

func TestMain(m *testing.M) {
	framework.
		NewSuite(m).
		RequireSingleCluster().
		Label(label.CustomSetup).
		Setup(istio.Setup(&i, func(t resource.Context, cfg *istio.Config) {
			cfg.ControlPlaneValues = `
values:
  pilot:
    env:
      VERIFY_CERTIFICATE_AT_CLIENT: "true"`
		})).
		Setup(setup).
		Run()
}

func setup(ctx resource.Context) (err error) {
	prom, err = prometheus.New(ctx, prometheus.Config{})
	return err
}

func TestOutboundTrafficPolicy_AllowAny(t *testing.T) {
	cases := []*TestCase{
		{
			Name:     "HTTP Traffic Egress",
			PortName: "http",
			Host:     "some-external-site.com",
			Expected: Expected{
				ResponseCode: []string{"503"},
				Metadata:     map[string]string{},
			},
		},
		{
			Name:     "HTTP H2 Traffic Egress",
			PortName: "http",
			HTTP2:    true,
			Host:     "some-external-site.com",
			Expected: Expected{

				ResponseCode: []string{"503"},
				Metadata:     map[string]string{},
			},
		},
	}

	RunExternalRequest(cases, prom, AllowAny, t)
}
