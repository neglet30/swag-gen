package swagger

import (
	"fmt"
	"reflect"
	"strings"
)

// SchemaBuilder is responsible for building JSON schemas from Go types.
type SchemaBuilder struct {
	schemas map[string]*Schema
}

// NewSchemaBuilder creates a new schema builder.
func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{
		schemas: make(map[string]*Schema),
	}
}

// BuildSchema builds a schema from a Go type string.
func (sb *SchemaBuilder) BuildSchema(typeStr string) *Schema {
	return sb.buildSchemaFromType(typeStr)
}

// buildSchemaFromType builds a schema from a type string.
func (sb *SchemaBuilder) buildSchemaFromType(typeStr string) *Schema {
	typeStr = strings.TrimSpace(typeStr)

	// Handle pointer types
	if strings.HasPrefix(typeStr, "*") {
		return sb.buildSchemaFromType(typeStr[1:])
	}

	// Handle array types
	if strings.HasPrefix(typeStr, "[]") {
		return &Schema{
			Type:  "array",
			Items: sb.buildSchemaFromType(typeStr[2:]),
		}
	}

	// Handle map types
	if strings.HasPrefix(typeStr, "map[") {
		return &Schema{
			Type: "object",
		}
	}

	// Handle basic types
	switch typeStr {
	case "string":
		return &Schema{Type: "string"}
	case "int", "int8", "int16", "int32", "int64":
		return &Schema{Type: "integer", Format: "int64"}
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return &Schema{Type: "integer", Format: "int64"}
	case "float32":
		return &Schema{Type: "number", Format: "float"}
	case "float64":
		return &Schema{Type: "number", Format: "double"}
	case "bool":
		return &Schema{Type: "boolean"}
	case "byte":
		return &Schema{Type: "string", Format: "byte"}
	case "rune":
		return &Schema{Type: "integer", Format: "int32"}
	case "time.Time":
		return &Schema{Type: "string", Format: "date-time"}
	case "time.Duration":
		return &Schema{Type: "string"}
	case "interface{}":
		return &Schema{}
	default:
		// For custom types, return a reference
		return &Schema{
			Ref: fmt.Sprintf("#/components/schemas/%s", typeStr),
		}
	}
}

// BuildStructSchema builds a schema from a struct definition.
func (sb *SchemaBuilder) BuildStructSchema(name string, fields map[string]string) *Schema {
	schema := &Schema{
		Type:       "object",
		Properties: make(map[string]*Schema),
		Required:   make([]string, 0),
	}

	for fieldName, fieldType := range fields {
		fieldSchema := sb.buildSchemaFromType(fieldType)
		schema.Properties[fieldName] = fieldSchema

		// Mark as required if not a pointer type
		if !strings.HasPrefix(fieldType, "*") {
			schema.Required = append(schema.Required, fieldName)
		}
	}

	sb.schemas[name] = schema
	return schema
}

// GetSchemas returns all built schemas.
func (sb *SchemaBuilder) GetSchemas() map[string]*Schema {
	return sb.schemas
}

// GetSchema returns a specific schema by name.
func (sb *SchemaBuilder) GetSchema(name string) *Schema {
	return sb.schemas[name]
}

// BuildFromReflect builds a schema from a Go reflect.Type.
func (sb *SchemaBuilder) BuildFromReflect(t reflect.Type) *Schema {
	if t == nil {
		return &Schema{}
	}

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		return sb.BuildFromReflect(t.Elem())
	}

	// Handle array types
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		return &Schema{
			Type:  "array",
			Items: sb.BuildFromReflect(t.Elem()),
		}
	}

	// Handle map types
	if t.Kind() == reflect.Map {
		return &Schema{
			Type: "object",
		}
	}

	// Handle struct types
	if t.Kind() == reflect.Struct {
		schema := &Schema{
			Type:       "object",
			Properties: make(map[string]*Schema),
			Required:   make([]string, 0),
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Anonymous {
				continue
			}

			fieldSchema := sb.BuildFromReflect(field.Type)
			schema.Properties[field.Name] = fieldSchema

			// Check if field is required (not a pointer)
			if field.Type.Kind() != reflect.Ptr {
				schema.Required = append(schema.Required, field.Name)
			}
		}

		return schema
	}

	// Handle basic types
	switch t.Kind() {
	case reflect.String:
		return &Schema{Type: "string"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &Schema{Type: "integer", Format: "int64"}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &Schema{Type: "integer", Format: "int64"}
	case reflect.Float32:
		return &Schema{Type: "number", Format: "float"}
	case reflect.Float64:
		return &Schema{Type: "number", Format: "double"}
	case reflect.Bool:
		return &Schema{Type: "boolean"}
	default:
		return &Schema{}
	}
}
