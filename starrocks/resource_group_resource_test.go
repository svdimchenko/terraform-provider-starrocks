package starrocks

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFormatMemLimit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"percentage without decimal", "80%", "80.0%"},
		{"percentage with decimal", "80.5%", "80.5%"},
		{"percentage with zero decimal", "80.0%", "80.0%"},
		{"non-percentage", "1024", "1024"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatMemLimit(tt.input)
			if result != tt.expected {
				t.Errorf("formatMemLimit(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestResourceGroupModel_Getters(t *testing.T) {
	model := &resourceGroupResourceModel{
		Name:                   types.StringValue("test_rg"),
		CPUWeight:              types.Int64Value(10),
		MemLimit:               types.StringValue("80%"),
		ConcurrencyLimit:       types.Int64Value(5),
		BigQueryMemLimit:       types.StringValue("1073741824"),
		BigQueryScanRowsLimit:  types.Int64Value(100000),
		BigQueryCPUSecondLimit: types.Int64Value(100),
	}

	if model.GetName().ValueString() != "test_rg" {
		t.Errorf("GetName() = %v, want test_rg", model.GetName().ValueString())
	}
	if model.GetCPUWeight().ValueInt64() != 10 {
		t.Errorf("GetCPUWeight() = %v, want 10", model.GetCPUWeight().ValueInt64())
	}
	if model.GetMemLimit().ValueString() != "80%" {
		t.Errorf("GetMemLimit() = %v, want 80%%", model.GetMemLimit().ValueString())
	}
	if model.GetConcurrencyLimit().ValueInt64() != 5 {
		t.Errorf("GetConcurrencyLimit() = %v, want 5", model.GetConcurrencyLimit().ValueInt64())
	}
}
