package tf6muxserver_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/tf6testserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

func TestMuxServerReadResource(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	testServer1 := &tf6testserver.TestServer{
		ResourceSchemas: map[string]*tfprotov6.Schema{
			"test_resource_server1": {},
		},
	}
	testServer2 := &tf6testserver.TestServer{
		ResourceSchemas: map[string]*tfprotov6.Schema{
			"test_resource_server2": {},
		},
	}

	servers := []func() tfprotov6.ProviderServer{testServer1.ProviderServer, testServer2.ProviderServer}
	muxServer, err := tf6muxserver.NewMuxServer(ctx, servers...)

	if err != nil {
		t.Fatalf("unexpected error setting up factory: %s", err)
	}

	_, err = muxServer.ProviderServer().ReadResource(ctx, &tfprotov6.ReadResourceRequest{
		TypeName: "test_resource_server1",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !testServer1.ReadResourceCalled["test_resource_server1"] {
		t.Errorf("expected test_resource_server1 ReadResource to be called on server1")
	}

	if testServer2.ReadResourceCalled["test_resource_server1"] {
		t.Errorf("unexpected test_resource_server1 ReadResource called on server2")
	}

	_, err = muxServer.ProviderServer().ReadResource(ctx, &tfprotov6.ReadResourceRequest{
		TypeName: "test_resource_server2",
	})

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if testServer1.ReadResourceCalled["test_resource_server2"] {
		t.Errorf("unexpected test_resource_server2 ReadResource called on server1")
	}

	if !testServer2.ReadResourceCalled["test_resource_server2"] {
		t.Errorf("expected test_resource_server2 ReadResource to be called on server2")
	}
}
