package provider

import (
	"fmt"
	ocicore "github.com/pulumi/pulumi-oci/sdk/go/oci/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The set of arguments for creating a StaticPage component resource.
type VcnArgs struct {
	CompartmentID         pulumi.StringInput `pulumi:"compartmentId"`
	CidrBlock             string             `pulumi:"cidrBlock"`
	DNSLabel              pulumi.StringInput `pulumi:"dnsLabel"`
	Ipv6Enabled           pulumi.BoolInput   `pulumi:"ipv6Enabled"`
	NumberOfSubnets       int                `pulumi:"numberOfSubnets"`
	CreateInternetGateway bool               `pulumi:"createInternetGateway"`
	CreateServiceGateway  bool               `pulumi:"createServiceGateway"`
	CreateNATGateway      bool               `pulumi:"createNatGateway"`
}

// The Vcn component resource.
type Vcn struct {
	pulumi.ResourceState

	VcnID             pulumi.IDOutput          `pulumi:"vcnId"`
	PublicSubnetIDs   pulumi.StringArrayOutput `pulumi:"publicSubnetIds"`
	PrivateSubnetIDs  pulumi.StringArrayOutput `pulumi:"privateSubnetIds"`
	InternetGatewayID pulumi.IDOutput          `pulumi:"internetGatewayId"`
	ServiceGatewayID  pulumi.IDOutput          `pulumi:"serviceGatewayId"`
	NATGatewayID      pulumi.IDOutput          `pulumi:"NatGatewayId"`
}

// NewVcn creates a new Vcn component resource.
func NewVcn(ctx *pulumi.Context,
	name string, args *VcnArgs, opts ...pulumi.ResourceOption) (*Vcn, error) {
	if args == nil {
		args = &VcnArgs{}
	}

	var err error

	component := &Vcn{}

	err = ctx.RegisterComponentResource("oci-vcn:index:Vcn", name, component, opts...)
	if err != nil {
		return nil, err
	}

	vcn, err := ocicore.NewVcn(ctx, name, &ocicore.VcnArgs{
		CompartmentId: args.CompartmentID,
		CidrBlocks: pulumi.StringArray{
			pulumi.String(args.CidrBlock),
		},
		DnsLabel:      args.DNSLabel,
		IsIpv6enabled: args.Ipv6Enabled,
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("error creating vcn: %v", err)
	}

	// set a default for number of subnets
	var subnetNumbers int
	if args.NumberOfSubnets == 0 {
		subnetNumbers = 1
	} else {
		subnetNumbers = args.NumberOfSubnets
	}

	// populate the subnets from the CidrBlock we set
	privateSubnets, publicSubnets, err := SubnetDistributor(args.CidrBlock, subnetNumbers)
	if err != nil {
		return nil, err
	}

	// create the private subnets
	var privateSubnetIds pulumi.StringArray
	var privateRouteTableRules ocicore.RouteTableRouteRuleArray

	for index, subnet := range privateSubnets {
		privateSubnet, err := ocicore.NewSubnet(ctx, fmt.Sprintf("%s-private-%d", name, index+1), &ocicore.SubnetArgs{
			VcnId:                  vcn.ID(),
			CompartmentId:          args.CompartmentID,
			CidrBlock:              pulumi.String(subnet),
			ProhibitPublicIpOnVnic: pulumi.Bool(true),
		}, pulumi.Parent(vcn))
		if err != nil {
			return nil, fmt.Errorf("error creating subnet: %v", err)
		}
		privateSubnetIds = append(privateSubnetIds, privateSubnet.ID().ToStringOutput())
	}

	var publicSubnetIds pulumi.StringArray
	var publicRouteTableRules ocicore.RouteTableRouteRuleArray
	for index, subnet := range publicSubnets {
		publicSubnet, err := ocicore.NewSubnet(ctx, fmt.Sprintf("%s-public-%d", name, index+1), &ocicore.SubnetArgs{
			VcnId:                  vcn.ID(),
			CompartmentId:          args.CompartmentID,
			CidrBlock:              pulumi.String(subnet),
			ProhibitPublicIpOnVnic: pulumi.Bool(false),
		}, pulumi.Parent(vcn))
		if err != nil {
			return nil, fmt.Errorf("error creating subnet: %v", err)
		}
		publicSubnetIds = append(publicSubnetIds, publicSubnet.ID().ToStringOutput())
	}

	if args.CreateInternetGateway {
		igw, err := ocicore.NewInternetGateway(ctx, name, &ocicore.InternetGatewayArgs{
			CompartmentId: args.CompartmentID,
			VcnId:         vcn.ID(),
		}, pulumi.Parent(vcn))
		if err != nil {
			return nil, fmt.Errorf("error creating internet gateway: %v", err)
		}
		component.InternetGatewayID = igw.ID()
		publicRouteTableRules = append(publicRouteTableRules, ocicore.RouteTableRouteRuleArgs{
			Destination:     pulumi.String("0.0.0.0/0"),
			Description:     pulumi.String("traffic to/from internet"),
			NetworkEntityId: igw.ID(),
		})
	}

	if args.CreateNATGateway {
		natGateway, err := ocicore.NewNatGateway(ctx, name, &ocicore.NatGatewayArgs{
			CompartmentId: args.CompartmentID,
			VcnId:         vcn.ID(),
		}, pulumi.Parent(vcn))
		if err != nil {
			return nil, fmt.Errorf("error creating NAT gateway: %v", err)
		}
		component.NATGatewayID = natGateway.ID()
		privateRouteTableRules = append(privateRouteTableRules, ocicore.RouteTableRouteRuleArgs{
			Destination:     pulumi.String("0.0.0.0/0"),
			Description:     pulumi.String("traffic to the internet"),
			NetworkEntityId: natGateway.ID(),
		})
	}

	if args.CreateServiceGateway {
		trueBool := true
		svcs, err := ocicore.GetServices(ctx, &ocicore.GetServicesArgs{
			Filters: []ocicore.GetServicesFilter{
				{
					Name:   "name",
					Values: []string{"All .* Services In Oracle Services Network"},
					Regex:  &trueBool,
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("error retrieveing services for service gateway: %v", err)
		}

		svcGateway, err := ocicore.NewServiceGateway(ctx, name, &ocicore.ServiceGatewayArgs{
			CompartmentId: args.CompartmentID,
			VcnId:         vcn.ID(),
			Services: ocicore.ServiceGatewayServiceArray{
				ocicore.ServiceGatewayServiceArgs{
					ServiceId: pulumi.String(svcs.Services[0].Id),
				},
			},
		}, pulumi.Parent(vcn))
		if err != nil {
			return nil, fmt.Errorf("error creating service gateway: %v", err)
		}
		component.ServiceGatewayID = svcGateway.ID()

		privateRouteTableRules = append(privateRouteTableRules, ocicore.RouteTableRouteRuleArgs{
			DestinationType: pulumi.String("SERVICE_CIDR_BLOCK"),
			Description:     pulumi.String("traffic to OCI services"),
			NetworkEntityId: svcGateway.ID(),
			Destination:     pulumi.String(svcs.Services[0].CidrBlock),
		})
	}

	// handle the route tables
	privateRouteTable, err := ocicore.NewRouteTable(ctx, fmt.Sprintf("%s-private", name), &ocicore.RouteTableArgs{
		CompartmentId: args.CompartmentID,
		VcnId:         vcn.ID(),
		RouteRules:    privateRouteTableRules,
	}, pulumi.Parent(vcn))
	if err != nil {
		return nil, fmt.Errorf("error creating private route table: %v", err)
	}

	publicRouteTable, err := ocicore.NewRouteTable(ctx, fmt.Sprintf("%s-public", name), &ocicore.RouteTableArgs{
		CompartmentId: args.CompartmentID,
		VcnId:         vcn.ID(),
		RouteRules:    publicRouteTableRules,
	}, pulumi.Parent(vcn))
	if err != nil {
		return nil, fmt.Errorf("error creating public route table: %v", err)
	}

	for index, subnet := range publicSubnetIds {
		_, err = ocicore.NewRouteTableAttachment(ctx, fmt.Sprintf("%s-public-%d", name, index), &ocicore.RouteTableAttachmentArgs{
			RouteTableId: publicRouteTable.ID(),
			SubnetId:     subnet,
		}, pulumi.Parent(publicRouteTable))
		if err != nil {
			return nil, fmt.Errorf("error creating route table attachment for subnet %v: %v", subnet, err)
		}
	}

	for index, subnet := range privateSubnetIds {
		_, err = ocicore.NewRouteTableAttachment(ctx, fmt.Sprintf("%s-private-%d", name, index), &ocicore.RouteTableAttachmentArgs{
			RouteTableId: privateRouteTable.ID(),
			SubnetId:     subnet,
		}, pulumi.Parent(privateRouteTable))
		if err != nil {
			return nil, fmt.Errorf("error creating route table attachment for subnet %v: %v", subnet, err)
		}
	}

	component.VcnID = vcn.ID()
	component.PrivateSubnetIDs = privateSubnetIds.ToStringArrayOutput()
	component.PublicSubnetIDs = publicSubnetIds.ToStringArrayOutput()

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{
		"vcnId":            vcn.ID(),
		"privateSubnetIds": privateSubnetIds.ToStringArrayOutput(),
		"publicSubnetIds":  publicSubnetIds.ToStringArrayOutput(),
	}); err != nil {
		return nil, err
	}

	return component, nil
}
