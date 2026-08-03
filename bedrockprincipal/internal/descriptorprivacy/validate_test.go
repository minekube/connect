package descriptorprivacy

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestV2ProposalAdditionsExposeOnlyOpaqueEnvelopeAndBindings(t *testing.T) {
	descriptor := approvedDescriptor()
	require.NoError(t, ValidateV2ProposalAdditions(descriptor))
	descriptor.Field = append(descriptor.Field, &descriptorpb.FieldDescriptorProto{Name: proto.String("canonical_xuid"), Number: proto.Int32(13), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum()})
	require.Error(t, ValidateV2ProposalAdditions(descriptor))
}

func approvedDescriptor() *descriptorpb.DescriptorProto {
	return &descriptorpb.DescriptorProto{Name: proto.String("Session"), Field: []*descriptorpb.FieldDescriptorProto{
		{Name: proto.String("endpoint_id"), Number: proto.Int32(7), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum()},
		{Name: proto.String("organization_id"), Number: proto.Int32(8), Type: descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum()},
		{Name: proto.String("connect_session_nonce"), Number: proto.Int32(9), Type: descriptorpb.FieldDescriptorProto_TYPE_BYTES.Enum()},
		{Name: proto.String("source_protocol_version"), Number: proto.Int32(10), Type: descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum()},
		{Name: proto.String("policy_revision"), Number: proto.Int32(11), Type: descriptorpb.FieldDescriptorProto_TYPE_INT64.Enum()},
		{Name: proto.String("signed_bedrock_principal_v2"), Number: proto.Int32(12), Type: descriptorpb.FieldDescriptorProto_TYPE_BYTES.Enum()},
	}}
}
