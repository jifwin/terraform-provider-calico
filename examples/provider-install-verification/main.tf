terraform {
  required_providers {
    calico = {
      source = "jifwin/calico"
    }
  }
}

provider "calico" {
  kubeconfig = file("/Users/gpietrus/tmp1/test-cluster.yml")
}

#TODO: this works when the calico apiserver is installed
#TODO: two options
  # install the calico apiserver - prefered
  # use the exactly same code that calicoctl uses (not the official client)

data "calico_coffees" "example" {}

output "test" {
  value = data.calico_coffees.example
}
