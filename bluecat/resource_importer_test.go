package bluecat

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestExternalHostRecordImporterSetsAbsoluteName(t *testing.T) {
	// A valid import uses the full External Host FQDN as both ID and absolute_name.
	resource := ResourceExternalHostRecord()
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	data.SetId("external.example.com")

	state, err := resource.Importer.State(data, nil)
	if err != nil {
		t.Fatalf("unexpected import error: %s", err)
	}
	if len(state) != 1 {
		t.Fatalf("expected one imported state, got %d", len(state))
	}

	imported := state[0]
	if got := imported.Id(); got != "external.example.com" {
		t.Fatalf("expected ID external.example.com, got %s", got)
	}
	if got := imported.Get("absolute_name").(string); got != "external.example.com" {
		t.Fatalf("expected absolute_name external.example.com, got %s", got)
	}
}

func TestExternalHostRecordImporterRejectsNonFQDN(t *testing.T) {
	// External Host imports need an FQDN because there is no separate zone field to infer.
	resource := ResourceExternalHostRecord()
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{})
	data.SetId("external")

	_, err := resource.Importer.State(data, nil)
	if err == nil {
		t.Fatal("expected import error for non-FQDN ID")
	}
	if !strings.Contains(err.Error(), "fully qualified domain name") {
		t.Fatalf("expected FQDN validation error, got %s", err)
	}
}

func TestGetAbsoluteNameFallsBackToAbsoluteNameWhenIDCannotBeParsed(t *testing.T) {
	// External Host Records do not depend on record + zone parsing when absolute_name is available.
	resource := ResourceExternalHostRecord()
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"absolute_name": "configured.example.com",
	})
	data.SetId("configured")

	got, err := getAbsoluteName(data)
	if err != nil {
		t.Fatalf("unexpected absolute_name error: %s", err)
	}
	if got != "configured.example.com" {
		t.Fatalf("expected configured.example.com, got %s", got)
	}
}
