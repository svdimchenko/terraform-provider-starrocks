resource "starrocks_resource_group" "example" {
  name                        = "example_rg"
  cpu_core_limit              = 10
  mem_limit                   = "80%"
  concurrency_limit           = 11
  big_query_mem_limit         = "1073741824"
  big_query_scan_rows_limit   = 100000
  big_query_cpu_second_limit  = 100

  classifiers = [
    {
      user = "username"
    },
    {
      role = "admin"
    },
    {
      query_type = "SELECT"
      db         = "analytics"
    }
  ]
}
