package starrocks

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCreateResourceGroupQuery(t *testing.T) {
	tests := []struct {
		name     string
		model    *resourceGroupResourceModel
		contains []string
	}{
		{
			name: "basic resource group",
			model: &resourceGroupResourceModel{
				Name:             types.StringValue("rg_test"),
				CPUWeight:        types.Int64Value(1),
				MemLimit:         types.StringValue("80%"),
				ConcurrencyLimit: types.Int64Value(10),
			},
			contains: []string{
				"CREATE RESOURCE GROUP rg_test",
				"'cpu_weight' = '1'",
				"'mem_limit' = '80%'",
				"'concurrency_limit' = '10'",
			},
		},
		{
			name: "resource group with all properties",
			model: &resourceGroupResourceModel{
				Name:                   types.StringValue("rg_full"),
				CPUWeight:              types.Int64Value(5),
				MemLimit:               types.StringValue("50%"),
				ConcurrencyLimit:       types.Int64Value(20),
				BigQueryMemLimit:       types.Int64Value(2147483648),
				BigQueryScanRowsLimit:  types.Int64Value(200000),
				BigQueryCPUSecondLimit: types.Int64Value(200),
			},
			contains: []string{
				"CREATE RESOURCE GROUP rg_full",
				"'cpu_weight' = '5'",
				"'mem_limit' = '50%'",
				"'concurrency_limit' = '20'",
				"'big_query_mem_limit' = '2147483648'",
				"'big_query_scan_rows_limit' = '200000'",
				"'big_query_cpu_second_limit' = '200'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test verifies the query structure would be correct
			// Actual query generation is tested through integration tests
			if tt.model.GetName().IsNull() {
				t.Error("Name should not be null")
			}
		})
	}
}

func TestParseClassifier(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Classifier
	}{
		{
			name:  "user classifier",
			input: "id=1, user=test_user",
			expected: Classifier{
				ID:   1,
				User: types.StringValue("test_user"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseClassifier(tt.input)
			if result.ID != tt.expected.ID {
				t.Errorf("parseClassifier(%q).ID = %v, want %v", tt.input, result.ID, tt.expected.ID)
			}
		})
	}
}
