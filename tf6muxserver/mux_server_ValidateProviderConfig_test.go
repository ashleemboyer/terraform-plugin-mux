package tf6muxserver_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf6testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func TestMuxServerValidateProviderConfig(t *testing.T) {
	t.Parallel()

	config, err := tfprotov6.NewDynamicValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"hello": tftypes.String,
		},
	}, tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"hello": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"hello": tftypes.NewValue(tftypes.String, "world"),
	}))

	if err != nil {
		t.Fatalf("error constructing config: %s", err)
	}

	testCases := map[string]struct {
		servers          []func() tfprotov6.ProviderServer
		expectedError    error
		expectedResponse *tfprotov6.ValidateProviderConfigResponse
	}{
		"error-once": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityError,
						Summary:  "test error summary",
						Detail:   "test error details",
					},
				},
			},
		},
		"error-multiple": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				}).ProviderServer,
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityError,
						Summary:  "test error summary",
						Detail:   "test error details",
					},
					{
						Severity: tfprotov6.DiagnosticSeverityError,
						Summary:  "test error summary",
						Detail:   "test error details",
					},
				},
			},
		},
		"warning-once": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityWarning,
						Summary:  "test warning summary",
						Detail:   "test warning details",
					},
				},
			},
		},
		"warning-multiple": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				}).ProviderServer,
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityWarning,
						Summary:  "test warning summary",
						Detail:   "test warning details",
					},
					{
						Severity: tfprotov6.DiagnosticSeverityWarning,
						Summary:  "test warning summary",
						Detail:   "test warning details",
					},
				},
			},
		},
		"warning-then-error": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				}).ProviderServer,
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityWarning,
						Summary:  "test warning summary",
						Detail:   "test warning details",
					},
					{
						Severity: tfprotov6.DiagnosticSeverityError,
						Summary:  "test error summary",
						Detail:   "test error details",
					},
				},
			},
		},
		"no-response": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
			},
		},
		"PreparedConfig-once": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				PreparedConfig: &config,
			},
		},
		"PreparedConfig-once-and-error": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "test error summary",
								Detail:   "test error details",
							},
						},
					},
				}).ProviderServer,
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityError,
						Summary:  "test error summary",
						Detail:   "test error details",
					},
				},
				PreparedConfig: &config,
			},
		},
		"PreparedConfig-once-and-warning": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityWarning,
								Summary:  "test warning summary",
								Detail:   "test warning details",
							},
						},
					},
				}).ProviderServer,
			},
			expectedResponse: &tfprotov6.ValidateProviderConfigResponse{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityWarning,
						Summary:  "test warning summary",
						Detail:   "test warning details",
					},
				},
				PreparedConfig: &config,
			},
		},
		"PreparedConfig-multiple": {
			servers: []func() tfprotov6.ProviderServer{
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				}).ProviderServer,
				(&tf6testserver.TestServer{}).ProviderServer,
				(&tf6testserver.TestServer{
					ValidateProviderConfigResponse: &tfprotov6.ValidateProviderConfigResponse{
						PreparedConfig: &config,
					},
				}).ProviderServer,
			},
			expectedError: fmt.Errorf("got a ValidateProviderConfig PreparedConfig response from multiple servers, not sure which to use"),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			muxServer, err := tf6muxserver.NewMuxServer(context.Background(), testCase.servers...)

			if err != nil {
				t.Fatalf("error setting up muxer: %s", err)
			}

			got, err := muxServer.ProviderServer().ValidateProviderConfig(context.Background(), &tfprotov6.ValidateProviderConfigRequest{
				Config: &config,
			})

			if err != nil {
				if testCase.expectedError == nil {
					t.Fatalf("wanted no error, got error: %s", err)
				}

				if !strings.Contains(err.Error(), testCase.expectedError.Error()) {
					t.Fatalf("wanted error %q, got error: %s", testCase.expectedError.Error(), err.Error())
				}
			}

			if err == nil && testCase.expectedError != nil {
				t.Fatalf("got no error, wanted err: %s", testCase.expectedError)
			}

			if !cmp.Equal(got, testCase.expectedResponse) {
				t.Errorf("unexpected response: %s", cmp.Diff(got, testCase.expectedResponse))
			}

			for num, server := range testCase.servers {
				if !server().(*tf6testserver.TestServer).ValidateProviderConfigCalled {
					t.Errorf("ValidateProviderConfig not called on server%d", num+1)
				}
			}
		})
	}
}
