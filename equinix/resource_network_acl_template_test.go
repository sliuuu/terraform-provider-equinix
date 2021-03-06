package equinix

import (
	"testing"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestNetworkACLTemplate_createFromResourceData(t *testing.T) {
	expected := ne.ACLTemplate{
		Name:        ne.String("test"),
		Description: ne.String("testTemplate"),
		MetroCode:   ne.String("SV"),
		InboundRules: []ne.ACLTemplateInboundRule{
			{
				SeqNo:    ne.Int(1),
				SrcType:  ne.String("SUBNET"),
				Subnets:  []string{"10.0.0.0/24", "1.1.1.1/32"},
				Protocol: ne.String("TCP"),
				SrcPort:  ne.String("any"),
				DstPort:  ne.String("8080"),
			},
		},
	}
	rawData := map[string]interface{}{
		networkACLTemplateSchemaNames["Name"]:        ne.StringValue(expected.Name),
		networkACLTemplateSchemaNames["Description"]: ne.StringValue(expected.Description),
		networkACLTemplateSchemaNames["MetroCode"]:   ne.StringValue(expected.MetroCode),
	}
	d := schema.TestResourceDataRaw(t, createNetworkACLTemplateSchema(), rawData)
	d.Set(networkACLTemplateSchemaNames["InboundRules"], flattenACLTemplateInboundRules(expected.InboundRules))
	//when
	result := createACLTemplate(d)
	//then
	assert.Equal(t, expected, result, "Created ACL Template matches expected result")
}

func TestNetworkACLTemplate_updateResourceData(t *testing.T) {
	input := &ne.ACLTemplate{
		Name:        ne.String("test"),
		Description: ne.String("testTemplate"),
		MetroCode:   ne.String("SV"),
		InboundRules: []ne.ACLTemplateInboundRule{
			{
				SeqNo:    ne.Int(1),
				SrcType:  ne.String("SUBNET"),
				Subnets:  []string{"10.0.0.0/24", "1.1.1.1/32"},
				Protocol: ne.String("TCP"),
				SrcPort:  ne.String("any"),
				DstPort:  ne.String("8080"),
			},
		},
	}
	d := schema.TestResourceDataRaw(t, createNetworkACLTemplateSchema(), make(map[string]interface{}))
	//when
	err := updateACLTemplateResource(input, d)
	//then
	assert.Nil(t, err, "Update of resource data does not return error")
	assert.Equal(t, ne.StringValue(input.Name), d.Get(networkACLTemplateSchemaNames["Name"]), "Name matches")
	assert.Equal(t, ne.StringValue(input.Description), d.Get(networkACLTemplateSchemaNames["Description"]), "Description matches")
	assert.Equal(t, ne.StringValue(input.MetroCode), d.Get(networkACLTemplateSchemaNames["MetroCode"]), "MetroCode matches")
	assert.Equal(t, input.InboundRules, expandACLTemplateInboundRules(d.Get(networkACLTemplateSchemaNames["InboundRules"]).([]interface{})), "InboundRules matches")
}

func TestNetworkACLTemplate_expandInboundRules(t *testing.T) {
	//given
	input := []interface{}{
		map[string]interface{}{
			networkACLTemplateInboundRuleSchemaNames["Subnets"]:  []interface{}{"10.0.0.0/24", "1.1.1.1/32"},
			networkACLTemplateInboundRuleSchemaNames["Protocol"]: "TCP",
			networkACLTemplateInboundRuleSchemaNames["SrcPort"]:  "any",
			networkACLTemplateInboundRuleSchemaNames["DstPort"]:  "8080",
		},
		map[string]interface{}{
			networkACLTemplateInboundRuleSchemaNames["Subnets"]:  []interface{}{"3.3.3.3/32"},
			networkACLTemplateInboundRuleSchemaNames["Protocol"]: "IP",
			networkACLTemplateInboundRuleSchemaNames["SrcPort"]:  "any",
			networkACLTemplateInboundRuleSchemaNames["DstPort"]:  "any",
		},
	}
	expected := []ne.ACLTemplateInboundRule{
		{
			SeqNo:    ne.Int(1),
			SrcType:  ne.String("SUBNET"),
			Subnets:  expandListToStringList(input[0].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["Subnets"]].([]interface{})),
			Protocol: ne.String(input[0].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["Protocol"]].(string)),
			SrcPort:  ne.String(input[0].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["SrcPort"]].(string)),
			DstPort:  ne.String(input[0].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["DstPort"]].(string)),
		},
		{
			SeqNo:    ne.Int(2),
			SrcType:  ne.String("SUBNET"),
			Subnets:  expandListToStringList(input[1].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["Subnets"]].([]interface{})),
			Protocol: ne.String(input[1].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["Protocol"]].(string)),
			SrcPort:  ne.String(input[1].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["SrcPort"]].(string)),
			DstPort:  ne.String(input[1].(map[string]interface{})[networkACLTemplateInboundRuleSchemaNames["DstPort"]].(string)),
		},
	}
	//when
	result := expandACLTemplateInboundRules(input)
	//then
	assert.Equal(t, expected, result, "Expanded ACL template inbound rules matches expected result")
}

func TestNetworkACLTemplate_flattenInboundRules(t *testing.T) {
	input := []ne.ACLTemplateInboundRule{
		{
			SeqNo:    ne.Int(1),
			SrcType:  ne.String("SUBNET"),
			Subnets:  []string{"10.0.0.0/24", "1.1.1.1/32"},
			Protocol: ne.String("TCP"),
			SrcPort:  ne.String("any"),
			DstPort:  ne.String("8080"),
		},
		{
			SeqNo:    ne.Int(2),
			SrcType:  ne.String("SUBNET"),
			Subnets:  []string{"3.3.3.3/32"},
			Protocol: ne.String("IP"),
			SrcPort:  ne.String("any"),
			DstPort:  ne.String("any"),
		},
	}
	expected := []interface{}{
		map[string]interface{}{
			networkACLTemplateInboundRuleSchemaNames["SeqNo"]:    input[0].SeqNo,
			networkACLTemplateInboundRuleSchemaNames["SrcType"]:  input[0].SrcType,
			networkACLTemplateInboundRuleSchemaNames["Subnets"]:  input[0].Subnets,
			networkACLTemplateInboundRuleSchemaNames["Protocol"]: input[0].Protocol,
			networkACLTemplateInboundRuleSchemaNames["SrcPort"]:  input[0].SrcPort,
			networkACLTemplateInboundRuleSchemaNames["DstPort"]:  input[0].DstPort,
		},
		map[string]interface{}{
			networkACLTemplateInboundRuleSchemaNames["SeqNo"]:    input[1].SeqNo,
			networkACLTemplateInboundRuleSchemaNames["SrcType"]:  input[1].SrcType,
			networkACLTemplateInboundRuleSchemaNames["Subnets"]:  input[1].Subnets,
			networkACLTemplateInboundRuleSchemaNames["Protocol"]: input[1].Protocol,
			networkACLTemplateInboundRuleSchemaNames["SrcPort"]:  input[1].SrcPort,
			networkACLTemplateInboundRuleSchemaNames["DstPort"]:  input[1].DstPort,
		},
	}
	//when
	result := flattenACLTemplateInboundRules(input)
	//then
	assert.Equal(t, expected, result, "Flattened ACL template inbound rules match expected result")
}
