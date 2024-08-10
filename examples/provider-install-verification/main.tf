terraform {
  required_providers {
    calico = {
      source = "jifwin/calico"
    }
  }
}

provider "calico" {}

data "calico_coffees" "example" {}
