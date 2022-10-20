import * as pulumi from "@pulumi/pulumi";
import * as vcn from "@lbrlabs/pulumi-oci-vcn"
import * as oci from "@pulumi/oci"
import * as init from "@pulumi/cloudinit"

const compartment = oci.identity.getCompartmentOutput({
    id: "ocid1.tenancy.oc1..aaaaaaaavh67gfytloujvijdue6pqwhfrlo2ms4jbkpozyez7sgdhf27e4yq",
  })

const v = new vcn.Vcn("ts-example", {
    compartmentId: compartment.id,
    cidrBlock: "172.16.0.0/22",
    createInternetGateway: true,
    createNatGateway: true,
    createServiceGateway: true,
    dnsLabel: "lbriggs",
})

export const vcnId = v.vcnId
export const publicSubnetIds = v.publicSubnetIds
export const privateSubnetIds = v.privateSubnetIds
