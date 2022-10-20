package main

import (
	"github.com/lbrlabs/pulumi-oci-vcn/sdk/go/vcn"
	"github.com/pulumi/pulumi-oci/sdk/go/oci/Identity"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		compartmentId := Identity.GetCompartment(ctx, &identity.GetCompartmentArgs{
			Id: "ocid1.tenancy.oc1..aaaaaaaavh67gfytloujvijdue6pqwhfrlo2ms4jbkpozyez7sgdhf27e4yq",
		}, nil).Id
		vcn, err := vcn.NewVcn(ctx, "vcn", &vcn.VcnArgs{
			CompartmentId:         pulumi.String(compartmentId),
			CidrBlock:             "172.16.0.0/22",
			CreateInternetGateway: true,
			CreateNatGateway:      true,
			CreateServiceGateway:  true,
			DnsLabel:              pulumi.String("lbriggs"),
		})
		if err != nil {
			return err
		}
		ctx.Export("vcnId", vcn.VcnId)
		return nil
	})
}
