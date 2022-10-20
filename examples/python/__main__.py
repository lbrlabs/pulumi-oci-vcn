"""A Python Pulumi program"""

import pulumi
import lbrlabs_pulumi_oci_vcn as lbrlabs
import pulumi_oci as oci

compartment = oci.identity.get_compartment_output(
    id="ocid1.tenancy.oc1..aaaaaaaavh67gfytloujvijdue6pqwhfrlo2ms4jbkpozyez7sgdhf27e4yq",
)

vcn = lbrlabs.Vcn(
    "example",
    compartment_id=compartment.id,
    cidr_block="172.16.0.0/22",
    create_internet_gateway=True,
    create_nat_gateway=True,
    create_service_gateway=True,
    dns_label="python",
)

pulumi.export("vcn_id", vcn.vcn_id)

