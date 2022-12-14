//go:generate go run ./generate.go

package main

import (
	"github.com/lbrlabs/pulumi-oci-vcn/pkg/provider"
	"github.com/lbrlabs/pulumi-oci-vcn/pkg/version"
)

var providerName = "oci-vcn"

func main() {
	provider.Serve(providerName, version.Version, pulumiSchema)
}
