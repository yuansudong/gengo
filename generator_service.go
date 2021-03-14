package gengo

import (
	"fmt"
	"strings"

	"github.com/yuansudong/gengo/descriptor"
)

// IsWellKnownType 用于判断是否是wellknown类型
func IsWellKnownType(typeName string) bool {
	_, ok := _WellKnownTypeConv[typeName]
	return ok
}

// Service 对protobuf的service类型再封装
type Service struct {
	// File 该Service定义于哪个文件
	File *File
	*descriptor.ServiceDescriptorProto
	// Methods 该服务下有哪些RPC的方法
	Methods []*Method
}

// FQSN 返回service的完整文件名
func (s *Service) FQSN() string {
	components := []string{""}
	if s.File.Package != nil {
		components = append(components, s.File.GetPackage())
	}
	components = append(components, s.GetName())
	return strings.Join(components, ".")
}

// Method 封装protobuf中的RPC的方法
type Method struct {
	// Service 这个RPC方法属于哪个Service,方便查找
	Service *Service
	*descriptor.MethodDescriptorProto
	// RequestType RPC方法的请求类型
	RequestType *Message
	// ResponseType RPC方法的响应类型
	ResponseType *Message
}

// FQMN 返回RPC的方法名
func (m *Method) FQMN() string {
	components := []string{}
	components = append(components, m.Service.FQSN())
	components = append(components, m.GetName())
	return strings.Join(components, ".")
}

// Field 封装字段方法
type Field struct {
	// Message 这个字段属于哪个消息
	Message *Message
	// FieldMessage 字段的消息类型.
	FieldMessage *Message
	*descriptor.FieldDescriptorProto
}

// Parameter 在RPC里的参数
type Parameter struct {
	// FieldPath 与字段映射
	FieldPath
	// Target 与目标映射.
	Target *Field
	// Method 这个参数属于RPC中的哪个方法
	Method *Method
}

// ConvertFuncExpr 获得转换函数
func (p Parameter) ConvertFuncExpr() (string, error) {
	tbl := _Proto3ConvertFuncs
	if !p.IsProto2() && p.IsRepeated() {
		tbl = _Proto3RepeatedConvertFuncs
	} else if p.IsProto2() && !p.IsRepeated() {
		tbl = _Proto2ConvertFuncs
	} else if p.IsProto2() && p.IsRepeated() {
		tbl = _Proto2RepeatedConvertFuncs
	}
	typ := p.Target.GetType()
	conv, ok := tbl[typ]
	if !ok {
		conv, ok = _WellKnownTypeConv[p.Target.GetTypeName()]
	}
	if !ok {
		return "", fmt.Errorf("unsupported field type %s of parameter %s in %s.%s", typ, p.FieldPath, p.Method.Service.GetName(), p.Method.GetName())
	}
	return conv, nil
}

// IsEnum 判断参数是否是枚举
func (p Parameter) IsEnum() bool {
	return p.Target.GetType() == descriptor.FieldDescriptorProto_TYPE_ENUM
}

// IsRepeated 判断参数是否是数组
func (p Parameter) IsRepeated() bool {
	return p.Target.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED
}

// IsProto2 返回是不是protobuf2协议
func (p Parameter) IsProto2() bool {
	return p.Target.Message.File.proto2()
}

// Body 描述一个http请求或者响应的body
type Body struct {
	// FieldPath 字段映射
	FieldPath FieldPath
}

// AssignableExpr 返回代码生成的表达式
func (b *Body) AssignableExpr(msgExpr string) string {
	return b.FieldPath.AssignableExpr(msgExpr)
}

// FieldPath 描述一个请求消息的映射结构
type FieldPath []FieldPathComponent

// String 返回一段文本描述
func (p FieldPath) String() string {
	var components []string
	for _, c := range p {
		components = append(components, c.Name)
	}
	return strings.Join(components, ".")
}

// IsNestedProto3 是否嵌套proto3协议
func (p FieldPath) IsNestedProto3() bool {
	if len(p) > 1 && !p[0].Target.Message.File.proto2() {
		return true
	}
	return false
}

// AssignableExpr 代码生成片段
func (p FieldPath) AssignableExpr(msgExpr string) string {
	l := len(p)
	if l == 0 {
		return msgExpr
	}

	var preparations []string
	components := msgExpr
	for i, c := range p {
		if c.Target.OneofIndex != nil {
			index := c.Target.OneofIndex
			msg := c.Target.Message
			oneOfName := Camel(msg.GetOneofDecl()[*index].GetName())
			oneofFieldName := msg.GetName() + "_" + c.AssignableExpr()
			components = components + "." + oneOfName
			s := `if %s == nil {
				%s =&%s{}
			} else if _, ok := %s.(*%s); !ok {
				return nil, metadata, status.Errorf(codes.InvalidArgument, "expect type: *%s, but: %%t\n",%s)
			}`

			preparations = append(preparations, fmt.Sprintf(s, components, components, oneofFieldName, components, oneofFieldName, oneofFieldName, components))
			components = components + ".(*" + oneofFieldName + ")"
		}

		if i == l-1 {
			components = components + "." + c.AssignableExpr()
			continue
		}
		components = components + "." + c.ValueExpr()
	}

	preparations = append(preparations, components)
	return strings.Join(preparations, "\n")
}

// FieldPathComponent ...
type FieldPathComponent struct {
	// Name protobuf字段中的名字
	Name   string
	Target *Field
}

// AssignableExpr 返回一个不可忽视的表达式,没有v3与v2的差别.
func (c FieldPathComponent) AssignableExpr() string {
	return Camel(c.Name)
}

// ValueExpr 为一个字段返回一个表达式.
func (c FieldPathComponent) ValueExpr() string {
	if c.Target.Message.File.proto2() {
		return fmt.Sprintf("Get%s()", Camel(c.Name))
	}
	return Camel(c.Name)
}
