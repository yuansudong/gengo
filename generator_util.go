package gengo

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	plugin "github.com/yuansudong/gengo/plugin"
	"google.golang.org/protobuf/proto"
)

// NameUpper 名字命名为大写
func NameUpper(str string) string {
	arr := strings.Split(str, "_")
	for index, val := range arr {
		arr[index] = strings.ToUpper(val[:1]) + val[1:]
	}
	return strings.Join(arr, "")
}

// GetRequest 用于从一个标准输入中,获取一个解析请求
func GetRequest(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	input, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("读取标准输入失败: %v", err)
	}
	req := new(plugin.CodeGeneratorRequest)
	if err = proto.Unmarshal(input, req); err != nil {
		return nil, fmt.Errorf("读取标准输入失败: %v", err)
	}
	return req, nil
}

// WriteFiles 用于写入一批响应到标准输出
func WriteFiles(files []*plugin.CodeGeneratorResponse_File) {
	WriteResponse(&plugin.CodeGeneratorResponse{File: files})
}

// WriteError 用于写一个错误到标准输出
func WriteError(err error) {
	WriteResponse(&plugin.CodeGeneratorResponse{Error: proto.String(err.Error())})
}

// WriteResponse 用于向标准输出中写数据
func WriteResponse(rsp *plugin.CodeGeneratorResponse) {
	buf, err := proto.Marshal(rsp)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		log.Fatalln(err.Error())
	}
}

// ReplaceArgs 用于替代参数
func ReplaceArgs(str string, args map[string]interface{}) string {

	for key, val := range args {
		var tmpV string
		switch val.(type) {
		default:
			tmpV = fmt.Sprint(val)
		}
		str = strings.Replace(str, key, tmpV, -1)
	}
	return str
}

// FieldDescriptorProto_TYPE_DOUBLE FieldDescriptorProto_Type = 1
//	FieldDescriptorProto_TYPE_FLOAT  FieldDescriptorProto_Type = 2

/*
	// FieldDescriptorProto_TYPE_DOUBLE FieldDescriptorProto_Type = 1
//	FieldDescriptorProto_TYPE_FLOAT  FieldDescriptorProto_Type = 2
	FieldDescriptorProto_TYPE_INT64  FieldDescriptorProto_Type = 3
	FieldDescriptorProto_TYPE_UINT64 FieldDescriptorProto_Type = 4
	FieldDescriptorProto_TYPE_INT32   FieldDescriptorProto_Type = 5
	FieldDescriptorProto_TYPE_FIXED64 FieldDescriptorProto_Type = 6
	FieldDescriptorProto_TYPE_FIXED32 FieldDescriptorProto_Type = 7
	FieldDescriptorProto_TYPE_BOOL    FieldDescriptorProto_Type = 8
	FieldDescriptorProto_TYPE_STRING  FieldDescriptorProto_Type = 9
	FieldDescriptorProto_TYPE_GROUP   FieldDescriptorProto_Type = 10
	FieldDescriptorProto_TYPE_MESSAGE FieldDescriptorProto_Type = 11 // Length-delimited aggregate.
	// New in version 2.
	FieldDescriptorProto_TYPE_BYTES    FieldDescriptorProto_Type = 12
	FieldDescriptorProto_TYPE_UINT32   FieldDescriptorProto_Type = 13
	FieldDescriptorProto_TYPE_ENUM     FieldDescriptorProto_Type = 14
	FieldDescriptorProto_TYPE_SFIXED32 FieldDescriptorProto_Type = 15
	FieldDescriptorProto_TYPE_SFIXED64 FieldDescriptorProto_Type = 16
	FieldDescriptorProto_TYPE_SINT32   FieldDescriptorProto_Type = 17 // Uses ZigZag encoding.
	FieldDescriptorProto_TYPE_SINT64   FieldDescriptorProto_Type = 18 // Uses ZigZag encoding.

*/
