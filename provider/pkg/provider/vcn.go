package provider

import (
	"fmt"
	ocicore "github.com/pulumi/pulumi-oci/sdk/go/oci/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The set of arguments for creating a StaticPage component resource.
type VcnArgs struct {
	CompartmentID pulumi.StringInput `pulumi:"compartmentId"`
	CidrBlock     string             `pulumi:"cidrBlock"`
}

// The Vcn component resource.
type Vcn struct {
	pulumi.ResourceState

	Vcn *ocicore.Vcn
}

// NewVcn creates a new Vcn component resource.
func NewVcn(ctx *pulumi.Context,
	name string, args *VcnArgs, opts ...pulumi.ResourceOption) (*Vcn, error) {
	if args == nil {
		args = &VcnArgs{}
	}

	var err error

	component := &Vcn{}

	vcn, err := ocicore.NewVcn(ctx, name, &ocicore.VcnArgs{
		CompartmentId: args.CompartmentID,
		CidrBlocks: pulumi.StringArray{
			pulumi.String(args.CidrBlock),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating vcn: %v", err)
	}

	err = ctx.RegisterComponentResource("oci-vcn:index:Vcn", name, component, opts...)
	if err != nil {
		return nil, err
	}

	component.Vcn = vcn
	// component.WebsiteUrl = bucket.WebsiteEndpoint

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{
		"vcn": vcn,
	}); err != nil {
		return nil, err
	}

	return component, nil
}
