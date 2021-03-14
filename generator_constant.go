package gengo

import descriptor "github.com/yuansudong/gengo/descriptor"

var (
	_Proto3ConvertFuncs = map[descriptor.FieldDescriptorProto_Type]string{
		descriptor.FieldDescriptorProto_TYPE_DOUBLE:  "runtime.Float64",
		descriptor.FieldDescriptorProto_TYPE_FLOAT:   "runtime.Float32",
		descriptor.FieldDescriptorProto_TYPE_INT64:   "runtime.Int64",
		descriptor.FieldDescriptorProto_TYPE_UINT64:  "runtime.Uint64",
		descriptor.FieldDescriptorProto_TYPE_INT32:   "runtime.Int32",
		descriptor.FieldDescriptorProto_TYPE_FIXED64: "runtime.Uint64",
		descriptor.FieldDescriptorProto_TYPE_FIXED32: "runtime.Uint32",
		descriptor.FieldDescriptorProto_TYPE_BOOL:    "runtime.Bool",
		descriptor.FieldDescriptorProto_TYPE_STRING:  "runtime.String",
		// FieldDescriptorProto_TYPE_GROUP
		// FieldDescriptorProto_TYPE_MESSAGE
		descriptor.FieldDescriptorProto_TYPE_BYTES:    "runtime.Bytes",
		descriptor.FieldDescriptorProto_TYPE_UINT32:   "runtime.Uint32",
		descriptor.FieldDescriptorProto_TYPE_ENUM:     "runtime.Enum",
		descriptor.FieldDescriptorProto_TYPE_SFIXED32: "runtime.Int32",
		descriptor.FieldDescriptorProto_TYPE_SFIXED64: "runtime.Int64",
		descriptor.FieldDescriptorProto_TYPE_SINT32:   "runtime.Int32",
		descriptor.FieldDescriptorProto_TYPE_SINT64:   "runtime.Int64",
	}

	_Proto3RepeatedConvertFuncs = map[descriptor.FieldDescriptorProto_Type]string{
		descriptor.FieldDescriptorProto_TYPE_DOUBLE:  "runtime.Float64Slice",
		descriptor.FieldDescriptorProto_TYPE_FLOAT:   "runtime.Float32Slice",
		descriptor.FieldDescriptorProto_TYPE_INT64:   "runtime.Int64Slice",
		descriptor.FieldDescriptorProto_TYPE_UINT64:  "runtime.Uint64Slice",
		descriptor.FieldDescriptorProto_TYPE_INT32:   "runtime.Int32Slice",
		descriptor.FieldDescriptorProto_TYPE_FIXED64: "runtime.Uint64Slice",
		descriptor.FieldDescriptorProto_TYPE_FIXED32: "runtime.Uint32Slice",
		descriptor.FieldDescriptorProto_TYPE_BOOL:    "runtime.BoolSlice",
		descriptor.FieldDescriptorProto_TYPE_STRING:  "runtime.StringSlice",
		// FieldDescriptorProto_TYPE_GROUP
		// FieldDescriptorProto_TYPE_MESSAGE
		descriptor.FieldDescriptorProto_TYPE_BYTES:    "runtime.BytesSlice",
		descriptor.FieldDescriptorProto_TYPE_UINT32:   "runtime.Uint32Slice",
		descriptor.FieldDescriptorProto_TYPE_ENUM:     "runtime.EnumSlice",
		descriptor.FieldDescriptorProto_TYPE_SFIXED32: "runtime.Int32Slice",
		descriptor.FieldDescriptorProto_TYPE_SFIXED64: "runtime.Int64Slice",
		descriptor.FieldDescriptorProto_TYPE_SINT32:   "runtime.Int32Slice",
		descriptor.FieldDescriptorProto_TYPE_SINT64:   "runtime.Int64Slice",
	}

	_Proto2ConvertFuncs = map[descriptor.FieldDescriptorProto_Type]string{
		descriptor.FieldDescriptorProto_TYPE_DOUBLE:  "runtime.Float64P",
		descriptor.FieldDescriptorProto_TYPE_FLOAT:   "runtime.Float32P",
		descriptor.FieldDescriptorProto_TYPE_INT64:   "runtime.Int64P",
		descriptor.FieldDescriptorProto_TYPE_UINT64:  "runtime.Uint64P",
		descriptor.FieldDescriptorProto_TYPE_INT32:   "runtime.Int32P",
		descriptor.FieldDescriptorProto_TYPE_FIXED64: "runtime.Uint64P",
		descriptor.FieldDescriptorProto_TYPE_FIXED32: "runtime.Uint32P",
		descriptor.FieldDescriptorProto_TYPE_BOOL:    "runtime.BoolP",
		descriptor.FieldDescriptorProto_TYPE_STRING:  "runtime.StringP",
		// FieldDescriptorProto_TYPE_GROUP
		// FieldDescriptorProto_TYPE_MESSAGE
		// FieldDescriptorProto_TYPE_BYTES
		// TODO(yugui) Handle bytes
		descriptor.FieldDescriptorProto_TYPE_UINT32:   "runtime.Uint32P",
		descriptor.FieldDescriptorProto_TYPE_ENUM:     "runtime.EnumP",
		descriptor.FieldDescriptorProto_TYPE_SFIXED32: "runtime.Int32P",
		descriptor.FieldDescriptorProto_TYPE_SFIXED64: "runtime.Int64P",
		descriptor.FieldDescriptorProto_TYPE_SINT32:   "runtime.Int32P",
		descriptor.FieldDescriptorProto_TYPE_SINT64:   "runtime.Int64P",
	}

	_Proto2RepeatedConvertFuncs = map[descriptor.FieldDescriptorProto_Type]string{
		descriptor.FieldDescriptorProto_TYPE_DOUBLE:  "runtime.Float64Slice",
		descriptor.FieldDescriptorProto_TYPE_FLOAT:   "runtime.Float32Slice",
		descriptor.FieldDescriptorProto_TYPE_INT64:   "runtime.Int64Slice",
		descriptor.FieldDescriptorProto_TYPE_UINT64:  "runtime.Uint64Slice",
		descriptor.FieldDescriptorProto_TYPE_INT32:   "runtime.Int32Slice",
		descriptor.FieldDescriptorProto_TYPE_FIXED64: "runtime.Uint64Slice",
		descriptor.FieldDescriptorProto_TYPE_FIXED32: "runtime.Uint32Slice",
		descriptor.FieldDescriptorProto_TYPE_BOOL:    "runtime.BoolSlice",
		descriptor.FieldDescriptorProto_TYPE_STRING:  "runtime.StringSlice",
		// FieldDescriptorProto_TYPE_GROUP
		// FieldDescriptorProto_TYPE_MESSAGE
		// FieldDescriptorProto_TYPE_BYTES
		// TODO(maros7) Handle bytes
		descriptor.FieldDescriptorProto_TYPE_UINT32:   "runtime.Uint32Slice",
		descriptor.FieldDescriptorProto_TYPE_ENUM:     "runtime.EnumSlice",
		descriptor.FieldDescriptorProto_TYPE_SFIXED32: "runtime.Int32Slice",
		descriptor.FieldDescriptorProto_TYPE_SFIXED64: "runtime.Int64Slice",
		descriptor.FieldDescriptorProto_TYPE_SINT32:   "runtime.Int32Slice",
		descriptor.FieldDescriptorProto_TYPE_SINT64:   "runtime.Int64Slice",
	}

	_WellKnownTypeConv = map[string]string{
		".google.protobuf.Timestamp":   "runtime.Timestamp",
		".google.protobuf.Duration":    "runtime.Duration",
		".google.protobuf.StringValue": "runtime.StringValue",
		".google.protobuf.FloatValue":  "runtime.FloatValue",
		".google.protobuf.DoubleValue": "runtime.DoubleValue",
		".google.protobuf.BoolValue":   "runtime.BoolValue",
		".google.protobuf.BytesValue":  "runtime.BytesValue",
		".google.protobuf.Int32Value":  "runtime.Int32Value",
		".google.protobuf.UInt32Value": "runtime.UInt32Value",
		".google.protobuf.Int64Value":  "runtime.Int64Value",
		".google.protobuf.UInt64Value": "runtime.UInt64Value",
	}
)
