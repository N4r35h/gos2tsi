package gos2tsi

import (
	"go/doc"
	"go/types"
	"reflect"
	"strings"

	"golang.org/x/tools/go/packages"
)

var GoTypeToTSType = map[string]string{
	"bool":      "boolean",
	"interface": "any",
	"any":       "any",
	"string":    "string",
	"int":       "number",
	"int8":      "number",
	"int16":     "number",
	"int32":     "number",
	"int64":     "number",
	"uint":      "number",
	"uint8":     "number",
	"uint16":    "number",
	"uint32":    "number",
	"uint64":    "number",
	"float32":   "number",
	"float64":   "number",
}

type ParsedField struct {
	Var       *types.Var
	Tag       string
	TSName    string
	TSType    string
	IsSlice   bool
	RefStruct ParsedStruct
}

type ParsedStruct struct {
	PackageName        string
	PackgePath         string
	ID                 string
	Name               string
	Required           bool
	Fields             []ParsedField
	IsSlice            bool
	GenericPopulations []ParsedField
}

type Converter struct {
	Indent               string
	Structs              map[string]ParsedStruct
	Docs                 map[string]string
	AlreadyParsedPackage map[string]bool
}

func New() *Converter {
	return &Converter{
		Structs:              map[string]ParsedStruct{},
		Docs:                 map[string]string{},
		AlreadyParsedPackage: map[string]bool{},
	}
}

func (c *Converter) ParseStruct(interf interface{}) ParsedStruct {
	reflectType := reflect.TypeOf(interf)
	IsSlice := false
	if reflectType.Kind() == reflect.Slice {
		IsSlice = true
		reflectType = reflectType.Elem()
	}
	pkgPath := reflectType.PkgPath()
	RequiredStruct := reflectType.Name()
	return c.ParseStructsInPackage(pkgPath, RequiredStruct, IsSlice)
}

func (c *Converter) ensureGenericPopulations(structName string) {
	genericPopulations := strings.Split(structName, "[")
	if len(genericPopulations) > 1 {
		for _, gp := range strings.Split(strings.ReplaceAll(genericPopulations[1], "]", ""), ",") {
			fullNameSegments := strings.Split(gp, ".")
			if len(fullNameSegments) > 1 {
				c.ParseStructsInPackage(fullNameSegments[0], fullNameSegments[1], false)
			}
		}
	}
}

func (c *Converter) ParseStructsInPackage(pkgPath, RequiredStruct string, IsSlice bool) ParsedStruct {
	RequestedStruct := ParsedStruct{}
	RequestedStruct.IsSlice = IsSlice
	c.ensureGenericPopulations(RequiredStruct)
	_, exists := c.AlreadyParsedPackage[pkgPath]
	if exists {
		rs := c.Structs[pkgPath+"."+removeGenericsPartFromStructName(RequiredStruct)]
		rs.Required = true
		rs.IsSlice = RequestedStruct.IsSlice
		c.Structs[pkgPath+"."+removeGenericsPartFromStructName(RequiredStruct)] = rs
		return rs
	}
	cfg := &packages.Config{
		Mode:  packages.NeedTypes | packages.NeedName | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedName | packages.NeedSyntax,
		Tests: false,
	}
	packages, _ := packages.Load(cfg, pkgPath)
	c.AlreadyParsedPackage[pkgPath] = true
	for _, pkg := range packages {
		docs, _ := doc.NewFromFiles(pkg.Fset, pkg.Syntax, "")
		for _, v := range docs.Types {
			c.Docs[v.Name] = v.Doc
		}
		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			switch item := obj.Type().Underlying().(type) {
			case *types.Struct:
				ObjectName := types.TypeString(obj.Type(), func(other *types.Package) string { return "" })
				st, _ := item.Underlying().(*types.Struct)
				parsedStruct := ParsedStruct{
					Required:           isStructWithPackagePathEqualToRequiredStructName(RequiredStruct, obj.Type().String()),
					PackageName:        pkg.Name,
					PackgePath:         pkgPath,
					ID:                 obj.Type().String(),
					Name:               ObjectName,
					GenericPopulations: c.getGenericPopulations(RequiredStruct),
				}
				for i := 0; i < st.NumFields(); i++ {
					pf := ParsedField{
						Var: st.Field(i),
						Tag: st.Tag(i),
					}
					fieldName := pf.Var.Name()
					typeName := pf.Var.Type().String()
					isTypeValid := false
					if convertedTypeName, ok := GoTypeToTSType[typeName]; ok {
						typeName = convertedTypeName
						isTypeValid = true
					}
					typeNameSegments := strings.Split(typeName, ".")
					if len(typeNameSegments) > 1 {
						typeName = typeNameSegments[len(typeNameSegments)-1]
					}
					if pf.Tag != "" {
						fieldTag := reflect.StructTag(pf.Tag)
						jsonTag := fieldTag.Get("json")
						if jsonTag != "" {
							fieldName = strings.Split(jsonTag, ",")[0]
						}
						tsTypeTag := fieldTag.Get("ts_type")
						if tsTypeTag != "" {
							typeName = tsTypeTag
							isTypeValid = true
						}
					}
					pf.IsSlice = reflect.TypeOf(st.Field(i).Type()).String() == "*types.Slice"
					pf.TSName = fieldName
					pf.TSType = strings.ReplaceAll(typeName, "[]", "")
					if pf.TSName != "-" {
						if !isTypeValid && len(typeNameSegments) > 1 && !pf.Var.Embedded() {
							FieldTypeName := pf.Var.Type().String()
							if pf.IsSlice {
								FieldTypeName = strings.Replace(FieldTypeName, "[]", "", 1)
							}
							DotSepFSPN := strings.Split(FieldTypeName, ".")
							StructName := DotSepFSPN[len(DotSepFSPN)-1]
							fieldPackagePath := strings.Replace(FieldTypeName, "."+StructName, "", 1)
							refStruct := c.ParseStructsInPackage(fieldPackagePath, StructName, pf.IsSlice)
							refStruct.Required = true
							refStruct.IsSlice = pf.IsSlice
							pf.RefStruct = refStruct
						}
						parsedStruct.Fields = append(parsedStruct.Fields, pf)
					}
				}
				if parsedStruct.Required {
					RequestedStruct = parsedStruct
				}
				/*
					a small workaround as we call this with in a package scan itself we could just ref this package
					which is marked as already scaned but its not really scaned fully so we can invertatly set a struct
					thats not parsed in this package as required and it would get overridden when it is parsed subsiquently
					we just need to persist this to that as well
				*/
				fullStructNameAfterRemovingGenerics := removeGenericsPartFromStructName(obj.Type().String())
				if oldVal, exists := c.Structs[fullStructNameAfterRemovingGenerics]; exists {
					parsedStruct.Required = oldVal.Required
				}
				c.Structs[fullStructNameAfterRemovingGenerics] = parsedStruct
			}
		}
	}
	return RequestedStruct
}

func (c *Converter) getGenericPopulations(structName string) []ParsedField {
	toRet := []ParsedField{}
	structNameSegments := strings.Split(structName, "[")
	if len(structNameSegments) > 1 {
		genericPopulations := strings.Replace(structNameSegments[1], "]", "", 1)
		for _, v := range strings.Split(genericPopulations, ",") {
			structSegments := strings.Split(v, ".")
			toRet = append(toRet, ParsedField{
				TSType: structSegments[len(structSegments)-1],
			})
		}
	}
	return toRet
}

func isStructWithPackagePathEqualToRequiredStructName(RequiredStructName, PackagePath string) bool {
	PackageStructSegments := strings.Split(PackagePath, ".")
	return removeGenericsPartFromStructName(RequiredStructName) == removeGenericsPartFromStructName(PackageStructSegments[len(PackageStructSegments)-1])
}

func removeGenericsPartFromStructName(name string) string {
	return strings.Split(name, "[")[0]
}

func (c *Converter) GetStructAsInterfaceString(ps ParsedStruct) string {
	var toRet string = ""
	doc, docExists := c.Docs[ps.Name]
	if docExists {
		if doc != "" {
			toRet += c.GetFormattedTSComment(doc) + "\n"
		}
	}
	if ps.Name == "" {
		return ""
	}
	toRet += "export interface " + GetFormattedInterfaceName(ps.Name) + " {"
	for _, v := range ps.Fields {
		if v.Var.Embedded() {
			for _, f := range c.Structs[v.Var.Pkg().Path()+"."+v.TSName].Fields {
				toRet += "\n" + c.Indent + f.TSName + ": " + f.TSType
				if f.IsSlice {
					toRet += "[]"
				}
			}
		} else {
			toRet += "\n" + c.Indent + v.TSName + ": " + v.TSType
			if v.IsSlice {
				toRet += "[]"
			}
		}
	}
	toRet += "\n}"
	return toRet
}

func (c *Converter) GetFormattedTSComment(commentContent string) string {
	result := "\n/**\n"
	for _, v := range strings.Split(commentContent, "\n") {
		if v != "" {
			result += c.Indent + v + "\n"
		}
	}
	result += "*/"
	return result
}

func GetFormattedInterfaceName(name string) string {
	// remove all typing like any and stuff from the generics []'s
	nameSegments := strings.Split(name, "[")
	if len(strings.Split(name, "[")) > 1 {
		genericsPart := strings.ReplaceAll(nameSegments[1], "]", "")
		name = nameSegments[0]
		name += "<"
		for i, v := range strings.Split(genericsPart, ",") {
			if i != 0 {
				name += ", "
			}
			name += strings.Split(strings.Trim(v, " "), " ")[0]
		}
		name += ">"
	} else {
		name = strings.ReplaceAll(name, "[", "<")
		name = strings.ReplaceAll(name, "]", ">")
	}
	return name
}
