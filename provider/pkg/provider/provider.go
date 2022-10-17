package provider

import (
	"github.com/pkg/errors"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/provider"
)

func construct(ctx *pulumi.Context, typ, name string, inputs provider.ConstructInputs,
	options pulumi.ResourceOption) (*provider.ConstructResult, error) {
	// TODO: Add support for additional component resources here.
	switch typ {
	case "oci-vcn:index:Vcn":
		return constructVcn(ctx, name, inputs, options)
	default:
		return nil, errors.Errorf("unknown resource type %s", typ)
	}
}

// constructVcn is an implementation of Construct for the example Vcn component.
// It demonstrates converting the raw ConstructInputs to the component's args struct, creating
// the component, and returning its URN and state (outputs).
func constructVcn(ctx *pulumi.Context, name string, inputs provider.ConstructInputs,
	options pulumi.ResourceOption) (*provider.ConstructResult, error) {

	// Copy the raw inputs to VcnArgs. `inputs.CopyTo` uses the types and `pulumi:` tags
	// on the struct's fields to convert the raw values to the appropriate Input types.
	args := &VcnArgs{}
	if err := inputs.CopyTo(args); err != nil {
		return nil, errors.Wrap(err, "setting args")
	}

	// Create the component resource.
	vcn, err := NewVcn(ctx, name, args, options)
	if err != nil {
		return nil, errors.Wrap(err, "creating component")
	}

	// Return the component resource's URN and state. `NewConstructResult` automatically sets the
	// ConstructResult's state based on resource struct fields tagged with `pulumi:` tags with a value
	// that is convertible to `pulumi.Input`.
	return provider.NewConstructResult(vcn)
}
