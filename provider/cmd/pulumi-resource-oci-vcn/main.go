//go:generate go run ./generate.go

package main

import (
	"github.com/lbrlabs/pulumi-oci-vcn/pkg/provider"
	"github.com/lbrlabs/pulumi-oci-vcn/pkg/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

var providerName = "oci-vcn"

func main() {
	kingpin.Version(version.Version)
	kingpin.Parse()
	provider.Serve(providerName, version.Version, pulumiSchema)
}
