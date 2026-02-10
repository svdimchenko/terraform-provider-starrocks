package starrocks

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestResourceGroupModel_Getters(t *testing.T) {
	model := &resourceGroupResourceModel{
		Name:                   types.StringValue("test_rg"),
		CPUWeight:              types.Int64Value(10),
		MemLimit:               types.StringValue("80.0%"),
		ConcurrencyLimit:       types.Int64Value(5),
		BigQueryMemLimit:       types.Int64Value(1073741824),
		BigQueryScanRowsLimit:  types.Int64Value(100000),
		BigQueryCPUSecondLimit: types.Int64Value(100),
	}

	if model.GetName().ValueString() != "test_rg" {
		t.Errorf("GetName() = %v, want test_rg", model.GetName().ValueString())
	}
	if model.GetCPUWeight().ValueInt64() != 10 {
		t.Errorf("GetCPUWeight() = %v, want 10", model.GetCPUWeight().ValueInt64())
	}
	if model.GetMemLimit().ValueString() != "80.0%" {
		t.Errorf("GetMemLimit() = %v, want 80.0%%", model.GetMemLimit().ValueString())
	}
	if model.GetConcurrencyLimit().ValueInt64() != 5 {
		t.Errorf("GetConcurrencyLimit() = %v, want 5", model.GetConcurrencyLimit().ValueInt64())
	}
}
