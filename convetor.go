package gos2tsi

import (
	"go/doc"
	"go/types"
	"reflect"
	"strings"

	"golang.org/x/tools/go/packages"
)

var GoTypeToTSType = map[string]string{
	"bool":        "boolean",
	"interface{}": "any",
	"any":         "any",
	"string":      "string",
	"int":         "number",
	"int8":        "number",
	"int16":       "number",
	"int32":       "number",
	"int64":       "number",
	"uint":        "number",
	"uint8":       "number",
	"uint16":      "number",
	"uint32":      "number",
	"uint64":      "number",
	"float32":     "number",
	"float64":     "number",
}

type ParsedField struct {
	Var       *types.Var
	Tag       string
	TSName    string
	TSType    string
	IsSlice   int
	RefStruct ParsedStruct
}

type ParsedStruct struct {
	PackageName        string
	PackgePath         string
	ID                 string
	Name               string
	Required           bool
	Fields             []ParsedField
	IsSlice            int
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
	IsSlice := 0
	if reflectType.Kind() == reflect.Slice {
		IsSlice = 1
		reflectType = reflectType.Elem()
	}
	pkgPath := reflectType.PkgPath()
	RequiredStruct := reflectType.Name()
	return c.ParseStructsInPackage(pkgPath, RequiredStruct, IsSlice)
}

func (c *Converter) ensureGenericPopulations(structName string) {
	arrayTrimmedStructName := strings.Trim(structName, "[]")
	genericPopulations := strings.SplitN(arrayTrimmedStructName, "[", 2)
	if len(genericPopulations) > 1 {
		closingBracketTrimmed := strings.Trim(genericPopulations[1], "]")
		for _, gp := range strings.Split(closingBracketTrimmed, ",") {
			arrayTrimmed := strings.Trim(gp, "[]")
			fullNameSegments := strings.Split(arrayTrimmed, ".")
			if len(fullNameSegments) > 1 {
				isSlice := 0
				if strings.HasPrefix(gp, "[]") {
					isSlice = 1
				}
				pkgPath := strings.Join(fullNameSegments[:len(fullNameSegments)-1], ".")
				structName := fullNameSegments[len(fullNameSegments)-1]
				c.ParseStructsInPackage(pkgPath, structName, isSlice)
			}
		}
	}
}

func (c *Converter) ParseStructsInPackage(pkgPath, RequiredStruct string, IsSlice int) ParsedStruct {
	RequestedStruct := ParsedStruct{}
	RequestedStruct.IsSlice = IsSlice
	RequestedStruct.GenericPopulations = c.getGenericPopulations(RequiredStruct)
	c.ensureGenericPopulations(RequiredStruct)
	_, exists := c.AlreadyParsedPackage[pkgPath]
	if exists {
		rs := c.Structs[pkgPath+"."+removeGenericsPartFromStructName(RequiredStruct)]
		rs.Required = true
		c.Structs[pkgPath+"."+removeGenericsPartFromStructName(RequiredStruct)] = rs
		rs.IsSlice = RequestedStruct.IsSlice
		rs.GenericPopulations = RequestedStruct.GenericPopulations
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
					FieldTypeName := pf.Var.Type().String()
					if reflect.TypeOf(st.Field(i).Type()).String() == "*types.Slice" {
						pf.IsSlice = countPrefixBrackets(FieldTypeName)
					}
					if strings.HasPrefix(typeName, "[]") {
						typeName = strings.Replace(typeName, "[]", "", pf.IsSlice)
					}
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
					pf.TSName = fieldName
					pf.TSType = typeName
					if pf.TSName != "-" {
						if !isTypeValid && len(typeNameSegments) > 1 && !pf.Var.Embedded() {
							if pf.IsSlice > 0 {
								FieldTypeName = strings.Replace(FieldTypeName, "[]", "", pf.IsSlice)
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
					parsedStruct.IsSlice = RequestedStruct.IsSlice
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

func countPrefixBrackets(line string) int {
	count := 0
	prefix := "[]"

	// Keep checking and removing the prefix "[]" from the start of the line
	for strings.HasPrefix(line, prefix) {
		count++
		line = strings.TrimPrefix(line, prefix)
	}

	return count
}

func (c *Converter) getGenericPopulations(structName string) []ParsedField {
	toRet := []ParsedField{}
	arrayTrimmedStructName := strings.Trim(structName, "[]")
	genericPopulations := strings.SplitN(arrayTrimmedStructName, "[", 2)
	if len(genericPopulations) > 1 {
		closingBracketTrimmed := strings.Trim(genericPopulations[1], "]")
		for _, gp := range strings.Split(closingBracketTrimmed, ",") {
			arrayTrimmed := strings.Trim(gp, "[]")
			fullNameSegments := strings.Split(arrayTrimmed, ".")
			isSlice := 0
			if strings.HasPrefix(gp, "[]") {
				isSlice = 1
			}
			fieldName := fullNameSegments[len(fullNameSegments)-1]
			toRet = append(toRet, ParsedField{
				IsSlice: isSlice,
				TSType:  fieldName,
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
			pkgPath, structName := c.GetPackagePathAndStructNameFromFullDenotation(v.Var.Type().String())
			ps := c.ParseStructsInPackage(pkgPath, structName, v.IsSlice)
			ps = c.SetGenericPopulationsToFields(ps)
			for _, f := range ps.Fields {
				toRet += c.GetFieldAsString(f)
			}
			// embeddedPackageName := v.Var.Pkg().Path() + "." + v.TSName
			// if strings.Contains(v.TSName, ".") {
			// 	embeddedPackageName = v.TSName
			// }
			// for _, f := range c.Structs[embeddedPackageName].Fields {
			// 	toRet += c.GetFieldAsString(f)
			// }
		} else {
			toRet += c.GetFieldAsString(v)
		}
	}
	toRet += "\n}"
	return toRet
}

func (c *Converter) GetPackagePathAndStructNameFromFullDenotation(fullPath string) (string, string) {
	woGenerics := strings.Split(fullPath, "[")[0]
	woGenericSegments := strings.Split(woGenerics, ".")
	StructName := woGenericSegments[len(woGenericSegments)-1]
	pkgPath := strings.Replace(woGenerics, StructName, "", 1)
	StructNameWithGenerics := strings.Replace(fullPath, pkgPath, "", 1)
	return pkgPath[:len(pkgPath)-1], StructNameWithGenerics
}

func (c *Converter) SetGenericPopulationsToFields(ps ParsedStruct) ParsedStruct {
	if !strings.Contains(ps.Name, "[") {
		return ps
	}
	var replaceMentMap map[string]string = make(map[string]string)
	genericSegments := strings.Split(ps.Name, "[")
	generics := strings.Replace(genericSegments[1], "]", "", 1)
	genericParts := strings.Split(generics, ",")
	for i, v := range genericParts {
		genericPartKey := strings.Split(strings.Trim(v, " "), " ")[0]
		replaceMentMap[genericPartKey] = ps.GenericPopulations[i].TSType
	}
	for i, field := range ps.Fields {
		typeTokens := strings.Split(field.TSType, "[")
		var reconstructed []string
		for _, tt := range typeTokens {
			typeTokens2 := strings.Split(tt, "]")
			var reconstructed2 []string
			for _, tt2 := range typeTokens2 {
				for k, v := range replaceMentMap {
					if tt2 == k {
						tt2 = v
					}
				}
				reconstructed2 = append(reconstructed2, tt2)
			}
			reconstructed = append(reconstructed, strings.Join(reconstructed2, "]"))
		}
		ps.Fields[i].TSType = strings.Join(reconstructed, "[")
	}
	return ps
}

func (c *Converter) GetFieldAsString(pf ParsedField) string {
	toRet := "\n" + c.Indent + pf.TSName + ": " + c.postProcessTSTypeName(pf.TSType)
	for i := 0; i < pf.IsSlice; i++ {
		toRet += "[]"
	}
	return toRet
}

func (c *Converter) postProcessTSTypeName(TSType string) string {
	if strings.HasPrefix(TSType, "map[") {
		return c.GetTSTypeFromMap(TSType)
	}
	return TSType
}

func (c *Converter) GetTSTypeFromMap(TSType string) string {
	arrayIndication := strings.Split(TSType, "map[")[0]
	if arrayIndication != "" {
		TSType = strings.Replace(TSType, arrayIndication, "", 1)
	}
	TSType = strings.Replace(TSType, "map[", "", 1)
	TSTypeSegments := strings.Split(TSType, "]")
	KeyType := TSTypeSegments[0]
	if c.isMap(KeyType) {
		KeyType = c.GetTSTypeFromMap(KeyType)
	}
	if convertedTypeName, ok := GoTypeToTSType[KeyType]; ok {
		KeyType = convertedTypeName
	}
	ValueType := strings.Join(TSTypeSegments[1:], "]")
	if c.isMap(ValueType) {
		ValueType = c.GetTSTypeFromMap(ValueType)
	}
	if convertedTypeName, ok := GoTypeToTSType[ValueType]; ok {
		ValueType = convertedTypeName
	}
	return `{[key: ` + KeyType + `]: ` + ValueType + `}` + arrayIndication
}

func (c *Converter) isMap(i string) bool {
	if strings.HasPrefix(i, "[]") {
		i = strings.Replace(i, "[]", "", 1)
	}
	if strings.HasPrefix(i, "map[") {
		return true
	}
	return false
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
