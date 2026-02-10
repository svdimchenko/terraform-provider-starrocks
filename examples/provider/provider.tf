terraform {
  required_providers {
    starrocks = {
      source = "svdimchenko/starrocks"
    }
  }
}

provider "starrocks" {
  host     = "localhost"
  port     = 9030
  username = "root"
  password = "password"
}
