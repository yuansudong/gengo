package gengo

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	descriptor "github.com/yuansudong/gengo/descriptor"
	plugin "github.com/yuansudong/gengo/plugin"
)

// Registry 是从请求中的信息提取
type Registry struct {
	// _Msgs 是所有message的集合
	_Msgs map[string]*Message

	// Enums 是所有的枚举聚合
	_Enums map[string]*Enum

	// _Files 是所有的文件集合
	_Files map[string]*File

	// _prefix 描述一个golang包名的前缀
	_Prefix string

	// _ImportPath 导入路径
	_ImportPath string

	// _PkgMap 一个文件路径到一个protobuff的映射
	_PkgMap map[string]string

	// _PkgAliases 包的别名集合
	_PkgAliases map[string]string
}

// NewRegistry 实例化一个
func NewRegistry() *Registry {
	return &Registry{
		_Msgs:       make(map[string]*Message),
		_Enums:      make(map[string]*Enum),
		_Files:      make(map[string]*File),
		_PkgMap:     make(map[string]string),
		_PkgAliases: make(map[string]string),
	}
}

// Load 加载protobuf的services, methods, messages, enumerations.
func (r *Registry) Load(req *plugin.CodeGeneratorRequest) error {
	for _, file := range req.GetProtoFile() {
		r._LoadFile(file)
	}

	var sTargetPkg string
	for _, name := range req.FileToGenerate {
		target := r._Files[name]
		if target == nil {
			return fmt.Errorf("no such file: %s", name)
		}
		name := r._PackageIdentityName(target.FileDescriptorProto)
		if sTargetPkg == "" {
			sTargetPkg = name
		} else {
			if sTargetPkg != name {
				return fmt.Errorf("inconsistent package names: %s %s", sTargetPkg, name)
			}
		}
		if err := r._LoadServices(target); err != nil {
			return err
		}
	}
	return nil
}

// _LoadFile 用于加载文件
func (r *Registry) _LoadFile(file *descriptor.FileDescriptorProto) {
	pkg := GoPackage{
		Path: r._GoPackagePath(file),
		Name: r._DefaultGoPackageName(file),
	}
	if err := r.ReserveGoPackageAlias(pkg.Name, pkg.Path); err != nil {
		for i := 0; ; i++ {
			alias := fmt.Sprintf("%s_%d", pkg.Name, i)
			if err := r.ReserveGoPackageAlias(alias, pkg.Path); err == nil {
				pkg.Alias = alias
				break
			}
		}
	}
	f := &File{
		FileDescriptorProto: file,
		GoPkg:               pkg,
	}

	r._Files[file.GetName()] = f
	r._RegisterMsg(f, nil, file.GetMessageType())
	r._RegisterEnum(f, nil, file.GetEnumType())
}

// _RegisterMsg 注册message类型
func (r *Registry) _RegisterMsg(file *File, outerPath []string, msgs []*descriptor.DescriptorProto) {
	for i, md := range msgs {
		m := &Message{
			File:            file,
			Outers:          outerPath,
			DescriptorProto: md,
			Index:           i,
		}
		for _, fd := range md.GetField() {
			m.Fields = append(m.Fields, &Field{
				Message:              m,
				FieldDescriptorProto: fd,
			})
		}
		file.Messages = append(file.Messages, m)
		r._Msgs[m.FQMN()] = m
		var outers []string
		outers = append(outers, outerPath...)
		outers = append(outers, m.GetName())
		r._RegisterMsg(file, outers, m.GetNestedType())
		r._RegisterEnum(file, outers, m.GetEnumType())
	}
}

// _RegisterEnum 增加枚举类型
func (r *Registry) _RegisterEnum(file *File, outerPath []string, enums []*descriptor.EnumDescriptorProto) {
	for i, ed := range enums {
		e := &Enum{
			File:                file,
			Outers:              outerPath,
			EnumDescriptorProto: ed,
			Index:               i,
		}
		file.Enums = append(file.Enums, e)
		r._Enums[e.FQEN()] = e
	}
}

// LookupMsg 查找Message
func (r *Registry) LookupMsg(location, name string) (*Message, error) {
	if strings.HasPrefix(name, ".") {
		m, ok := r._Msgs[name]
		if !ok {
			return nil, fmt.Errorf("no message found: %s", name)
		}
		return m, nil
	}

	if !strings.HasPrefix(location, ".") {
		location = fmt.Sprintf(".%s", location)
	}
	components := strings.Split(location, ".")
	for len(components) > 0 {
		fqmn := strings.Join(append(components, name), ".")
		if m, ok := r._Msgs[fqmn]; ok {
			return m, nil
		}
		components = components[:len(components)-1]
	}
	return nil, fmt.Errorf("no message found: %s", name)
}

// LookupEnum 根据名字查找枚举
func (r *Registry) LookupEnum(location, name string) (*Enum, error) {
	if strings.HasPrefix(name, ".") {
		e, ok := r._Enums[name]
		if !ok {
			return nil, fmt.Errorf("no enum found: %s", name)
		}
		return e, nil
	}

	if !strings.HasPrefix(location, ".") {
		location = fmt.Sprintf(".%s", location)
	}
	components := strings.Split(location, ".")
	for len(components) > 0 {
		fqen := strings.Join(append(components, name), ".")
		if e, ok := r._Enums[fqen]; ok {
			return e, nil
		}
		components = components[:len(components)-1]
	}
	return nil, fmt.Errorf("no enum found: %s", name)
}

// LookupFile 通过名字查找文件.
func (r *Registry) LookupFile(name string) (*File, error) {
	f, ok := r._Files[name]
	if !ok {
		return nil, fmt.Errorf("没有找到这个File: %s", name)
	}
	return f, nil
}

// AddPkgMap 增加包名到protobuff包名的映射
func (r *Registry) AddPkgMap(file, protoPkg string) {
	r._PkgMap[file] = protoPkg
}

// SetPrefix 设置前缀
func (r *Registry) SetPrefix(prefix string) {
	r._Prefix = prefix
}

// SetImportPath 设置导入路径
func (r *Registry) SetImportPath(importPath string) {
	r._ImportPath = importPath
}

// ReserveGoPackageAlias 用于检索,alias 对应的包路径是否存在
func (r *Registry) ReserveGoPackageAlias(alias, pkgpath string) error {
	if taken, ok := r._PkgAliases[alias]; ok {
		if taken == pkgpath {
			return nil
		}
		return fmt.Errorf("package name %s is already taken. Use another alias", alias)
	}
	r._PkgAliases[alias] = pkgpath
	return nil
}

// _GoPackagePath 返回go的包名路径
func (r *Registry) _GoPackagePath(f *descriptor.FileDescriptorProto) string {
	name := f.GetName()
	if pkg, ok := r._PkgMap[name]; ok {
		return path.Join(r._Prefix, pkg)
	}

	gopkg := f.Options.GetGoPackage()
	idx := strings.LastIndex(gopkg, "/")
	if idx >= 0 {
		if sc := strings.LastIndex(gopkg, ";"); sc > 0 {
			gopkg = gopkg[:sc+1-1]
		}
		return gopkg
	}

	return path.Join(r._Prefix, path.Dir(name))
}

// GetAllFQMNs 返回所有的消息类型
func (r *Registry) GetAllFQMNs() []string {
	var keys []string
	for k := range r._Msgs {
		keys = append(keys, k)
	}
	return keys
}

// GetAllFQENs 返回所有的枚举类型
func (r *Registry) GetAllFQENs() []string {
	var keys []string
	for k := range r._Enums {
		keys = append(keys, k)
	}
	return keys
}

// _SanitizePackageName 整理包名
func _SanitizePackageName(pkgName string) string {
	pkgName = strings.Replace(pkgName, ".", "_", -1)
	pkgName = strings.Replace(pkgName, "-", "_", -1)
	return pkgName
}

// _DefaultGoPackageName returns 用于返回默认的go包名
func (r *Registry) _DefaultGoPackageName(f *descriptor.FileDescriptorProto) string {
	name := r._PackageIdentityName(f)
	return _SanitizePackageName(name)
}

// _PackageIdentityName 用于返回包的标识.
func (r *Registry) _PackageIdentityName(f *descriptor.FileDescriptorProto) string {
	if f.Options != nil && f.Options.GoPackage != nil {
		gopkg := f.Options.GetGoPackage()
		idx := strings.LastIndex(gopkg, "/")
		if idx < 0 {
			gopkg = gopkg[idx+1:]
		}

		gopkg = gopkg[idx+1:]
		// package name is overrided with the string after the
		// ';' character
		sc := strings.IndexByte(gopkg, ';')
		if sc < 0 {
			return _SanitizePackageName(gopkg)

		}
		return _SanitizePackageName(gopkg[sc+1:])
	}
	if p := r._ImportPath; len(p) != 0 {
		if i := strings.LastIndex(p, "/"); i >= 0 {
			p = p[i+1:]
		}
		return p
	}
	if f.Package == nil {
		base := filepath.Base(f.GetName())
		ext := filepath.Ext(base)
		return strings.TrimSuffix(base, ext)
	}
	return f.GetPackage()
}

// loadServices 从文件中解析服务
func (r *Registry) _LoadServices(file *File) error {
	var svcs []*Service
	for _, sd := range file.GetService() {
		svc := &Service{
			File:                   file,
			ServiceDescriptorProto: sd,
		}
		for _, md := range sd.GetMethod() {
			meth, err := r._NewMethod(svc, md)
			if err != nil {
				return err
			}
			svc.Methods = append(svc.Methods, meth)
		}
		if len(svc.Methods) == 0 {
			continue
		}
		svcs = append(svcs, svc)
	}
	file.Services = svcs
	return nil
}

func (r *Registry) _NewMethod(svc *Service, md *descriptor.MethodDescriptorProto) (*Method, error) {
	requestType, err := r.LookupMsg(svc.File.GetPackage(), md.GetInputType())
	if err != nil {
		return nil, err
	}
	responseType, err := r.LookupMsg(svc.File.GetPackage(), md.GetOutputType())
	if err != nil {
		return nil, err
	}
	meth := &Method{
		Service:               svc,
		MethodDescriptorProto: md,
		RequestType:           requestType,
		ResponseType:          responseType,
	}

	return meth, nil
}

// _NewParam RPC的参数
func (r *Registry) _NewParam(meth *Method, path string) (Parameter, error) {
	msg := meth.RequestType
	fields, err := r._ResolveFieldPath(msg, path, true)
	if err != nil {
		return Parameter{}, err
	}
	l := len(fields)
	if l == 0 {
		return Parameter{}, fmt.Errorf("invalid field access list for %s", path)
	}
	target := fields[l-1].Target
	switch target.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE, descriptor.FieldDescriptorProto_TYPE_GROUP:
		if IsWellKnownType(*target.TypeName) {
		} else {
			return Parameter{}, fmt.Errorf("aggregate type %s in parameter of %s.%s: %s", target.Type, meth.Service.GetName(), meth.GetName(), path)
		}
	}
	return Parameter{
		FieldPath: FieldPath(fields),
		Method:    meth,
		Target:    fields[l-1].Target,
	}, nil
}

// _NewBody 新实例化一个body,属于RPC中为http提供时,所需要指定的body
func (r *Registry) _NewBody(meth *Method, path string) (*Body, error) {
	msg := meth.RequestType
	switch path {
	case "":
		return nil, nil
	case "*":
		return &Body{FieldPath: nil}, nil
	}
	fields, err := r._ResolveFieldPath(msg, path, false)
	if err != nil {
		return nil, err
	}
	return &Body{FieldPath: FieldPath(fields)}, nil
}

// _NewResponse RPC的响应
func (r *Registry) _NewResponse(meth *Method, path string) (*Body, error) {
	msg := meth.ResponseType
	switch path {
	case "", "*":
		return nil, nil
	}
	fields, err := r._ResolveFieldPath(msg, path, false)
	if err != nil {
		return nil, err
	}
	return &Body{FieldPath: FieldPath(fields)}, nil
}

// _ResolveFieldPath 解析Message的路径
func (r *Registry) _ResolveFieldPath(msg *Message, path string, isPathParam bool) ([]FieldPathComponent, error) {
	if path == "" {
		return nil, nil
	}

	root := msg
	var result []FieldPathComponent
	for i, c := range strings.Split(path, ".") {
		if i > 0 {
			f := result[i-1].Target
			switch f.GetType() {
			case descriptor.FieldDescriptorProto_TYPE_MESSAGE, descriptor.FieldDescriptorProto_TYPE_GROUP:
				var err error
				msg, err = r.LookupMsg(msg.FQMN(), f.GetTypeName())
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("not an aggregate type: %s in %s", f.GetName(), path)
			}
		}
		f := msg.LookupField(c)
		if f == nil {
			return nil, fmt.Errorf("no field %q found in %s", path, root.GetName())
		}
		result = append(result, FieldPathComponent{Name: c, Target: f})
	}
	return result, nil
}
