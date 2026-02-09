package client

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Client struct {
	db *sql.DB
}

type ResourceGroup struct {
	Name                     types.String
	CPUCoreLimit             types.Int64
	MemLimit                 types.String
	ConcurrencyLimit         types.Int64
	BigQueryMemLimit         types.String
	BigQueryScanRowsLimit    types.Int64
	BigQueryCPUSecondLimit   types.Int64
	Classifiers              types.Set
}

type Classifier struct {
	ID        int64
	User      types.String
	Role      types.String
	QueryType types.String
	SourceIP  types.String
	DB        types.String
}

type ResourceGroupModel interface {
	GetName() types.String
	GetCPUCoreLimit() types.Int64
	GetMemLimit() types.String
	GetConcurrencyLimit() types.Int64
	GetBigQueryMemLimit() types.String
	GetBigQueryScanRowsLimit() types.Int64
	GetBigQueryCPUSecondLimit() types.Int64
	GetClassifiers() types.Set
}

func NewClient(host, username, password string) (*Client, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", username, password, host)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &Client{db: db}, nil
}

func (c *Client) CreateResourceGroup(rg ResourceGroupModel) error {
	var parts []string
	parts = append(parts, fmt.Sprintf("CREATE RESOURCE GROUP '%s'", rg.GetName().ValueString()))

	if !rg.GetCPUCoreLimit().IsNull() {
		parts = append(parts, fmt.Sprintf("'cpu_core_limit' = '%d'", rg.GetCPUCoreLimit().ValueInt64()))
	}
	if !rg.GetMemLimit().IsNull() {
		parts = append(parts, fmt.Sprintf("'mem_limit' = '%s'", rg.GetMemLimit().ValueString()))
	}
	if !rg.GetConcurrencyLimit().IsNull() {
		parts = append(parts, fmt.Sprintf("'concurrency_limit' = '%d'", rg.GetConcurrencyLimit().ValueInt64()))
	}
	if !rg.GetBigQueryMemLimit().IsNull() {
		parts = append(parts, fmt.Sprintf("'big_query_mem_limit' = '%s'", rg.GetBigQueryMemLimit().ValueString()))
	}
	if !rg.GetBigQueryScanRowsLimit().IsNull() {
		parts = append(parts, fmt.Sprintf("'big_query_scan_rows_limit' = '%d'", rg.GetBigQueryScanRowsLimit().ValueInt64()))
	}
	if !rg.GetBigQueryCPUSecondLimit().IsNull() {
		parts = append(parts, fmt.Sprintf("'big_query_cpu_second_limit' = '%d'", rg.GetBigQueryCPUSecondLimit().ValueInt64()))
	}

	query := parts[0]
	if len(parts) > 1 {
		query += " WITH (" + strings.Join(parts[1:], ", ") + ")"
	}

	if _, err := c.db.Exec(query); err != nil {
		return err
	}

	// Add classifiers if provided
	if !rg.GetClassifiers().IsNull() && len(rg.GetClassifiers().Elements()) > 0 {
		// Parse and add classifiers
	}

	return nil
}

func (c *Client) GetResourceGroup(name string) (*ResourceGroup, error) {
	query := fmt.Sprintf("SHOW RESOURCE GROUP '%s'", name)
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rg := &ResourceGroup{Name: types.StringValue(name)}
	var classifiers []Classifier

	for rows.Next() {
		var name, id, cpuWeight, exclusiveCPUCores, memLimit, bigQueryCPUSecondLimit, bigQueryScanRowsLimit, bigQueryMemLimit, concurrencyLimit, spillMemLimitThreshold, classifiersStr string
		if err := rows.Scan(&name, &id, &cpuWeight, &exclusiveCPUCores, &memLimit, &bigQueryCPUSecondLimit, &bigQueryScanRowsLimit, &bigQueryMemLimit, &concurrencyLimit, &spillMemLimitThreshold, &classifiersStr); err != nil {
			return nil, err
		}

		if rg.MemLimit.IsNull() {
			rg.MemLimit = types.StringValue(memLimit)
		}
		if rg.ConcurrencyLimit.IsNull() {
			if v, err := strconv.ParseInt(concurrencyLimit, 10, 64); err == nil {
				rg.ConcurrencyLimit = types.Int64Value(v)
			}
		}
		if rg.BigQueryMemLimit.IsNull() {
			rg.BigQueryMemLimit = types.StringValue(bigQueryMemLimit)
		}
		if rg.BigQueryScanRowsLimit.IsNull() {
			if v, err := strconv.ParseInt(bigQueryScanRowsLimit, 10, 64); err == nil {
				rg.BigQueryScanRowsLimit = types.Int64Value(v)
			}
		}
		if rg.BigQueryCPUSecondLimit.IsNull() {
			if v, err := strconv.ParseInt(bigQueryCPUSecondLimit, 10, 64); err == nil {
				rg.BigQueryCPUSecondLimit = types.Int64Value(v)
			}
		}

		if classifiersStr != "" {
			classifier := parseClassifier(classifiersStr)
			classifiers = append(classifiers, classifier)
		}
	}

	return rg, nil
}

func parseClassifier(s string) Classifier {
	re := regexp.MustCompile(`id=(\d+).*?user=([^,)]+)|role=([^,)]+)|query_type=([^,)]+)|source_ip=([^,)]+)|db=([^,)]+)`)
	matches := re.FindStringSubmatch(s)
	c := Classifier{}
	if len(matches) > 1 {
		c.ID, _ = strconv.ParseInt(matches[1], 10, 64)
	}
	if len(matches) > 2 && matches[2] != "" {
		c.User = types.StringValue(matches[2])
	}
	if len(matches) > 3 && matches[3] != "" {
		c.Role = types.StringValue(matches[3])
	}
	if len(matches) > 4 && matches[4] != "" {
		c.QueryType = types.StringValue(matches[4])
	}
	if len(matches) > 5 && matches[5] != "" {
		c.SourceIP = types.StringValue(matches[5])
	}
	if len(matches) > 6 && matches[6] != "" {
		c.DB = types.StringValue(matches[6])
	}
	return c
}

func (c *Client) UpdateResourceGroup(rg ResourceGroupModel, oldRg ResourceGroupModel) error {
	var oldClassifiers, newClassifiers []Classifier

	// Drop old classifiers
	if !oldRg.GetClassifiers().IsNull() {
		// Parse old classifiers to get IDs for dropping
	}

	for _, classifier := range oldClassifiers {
		query := fmt.Sprintf("ALTER RESOURCE GROUP '%s' DROP (%d)", rg.GetName().ValueString(), classifier.ID)
		if _, err := c.db.Exec(query); err != nil {
			return err
		}
	}

	var parts []string

	if !rg.GetCPUCoreLimit().IsNull() {
		parts = append(parts, fmt.Sprintf("'cpu_core_limit' = '%d'", rg.GetCPUCoreLimit().ValueInt64()))
	}
	if !rg.GetMemLimit().IsNull() {
		parts = append(parts, fmt.Sprintf("'mem_limit' = '%s'", rg.GetMemLimit().ValueString()))
	}
	if !rg.GetConcurrencyLimit().IsNull() {
		parts = append(parts, fmt.Sprintf("'concurrency_limit' = '%d'", rg.GetConcurrencyLimit().ValueInt64()))
	}
	if !rg.GetBigQueryMemLimit().IsNull() {
		parts = append(parts, fmt.Sprintf("'big_query_mem_limit' = '%s'", rg.GetBigQueryMemLimit().ValueString()))
	}
	if !rg.GetBigQueryScanRowsLimit().IsNull() {
		parts = append(parts, fmt.Sprintf("'big_query_scan_rows_limit' = '%d'", rg.GetBigQueryScanRowsLimit().ValueInt64()))
	}
	if !rg.GetBigQueryCPUSecondLimit().IsNull() {
		parts = append(parts, fmt.Sprintf("'big_query_cpu_second_limit' = '%d'", rg.GetBigQueryCPUSecondLimit().ValueInt64()))
	}

	if len(parts) > 0 {
		query := fmt.Sprintf("ALTER RESOURCE GROUP '%s' WITH (%s)", rg.GetName().ValueString(), strings.Join(parts, ", "))
		if _, err := c.db.Exec(query); err != nil {
			return err
		}
	}

	// Add new classifiers
	for _, classifier := range newClassifiers {
		var conditions []string
		if !classifier.User.IsNull() {
			conditions = append(conditions, fmt.Sprintf("user='%s'", classifier.User.ValueString()))
		}
		if !classifier.Role.IsNull() {
			conditions = append(conditions, fmt.Sprintf("role='%s'", classifier.Role.ValueString()))
		}
		if !classifier.QueryType.IsNull() {
			conditions = append(conditions, fmt.Sprintf("query_type='%s'", classifier.QueryType.ValueString()))
		}
		if !classifier.SourceIP.IsNull() {
			conditions = append(conditions, fmt.Sprintf("source_ip='%s'", classifier.SourceIP.ValueString()))
		}
		if !classifier.DB.IsNull() {
			conditions = append(conditions, fmt.Sprintf("db='%s'", classifier.DB.ValueString()))
		}

		query := fmt.Sprintf("ALTER RESOURCE GROUP '%s' ADD (%s)", rg.GetName().ValueString(), strings.Join(conditions, ", "))
		if _, err := c.db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) DeleteResourceGroup(name string) error {
	query := fmt.Sprintf("DROP RESOURCE GROUP '%s'", name)
	_, err := c.db.Exec(query)
	return err
}
