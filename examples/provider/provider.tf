terraform {
  required_providers {
    starrocks = {
      source = "hashicorp/starrocks"
    }
  }
}

provider "starrocks" {
  host     = "localhost:9030"
  username = "root"
  password = "password"
}
