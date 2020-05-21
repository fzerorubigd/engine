package main

import (
	"fmt"
	"reflect"
	"strings"

	extrapb "github.com/fzerorubigd/protobuf/extra"
	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gogo/protobuf/vanity"
	"github.com/jinzhu/inflection"
	"github.com/kr/pretty"
)

type plugin struct {
	*generator.Generator
	generator.PluginImports
	useGogoImport bool

	modelImport     generator.Single
	timeImport      generator.Single
	initImport      generator.Single
	uuidImport      generator.Single
	protoTimeImport generator.Single
	contextImport   generator.Single
}

type modelData struct {
	table    string
	schema   string
	model    string
	receiver string
	dbFields []string
	goFields []string

	types     map[string]string
	createdAt bool
	updatedAt bool
	idType    string
}

func newPlugin(useGogoImport bool) generator.Plugin {
	return &plugin{
		useGogoImport: useGogoImport,
	}
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

func getVariableName(s string) string {
	if str := strings.ToLower(strings.Trim(s, "\n\t ")); len(str) < 3 {
		return str
	}
	s = CamelToSnake(s)
	arr := strings.Split(strings.ToLower(s), "_")
	res := ""
	for _, i := range arr {
		i = strings.Trim(i, " \n\t\"")
		if i != "" {
			res += i[0:1]
		}
	}
	if res == "m" {
		return "mm"
	}
	return res
}

func (p *plugin) Name() string {
	return "model"
}

func (p *plugin) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *plugin) Generate(file *generator.FileDescriptor) {
	if !p.useGogoImport {
		vanity.TurnOffGogoImport(file.FileDescriptorProto)
	}
	p.PluginImports = generator.NewPluginImports(p.Generator)
	p.modelImport = p.NewImport("elbix.dev/engine/pkg/postgres/model")
	p.initImport = p.NewImport("elbix.dev/engine/pkg/initializer")
	p.timeImport = p.NewImport("time")
	p.uuidImport = p.NewImport("github.com/google/uuid")
	p.protoTimeImport = p.NewImport("github.com/fzerorubigd/protobuf/types")
	p.contextImport = p.NewImport("context")

	var models []modelData

	for _, msg := range file.Messages() {
		if !proto.GetBoolExtension(msg.Options, extrapb.E_IsModel, false) {
			continue
		}
		model := modelData{}
		model.model = msg.GetName()
		model.receiver = getVariableName(model.model)
		model.schema = getString(msg.Options, extrapb.E_SchemaName, getString(file.Options, extrapb.E_SchemaNameAll, "public"))
		model.table = getString(msg.Options, extrapb.E_TableName, strings.ToLower(inflection.Plural(msg.GetName())))

		model.dbFields = p.getDbFields(msg)
		model.goFields = p.getGoFields(msg)
		model.types = p.getExtraMap(msg)
		model.createdAt = p.hasTimeField(msg, "created_at")
		model.updatedAt = p.hasTimeField(msg, "updated_at")
		model.idType = p.hasIDField(msg)
		models = append(models, model)
	}
	if len(models) == 0 {
		return
	}

	p.P("const (")
	for i := range models {
		p.createConstants(models[i])
	}
	p.P(")")

	p.createManager()

	for i := range models {
		p.P("/*")
		p.P(pretty.Sprint(models[i]))
		p.P("*/")
		if models[i].idType != "" {
			p.createFunction(models[i])
			p.updateFunction(models[i])
			p.byPrimaryFunction(models[i])
		}
		p.scanModel(models[i])
		p.getFiledMethod(models[i])
	}
}

func (p *plugin) createConstants(model modelData) {
	p.In()
	p.P(model.model, "Schema = ", `"`, model.schema, `"`)
	p.P(model.model, "Table = ", `"`, model.table, `"`)
	//	UserTableFull = UserSchema + "." + UserTable
	p.P(model.model, "TableFull = ", model.model, `Schema + "." +`, model.model, "Table")
	p.Out()
}

func (p *plugin) createManager() {
	p.P("type Manager struct {")
	p.In()
	p.P(p.modelImport.Use(), ".Manager")
	p.Out()
	p.P("}")
	p.P()
	p.P("func NewManager() *Manager {")
	p.In()
	p.P("return &Manager{}")
	p.Out()
	p.P("}")
	p.P()
	p.P("func NewManagerFromTransaction(tx ", p.modelImport.Use(), ".DBX) (*Manager, error) {")
	p.In()
	p.P("m := &Manager{}")
	p.P("err := m.Hijack(tx)")
	p.P()
	p.P("if err != nil {")
	p.In()
	p.P("return nil, err")
	p.Out()
	p.P("}")
	p.P()
	p.P("return m, nil")
	p.Out()
	p.P("}")
}

func (p *plugin) createFunction(msg modelData) {
	p.P()
	p.P("func (m *Manager) Create", msg.model, "(ctx ", p.contextImport.Use(), ".Context, ", msg.receiver, " *", msg.model, ") error {")
	p.In()
	// p.P("var err error")
	if msg.updatedAt || msg.createdAt {
		p.P("now := ", p.timeImport.Use(), ".Now()")
		if msg.createdAt {
			p.P(msg.receiver, ".CreatedAt = ", p.protoTimeImport.Use(), ".TimestampProto(now)")
		}
		if msg.updatedAt {
			p.P(msg.receiver, ".UpdatedAt = ", p.protoTimeImport.Use(), ".TimestampProto(now)")
		}
	}
	p.initClosure(msg.receiver, "PreInsert")

	l := len(msg.dbFields)
	if msg.idType != "" {
		l--
	}
	args := make([]string, l)
	values := make([]string, l)
	flds := make([]string, l)
	i := 0
	for j, dbf := range msg.dbFields {
		if dbf == "id" {
			continue
		}
		flds[i] = dbf
		args[i] = fmt.Sprintf("$%d", i+1)
		v := msg.types[dbf]
		if v == "Timestamp" {
			p.P(msg.goFields[j], ",err := ", p.protoTimeImport.Use(), ".TimestampFromProto(", msg.receiver, ".", msg.goFields[j], ")")
			p.ifErr()
			values[i] = msg.goFields[j]
		} else {
			values[i] = fmt.Sprintf("%s.%s", msg.receiver, msg.goFields[j])
		}
		i++
	}
	p.P(
		"q := `INSERT INTO ", msg.schema, ".", msg.table, "(",
		strings.Join(flds, ", "), ") VALUES (",
		strings.Join(args, ", "), ") RETURNING id`",
	)
	p.P("row := m.GetDbMap().QueryRowxContext(ctx, q, ", strings.Join(values, ", "), ")")
	p.P("return row.Scan(&", msg.receiver, ".Id)")
	p.Out()
	p.P("}")
}

func (p *plugin) updateFunction(msg modelData) {
	p.P()
	p.P("func (m *Manager) Update", msg.model, "(ctx ", p.contextImport.Use(), ".Context, ", msg.receiver, " *", msg.model, ") error {")
	p.In()
	p.P("var err error")
	if msg.updatedAt {
		p.P("now := ", p.timeImport.Use(), ".Now()")
		p.P(msg.receiver, ".UpdatedAt = ", p.protoTimeImport.Use(), ".TimestampProto(now)")
	}
	p.initClosure(msg.receiver, "PreUpdate")
	l := len(msg.dbFields)
	if msg.idType != "" {
		l--
	}
	args := make([]string, l)
	values := make([]string, l)
	i := 0
	for j, dbf := range msg.dbFields {
		if dbf == "id" {
			continue
		}
		args[i] = fmt.Sprintf("%s = $%d", dbf, i+1)
		v := msg.types[dbf]
		if v == "Timestamp" {
			p.P(msg.goFields[j], ",err := ", p.protoTimeImport.Use(), ".TimestampFromProto(", msg.receiver, ".", msg.goFields[j], ")")
			p.ifErr()
			values[i] = msg.goFields[j]
		} else {
			values[i] = fmt.Sprintf("%s.%s", msg.receiver, msg.goFields[j])
		}
		i++
	}

	p.P(
		"q := `UPDATE ", msg.schema, ".", msg.table, " SET ",
		strings.Join(args, ", "), " WHERE id = $", i+1, "`",
	)
	p.P("_, err = m.GetDbMap().ExecContext(ctx, q, ", strings.Join(values, ", "), ", ", msg.receiver, ".Id )")
	p.P("return err")
	p.Out()
	p.P("}")
}

func (p *plugin) byPrimaryFunction(msg modelData) {
	p.P()
	p.P("func (m *Manager) Get", msg.model, "ByPrimary(ctx ", p.contextImport.Use(), ".Context, id ", msg.idType, ") (*", msg.model, ", error){")
	p.In()
	p.P("q := `SELECT ", strings.Join(msg.dbFields, ", "), " FROM ", msg.schema, ".", msg.table, " WHERE id = $1`")

	p.P("row := m.GetDbMap().QueryRowxContext(ctx, q, id)")
	p.P()
	p.P("return m.scan", msg.model, "(row)")
	p.Out()
	p.P("}")
}

func (p *plugin) scanModel(msg modelData) {
	p.P()
	p.P("func (m *Manager) scan", msg.model, "(row ", p.modelImport.Use(), ".Scanner, extra ...interface{}) (*", msg.model, ", error){")
	p.P("var ", msg.receiver, " ", msg.model)
	values := make([]string, len(msg.dbFields))
	for j, dbf := range msg.dbFields {
		v := msg.types[dbf]
		if v == "Timestamp" {
			p.P("var ", msg.goFields[j], " ", p.timeImport.Use(), ".Time")
			values[j] = fmt.Sprintf("&%s", msg.goFields[j])
		} else {
			values[j] = fmt.Sprintf("&%s.%s", msg.receiver, msg.goFields[j])
		}
	}
	p.P("all := append([]interface{}{", strings.Join(values, ", "), "}, extra...)")
	p.P("err := row.Scan(all...)")
	p.P("if err != nil {")
	p.In()
	p.P("return nil, err")
	p.Out()
	p.P("}")

	for j, dbf := range msg.dbFields {
		v := msg.types[dbf]
		if v == "Timestamp" {
			p.P(msg.receiver, ".", msg.goFields[j], ", _ = ", p.protoTimeImport.Use(), ".TimestampProto(", msg.goFields[j], ")")
		}
	}
	p.P("return &", msg.receiver, ", nil")
	p.Out()
	p.P("}")

}

func (p *plugin) getFiledMethod(msg modelData) {
	p.P()
	p.P("func (m *Manager) get", msg.model, "Fields() []string {")
	p.In()
	p.P("return ", pretty.Sprint(msg.dbFields))
	p.Out()
	p.P("}")
}

func (p *plugin) ifErr() {
	p.P("if err != nil {")
	p.In()
	p.P("return err")
	p.Out()
	p.P("}")
}

func (p *plugin) getDbFields(msg *generator.Descriptor) []string {
	var (
		db []string
	)
	for _, f := range msg.Field {
		tag := strings.Trim(getString(f, gogoproto.E_Moretags, `"db:"`+strings.ToLower(f.GetName())), `"`) + `"`

		db = append(db, reflect.StructTag(tag).Get("db"))
	}

	return db
}

func (p *plugin) getGoFields(msg *generator.Descriptor) []string {
	var (
		fl []string
	)
	for _, f := range msg.Field {
		fl = append(fl, SnakeToCamel(f.GetName()))
	}

	return fl
}

func (p *plugin) getExtraMap(msg *generator.Descriptor) map[string]string {
	var (
		tps = map[string]string{}
	)
	for _, f := range msg.Field {
		if f.GetTypeName() == ".google.protobuf.Timestamp" {
			tps[f.GetName()] = "Timestamp"
		}

	}
	return tps
}

func (p *plugin) hasTimeField(msg *generator.Descriptor, s string) bool {
	for _, f := range msg.Field {
		if f.GetTypeName() == ".types.Timestamp" {
			if f.GetName() == s {
				return true
			}
		}
	}

	return false
}

func (p *plugin) hasIDField(msg *generator.Descriptor) string {
	for _, f := range msg.Field {
		if f.GetName() == "id" {
			switch f.GetType() {
			case descriptor.FieldDescriptorProto_TYPE_INT64:
				return "int64"
			case descriptor.FieldDescriptorProto_TYPE_STRING:
				return "string"
			default:
				return ""
			}
		}
	}

	return ""
}

func (p *plugin) initClosure(rec, fn string) {
	p.P("func(in interface{}) {")
	p.In()
	p.P("if o, ok := in.(interface{ ", fn, "() }); ok {")
	p.In()
	p.P("o.", fn, "()")
	p.Out()
	p.P("}")
	p.Out()
	p.P("}(", rec, ")")
}
