package casesHelper

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	_stringUtils "github.com/aaronchen2k/deeptest/pkg/lib/string"
	"github.com/getkin/kin-openapi/openapi3"
	"math"
)

func addPropCase(propName string, propVal *openapi3.Schema, requires []string, parent *AlternativeCase) {
	if propVal.Type == OasFieldTypeArray.String() {
		arrCase := &AlternativeCase{
			Title:    "items",
			Category: consts.AlternativeCaseProp,
			IsDir:    true,
		}

		addPropCase(propName, propVal.Items.Value, nil, arrCase)

		parent.Children = append(parent.Children, arrCase)

		return

	} else if propVal.Type == OasFieldTypeObject.String() {
		objCase := &AlternativeCase{
			Title:    "object",
			Category: consts.AlternativeCaseObject,
			IsDir:    true,
		}

		for propName, propRef := range propVal.Properties {
			addPropCase(propName, propRef.Value, propVal.Required, objCase)
		}

		parent.Children = append(parent.Children, objCase)

		return
	}

	addPropRequiredCase(propName, propVal, requires, parent)
	addPropTypeCase(propVal, parent)
	addPropEnumCase(propVal, parent)
	addPropFormatCase(propVal, parent)
	addPropRuleCase(propVal, parent)
}

func addPropRequiredCase(propName string, schemaVal *openapi3.Schema, requires []string, parent *AlternativeCase) {
	if !_stringUtils.StrInArr(propName, requires) {
		return
	}

	required := &AlternativeCase{
		Sample: ExampleEmpty,

		Category:      consts.AlternativeCaseCase,
		Type:          consts.AlternativeCaseRequired,
		FieldRequired: true,
		IsDir:         false,
	}

	parent.Children = append(parent.Children, required)
}

func addPropTypeCase(schema *openapi3.Schema, parent *AlternativeCase) {
	typ := OasFieldType(schema.Type)

	if typ == OasFieldTypeAny || typ == OasFieldTypeString {
		return
	}

	var sample interface{}
	if typ == OasFieldTypeBoolean || typ == OasFieldTypeNumber || typ == OasFieldTypeArray {
		sample = RandStr()
	} else if typ == OasFieldTypeInteger {
		sample = RandFloat32()
	}

	if sample == nil {
		return
	}

	typeCase := &AlternativeCase{
		Sample: sample,

		Category:  consts.AlternativeCaseCase,
		Type:      consts.AlternativeCaseTyped,
		FieldType: typ,
		IsDir:     false,
	}

	parent.Children = append(parent.Children, typeCase)
}

func addPropEnumCase(schema *openapi3.Schema, parent *AlternativeCase) {
	enum := schema.Enum

	if enum == nil {
		return
	}

	enumCase := &AlternativeCase{
		Sample: RandStr(),

		Category:  consts.AlternativeCaseCase,
		Type:      consts.AlternativeCaseEnum,
		FieldType: OasFieldType(schema.Type),
		IsDir:     false,
	}

	parent.Children = append(parent.Children, enumCase)
}

func addPropFormatCase(schema *openapi3.Schema, parent *AlternativeCase) {
	typ := OasFieldType(schema.Type)
	format := OasFieldFormat(schema.Format)

	if format == "" {
		return
	}

	var sample interface{}
	if typ == OasFieldTypeInteger {
		if format == OasFieldFormatInt32 {
			sample = RandInt64()
		} else if format == OasFieldFormatInt64 {
			sample = RandStr()
		}
	} else if typ == OasFieldTypeNumber {
		if format == OasFieldFormatFloat {
			sample = RandFloat64()
		} else if format == OasFieldFormatDouble {
			sample = RandStr()
		}
	} else {
		sample = RandStr()
	}

	formatCase := &AlternativeCase{
		Sample: sample,

		Category:  consts.AlternativeCaseCase,
		Type:      consts.AlternativeCaseFormat,
		FieldType: OasFieldType(schema.Type),
		IsDir:     false,
	}

	parent.Children = append(parent.Children, formatCase)
}

func addPropRuleCase(schema *openapi3.Schema, parent *AlternativeCase) {
	typ := OasFieldType(schema.Type)

	var sample interface{}
	if typ == OasFieldTypeInteger || typ == OasFieldTypeNumber {
		if schema.Min != nil && *schema.Min != 0 {
			sample = *schema.Min - 1
		} else if schema.Max != nil && *schema.Max != 0 {
			sample = *schema.Max + 1
		} else if schema.MaxLength != nil && *schema.MaxLength > 0 {
			if typ == OasFieldTypeInteger {
				sample = 1 * math.Pow(10, float64(*schema.MaxLength))
			} else {
				sample = 1 / math.Pow(10, float64(*schema.MaxLength-1))
			}
		} else if schema.MinLength > 0 {
			sample = 1
		} else if schema.MultipleOf != nil && *schema.MultipleOf != 0 {
			if typ == OasFieldTypeInteger {
				sample = *schema.MultipleOf + 1
			} else {
				sample = *schema.MultipleOf + *schema.MultipleOf*0.1
			}
		} else if schema.ExclusiveMin {
			sample = *schema.Min
		} else if schema.ExclusiveMax {
			sample = *schema.Max
		}

	} else {
		if schema.Pattern != "" {
			sample = RandStrSpecial()
		} else if schema.MaxLength != nil && *(schema.MaxLength) > 0 {
			sample = RandStrWithLen(int(*(schema.MaxLength) + 1))
		} else if schema.MinLength > 0 {
			sample = RandStrWithLen(int(schema.MinLength - 1))
		}
	}

	if sample == nil {
		return
	}

	ruleCase := &AlternativeCase{
		Sample: sample,

		Category:  consts.AlternativeCaseCase,
		Type:      consts.AlternativeCaseRule,
		FieldType: OasFieldType(schema.Type),
		IsDir:     false,
	}

	parent.Children = append(parent.Children, ruleCase)
}
