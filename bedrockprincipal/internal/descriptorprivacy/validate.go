package descriptorprivacy

import (
	"fmt"

	"google.golang.org/protobuf/types/descriptorpb"
)

type fieldContract struct {
	name      string
	label     descriptorpb.FieldDescriptorProto_Label
	typeValue descriptorpb.FieldDescriptorProto_Type
}

var v2ProposalAdditions = map[int32]fieldContract{
	7:  {name: "endpoint_id", label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, typeValue: descriptorpb.FieldDescriptorProto_TYPE_STRING},
	8:  {name: "organization_id", label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, typeValue: descriptorpb.FieldDescriptorProto_TYPE_STRING},
	9:  {name: "connect_session_nonce", label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, typeValue: descriptorpb.FieldDescriptorProto_TYPE_BYTES},
	10: {name: "source_protocol_version", label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, typeValue: descriptorpb.FieldDescriptorProto_TYPE_INT32},
	11: {name: "policy_revision", label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, typeValue: descriptorpb.FieldDescriptorProto_TYPE_INT64},
	12: {name: "signed_bedrock_principal_v2", label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL, typeValue: descriptorpb.FieldDescriptorProto_TYPE_BYTES},
}

// ValidateV2ProposalAdditions checks the six frozen additive fields exactly.
// Existing legacy fields 1-6 are outside this v2 structural boundary.
func ValidateV2ProposalAdditions(descriptor *descriptorpb.DescriptorProto) error {
	seen := map[int32]bool{}
	for _, field := range descriptor.GetField() {
		number := field.GetNumber()
		if number < 7 {
			continue
		}
		contract, ok := v2ProposalAdditions[number]
		if !ok {
			return fmt.Errorf("unexpected v2 proposal field %d %q", number, field.GetName())
		}
		if seen[number] || field.GetName() != contract.name || field.GetLabel() != contract.label || field.GetType() != contract.typeValue {
			return fmt.Errorf("v2 proposal field %d differs from frozen contract", number)
		}
		seen[number] = true
	}
	for number := range v2ProposalAdditions {
		if !seen[number] {
			return fmt.Errorf("missing v2 proposal field %d", number)
		}
	}
	return nil
}
