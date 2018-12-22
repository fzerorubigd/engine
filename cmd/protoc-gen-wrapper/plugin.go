package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fzerorubigd/protobuf/extra"

	"github.com/fzerorubigd/balloon/pkg/assert"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gogo/protobuf/vanity"
)

type plugin struct {
	*generator.Generator
	generator.PluginImports
	useGogoImport bool

	grpcgwImport     generator.Single
	validatorImport  generator.Single
	contextImport    generator.Single
	inprocgrpcImport generator.Single
	runtimeImport    generator.Single
	assertImport     generator.Single
	logImport        generator.Single
	errorsImport     generator.Single
	resourcesImport  generator.Single
}

func newPlugin(useGogoImport bool) generator.Plugin {
	return &plugin{
		useGogoImport: useGogoImport,
	}
}

func (p *plugin) Name() string {
	return "wrapper"
}

func (p *plugin) Init(g *generator.Generator) {
	p.Generator = g
}
func getString(msg proto.Message, extension *proto.ExtensionDesc, def string) string {
	ss, err := proto.GetExtension(msg, extension)
	if err == nil {
		if str, ok := ss.(*string); ok {
			return *str
		}
	}

	return def
}

func (p *plugin) Generate(file *generator.FileDescriptor) {
	if !p.useGogoImport {
		vanity.TurnOffGogoImport(file.FileDescriptorProto)
	}
	p.PluginImports = generator.NewPluginImports(p.Generator)
	p.grpcgwImport = p.NewImport("github.com/fzerorubigd/balloon/pkg/grpcgw")
	p.validatorImport = p.NewImport("gopkg.in/go-playground/validator.v9")
	p.contextImport = p.NewImport("golang.org/x/net/context")
	p.inprocgrpcImport = p.NewImport("github.com/fullstorydev/grpchan/inprocgrpc")
	p.runtimeImport = p.NewImport("github.com/grpc-ecosystem/grpc-gateway/runtime")
	p.assertImport = p.NewImport("github.com/fzerorubigd/balloon/pkg/assert")
	p.logImport = p.NewImport("github.com/fzerorubigd/balloon/pkg/log")
	p.resourcesImport = p.NewImport("github.com/fzerorubigd/balloon/pkg/resources")
	p.errorsImport = p.NewImport("github.com/pkg/errors")

	resMap := make(map[string]string)
	var order []string
	for _, svc := range file.GetService() {
		p.createWrappedInterface(svc.GetName())
		p.P()
		p.createWrappedStruct(svc.GetName())
		p.P()
		p.createInitFunction(svc.GetName())
		//p.P("/*")
		for _, m := range svc.GetMethod() {
			p.P()
			p.createMethod(svc.GetName(), m.GetName(), m.GetInputType(), m.GetOutputType())
			res := getString(m.Options, extrapb.E_Resource, "__DEFAULT__")
			if res != "__DEFAULT__" {
				m := getFullName(file.GetName(), svc.GetName(), m.GetName())
				order = append(order, m)
				resMap[m] = res
			}
		}
		p.P()
		p.createNewFunction(svc.GetName())
		//p.P("*/")
	}

	p.P("func init() {")
	p.In()
	for _, i := range order {
		p.P(p.resourcesImport.Use(), ".RegisterResource(", fmt.Sprintf("%q, %q", i, resMap[i]), ")")
	}
	p.Out()
	p.P("}")
}

func getFullName(fl, svc, meth string) string {
	parts := strings.Split(fl, "/")
	f := parts[len(parts)-1]
	f = strings.TrimSuffix(f, filepath.Ext(f))

	return "/" + f + "." + svc + "/" + meth
}

func (p *plugin) createWrappedInterface(class string) {
	p.P("type Wrapped", class, "Controller interface {")
	p.In()
	p.P(class, "Server")
	p.P(p.grpcgwImport.Use(), ".Controller")
	p.Out()
	p.P("}")
}

func (p *plugin) createWrappedStruct(class string) {
	p.P("type wrapped", class, "Server struct {")
	p.In()
	p.P("original ", class, "Server")
	p.P("v *", p.validatorImport.Use(), ".Validate")
	p.Out()
	p.P("}")
}

func (p *plugin) createInitFunction(class string) {
	p.P("func (w *wrapped", class, "Server) Init(ctx ", p.contextImport.Use(), ".Context, ch *", p.inprocgrpcImport.Use(), ".Channel, mux *", p.runtimeImport.Use(), ".ServeMux) {")
	p.In()
	p.P("RegisterHandler", class, "(ch, w)")
	p.P("cl := New", class, "ChannelClient(ch)")
	p.P()
	p.P(p.assertImport.Use(), ".Nil(Register", class, "HandlerClient(ctx, mux, cl))")
	p.Out()
	p.P("}")
}

func (p *plugin) createMethod(class, method, in, out string) {
	inp := strings.Split(in, ".")
	assert.True(len(inp) == 3)
	assert.Empty(inp[0])

	oup := strings.Split(out, ".")
	assert.True(len(oup) == 3)
	assert.Empty(oup[0])

	p.P("func (w *wrapped", class, "Server) ", method, "(ctx ", p.contextImport.Use(), ".Context, req *", inp[2], ") (res *", oup[2], ", err error) {")
	p.In()
	p.P(p.logImport.Use(), `.Info("`, class, ".", method, ` request")`)
	p.P("defer func() {")
	p.In()
	p.P("e := recover()")
	p.P("if e == nil {")
	p.In()
	p.P("return")
	p.Out()
	p.P("}")
	p.P(p.logImport.Use(), `.Error("Recovering from panic", `, p.logImport.Use(), `.Any("panic", e))`)
	p.P("res, err = nil, ", p.errorsImport.Use(), `.New("internal server error")`)
	p.Out()
	p.P("}()")
	p.P("ctx, err = ", p.grpcgwImport.Use(), ".ExecuteMiddleware(ctx, w.original)")
	p.P("if err != nil {")
	p.In()
	p.P("return nil, err")
	p.Out()
	p.P("}")
	p.P("if err = w.v.Struct(req); err != nil {")
	p.In()
	p.P("return nil, ", p.grpcgwImport.Use(), `.NewBadRequest(err, "validation failed")`)
	p.Out()
	p.P("}")
	p.P()
	p.P("res, err = w.original.", method, "(ctx, req)")
	p.P("return")
	p.Out()
	p.P("}")
}

func (p *plugin) createNewFunction(class string) {

	p.P("func NewWrapped", class, "Server(server ", class, "Server) Wrapped", class, "Controller {")
	p.In()
	p.P("return &wrapped", class, "Server{")
	p.In()
	p.P("original: server,")
	p.P("v:", p.validatorImport.Use(), ".New(),")
	p.Out()
	p.P("}")
	p.Out()
	p.P("}")
}
