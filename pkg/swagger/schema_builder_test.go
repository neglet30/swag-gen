package swagger

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSchemaBuilder(t *testing.T) {
	sb := NewSchemaBuilder()

	assert.NotNil(t, sb)
	assert.NotNil(t, sb.schemas)
	assert.Empty(t, sb.schemas)
}

func TestSchemaBuilderBuildSchema_BasicTypes(t *testing.T) {
	tests := []struct {
		name     string
		typeStr  string
		wantType string
		wantFmt  string
	}{
		{"string", "string", "string", ""},
		{"int", "int", "integer", "int64"},
		{"int8", "int8", "integer", "int64"},
		{"int16", "int16", "integer", "int64"},
		{"int32", "int32", "integer", "int64"},
		{"int64", "int64", "integer", "int64"},
		{"uint", "uint", "integer", "int64"},
		{"uint8", "uint8", "integer", "int64"},
		{"uint16", "uint16", "integer", "int64"},
		{"uint32", "uint32", "integer", "int64"},
		{"uint64", "uint64", "integer", "int64"},
		{"float32", "float32", "number", "float"},
		{"float64", "float64", "number", "double"},
		{"bool", "bool", "boolean", ""},
		{"byte", "byte", "string", "byte"},
		{"rune", "rune", "integer", "int32"},
		{"time.Time", "time.Time", "string", "date-time"},
		{"time.Duration", "time.Duration", "string", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewSchemaBuilder()
			schema := sb.BuildSchema(tt.typeStr)

			assert.NotNil(t, schema)
			assert.Equal(t, tt.wantType, schema.Type)
			if tt.wantFmt != "" {
				assert.Equal(t, tt.wantFmt, schema.Format)
			}
		})
	}
}

func TestSchemaBuilderBuildSchema_PointerTypes(t *testing.T) {
	sb := NewSchemaBuilder()

	schema := sb.BuildSchema("*string")
	assert.NotNil(t, schema)
	assert.Equal(t, "string", schema.Type)
}

func TestSchemaBuilderBuildSchema_ArrayTypes(t *testing.T) {
	sb := NewSchemaBuilder()

	schema := sb.BuildSchema("[]string")
	assert.NotNil(t, schema)
	assert.Equal(t, "array", schema.Type)
	assert.NotNil(t, schema.Items)
	assert.Equal(t, "string", schema.Items.Type)
}

func TestSchemaBuilderBuildSchema_NestedArrayTypes(t *testing.T) {
	sb := NewSchemaBuilder()

	schema := sb.BuildSchema("[][]int")
	assert.NotNil(t, schema)
	assert.Equal(t, "array", schema.Type)
	assert.NotNil(t, schema.Items)
	assert.Equal(t, "array", schema.Items.Type)
	assert.NotNil(t, schema.Items.Items)
	assert.Equal(t, "integer", schema.Items.Items.Type)
}

func TestSchemaBuilderBuildSchema_MapTypes(t *testing.T) {
	sb := NewSchemaBuilder()

	schema := sb.BuildSchema("map[string]interface{}")
	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
}

func TestSchemaBuilderBuildSchema_CustomTypes(t *testing.T) {
	sb := NewSchemaBuilder()

	schema := sb.BuildSchema("User")
	assert.NotNil(t, schema)
	assert.Equal(t, "#/components/schemas/User", schema.Ref)
}

func TestSchemaBuilderBuildSchema_InterfaceType(t *testing.T) {
	sb := NewSchemaBuilder()

	schema := sb.BuildSchema("interface{}")
	assert.NotNil(t, schema)
	assert.Equal(t, "", schema.Type)
}

func TestSchemaBuilderBuildStructSchema(t *testing.T) {
	sb := NewSchemaBuilder()

	fields := map[string]string{
		"id":   "int",
		"name": "string",
		"age":  "int",
	}

	schema := sb.BuildStructSchema("User", fields)

	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
	assert.Len(t, schema.Properties, 3)
	assert.NotNil(t, schema.Properties["id"])
	assert.NotNil(t, schema.Properties["name"])
	assert.NotNil(t, schema.Properties["age"])
	assert.Len(t, schema.Required, 3)
}

func TestSchemaBuilderBuildStructSchema_WithPointers(t *testing.T) {
	sb := NewSchemaBuilder()

	fields := map[string]string{
		"id":    "int",
		"name":  "string",
		"email": "*string",
		"phone": "*string",
	}

	schema := sb.BuildStructSchema("User", fields)

	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
	assert.Len(t, schema.Properties, 4)
	// Only non-pointer fields should be required
	assert.Len(t, schema.Required, 2)
	assert.Contains(t, schema.Required, "id")
	assert.Contains(t, schema.Required, "name")
}

func TestSchemaBuilderGetSchemas(t *testing.T) {
	sb := NewSchemaBuilder()

	fields1 := map[string]string{"id": "int", "name": "string"}
	fields2 := map[string]string{"id": "int", "email": "string"}

	sb.BuildStructSchema("User", fields1)
	sb.BuildStructSchema("Admin", fields2)

	schemas := sb.GetSchemas()
	assert.Len(t, schemas, 2)
	assert.NotNil(t, schemas["User"])
	assert.NotNil(t, schemas["Admin"])
}

func TestSchemaBuilderGetSchema(t *testing.T) {
	sb := NewSchemaBuilder()

	fields := map[string]string{"id": "int", "name": "string"}
	sb.BuildStructSchema("User", fields)

	schema := sb.GetSchema("User")
	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
	assert.Len(t, schema.Properties, 2)
}

func TestSchemaBuilderGetSchema_NotFound(t *testing.T) {
	sb := NewSchemaBuilder()

	schema := sb.GetSchema("NonExistent")
	assert.Nil(t, schema)
}

func TestSchemaBuilderBuildSchema_WithWhitespace(t *testing.T) {
	sb := NewSchemaBuilder()

	schema := sb.BuildSchema("  string  ")
	assert.NotNil(t, schema)
	assert.Equal(t, "string", schema.Type)
}

func TestSchemaBuilderBuildSchema_PointerToArray(t *testing.T) {
	sb := NewSchemaBuilder()

	schema := sb.BuildSchema("*[]string")
	assert.NotNil(t, schema)
	assert.Equal(t, "array", schema.Type)
	assert.NotNil(t, schema.Items)
	assert.Equal(t, "string", schema.Items.Type)
}

func TestSchemaBuilderBuildStructSchema_Empty(t *testing.T) {
	sb := NewSchemaBuilder()

	schema := sb.BuildStructSchema("Empty", map[string]string{})

	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
	assert.Empty(t, schema.Properties)
	assert.Empty(t, schema.Required)
}

func TestSchemaBuilderBuildStructSchema_ComplexTypes(t *testing.T) {
	sb := NewSchemaBuilder()

	fields := map[string]string{
		"id":       "int",
		"name":     "string",
		"tags":     "[]string",
		"metadata": "map[string]interface{}",
		"user":     "User",
	}

	schema := sb.BuildStructSchema("Post", fields)

	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
	assert.Len(t, schema.Properties, 5)

	// Check array type
	assert.Equal(t, "array", schema.Properties["tags"].Type)
	assert.Equal(t, "string", schema.Properties["tags"].Items.Type)

	// Check map type
	assert.Equal(t, "object", schema.Properties["metadata"].Type)

	// Check custom type reference
	assert.Equal(t, "#/components/schemas/User", schema.Properties["user"].Ref)
}

func TestSchemaBuilderMultipleStructs(t *testing.T) {
	sb := NewSchemaBuilder()

	userFields := map[string]string{
		"id":   "int",
		"name": "string",
	}

	postFields := map[string]string{
		"id":     "int",
		"title":  "string",
		"author": "User",
	}

	sb.BuildStructSchema("User", userFields)
	sb.BuildStructSchema("Post", postFields)

	schemas := sb.GetSchemas()
	assert.Len(t, schemas, 2)

	userSchema := sb.GetSchema("User")
	assert.Len(t, userSchema.Properties, 2)

	postSchema := sb.GetSchema("Post")
	assert.Len(t, postSchema.Properties, 3)
	assert.Equal(t, "#/components/schemas/User", postSchema.Properties["author"].Ref)
}

func TestSchemaBuilderBuildSchema_ArrayOfPointers(t *testing.T) {
	sb := NewSchemaBuilder()

	schema := sb.BuildSchema("[]*User")
	assert.NotNil(t, schema)
	assert.Equal(t, "array", schema.Type)
	assert.NotNil(t, schema.Items)
	assert.Equal(t, "#/components/schemas/User", schema.Items.Ref)
}

func TestSchemaBuilderBuildSchema_PointerToArrayOfPointers(t *testing.T) {
	sb := NewSchemaBuilder()

	schema := sb.BuildSchema("*[]*string")
	assert.NotNil(t, schema)
	assert.Equal(t, "array", schema.Type)
	assert.NotNil(t, schema.Items)
	assert.Equal(t, "string", schema.Items.Type)
}

func TestSchemaBuilderBuildFromReflect_String(t *testing.T) {
	sb := NewSchemaBuilder()

	var s string
	schema := sb.BuildFromReflect(reflect.TypeOf(s))

	assert.NotNil(t, schema)
	assert.Equal(t, "string", schema.Type)
}

func TestSchemaBuilderBuildFromReflect_Int(t *testing.T) {
	sb := NewSchemaBuilder()

	var i int
	schema := sb.BuildFromReflect(reflect.TypeOf(i))

	assert.NotNil(t, schema)
	assert.Equal(t, "integer", schema.Type)
	assert.Equal(t, "int64", schema.Format)
}

func TestSchemaBuilderBuildFromReflect_Float64(t *testing.T) {
	sb := NewSchemaBuilder()

	var f float64
	schema := sb.BuildFromReflect(reflect.TypeOf(f))

	assert.NotNil(t, schema)
	assert.Equal(t, "number", schema.Type)
	assert.Equal(t, "double", schema.Format)
}

func TestSchemaBuilderBuildFromReflect_Bool(t *testing.T) {
	sb := NewSchemaBuilder()

	var b bool
	schema := sb.BuildFromReflect(reflect.TypeOf(b))

	assert.NotNil(t, schema)
	assert.Equal(t, "boolean", schema.Type)
}

func TestSchemaBuilderBuildFromReflect_Slice(t *testing.T) {
	sb := NewSchemaBuilder()

	var slice []string
	schema := sb.BuildFromReflect(reflect.TypeOf(slice))

	assert.NotNil(t, schema)
	assert.Equal(t, "array", schema.Type)
	assert.NotNil(t, schema.Items)
	assert.Equal(t, "string", schema.Items.Type)
}

func TestSchemaBuilderBuildFromReflect_Array(t *testing.T) {
	sb := NewSchemaBuilder()

	var arr [5]int
	schema := sb.BuildFromReflect(reflect.TypeOf(arr))

	assert.NotNil(t, schema)
	assert.Equal(t, "array", schema.Type)
	assert.NotNil(t, schema.Items)
	assert.Equal(t, "integer", schema.Items.Type)
}

func TestSchemaBuilderBuildFromReflect_Map(t *testing.T) {
	sb := NewSchemaBuilder()

	var m map[string]interface{}
	schema := sb.BuildFromReflect(reflect.TypeOf(m))

	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
}

func TestSchemaBuilderBuildFromReflect_Pointer(t *testing.T) {
	sb := NewSchemaBuilder()

	var s *string
	schema := sb.BuildFromReflect(reflect.TypeOf(s))

	assert.NotNil(t, schema)
	assert.Equal(t, "string", schema.Type)
}

func TestSchemaBuilderBuildFromReflect_Struct(t *testing.T) {
	sb := NewSchemaBuilder()

	type User struct {
		ID   int
		Name string
	}

	schema := sb.BuildFromReflect(reflect.TypeOf(User{}))

	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
	assert.Len(t, schema.Properties, 2)
	assert.NotNil(t, schema.Properties["ID"])
	assert.NotNil(t, schema.Properties["Name"])
	assert.Len(t, schema.Required, 2)
}

func TestSchemaBuilderBuildFromReflect_StructWithPointers(t *testing.T) {
	sb := NewSchemaBuilder()

	type User struct {
		ID    int
		Name  string
		Email *string
	}

	schema := sb.BuildFromReflect(reflect.TypeOf(User{}))

	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
	assert.Len(t, schema.Properties, 3)
	assert.Len(t, schema.Required, 2)
	assert.Contains(t, schema.Required, "ID")
	assert.Contains(t, schema.Required, "Name")
}

func TestSchemaBuilderBuildFromReflect_Nil(t *testing.T) {
	sb := NewSchemaBuilder()

	schema := sb.BuildFromReflect(nil)

	assert.NotNil(t, schema)
	assert.Equal(t, "", schema.Type)
}

func TestSchemaBuilderBuildFromReflect_NestedStruct(t *testing.T) {
	sb := NewSchemaBuilder()

	type Address struct {
		Street string
		City   string
	}

	type User struct {
		ID      int
		Name    string
		Address Address
	}

	schema := sb.BuildFromReflect(reflect.TypeOf(User{}))

	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
	assert.Len(t, schema.Properties, 3)
	assert.NotNil(t, schema.Properties["Address"])
	assert.Equal(t, "object", schema.Properties["Address"].Type)
}

func TestSchemaBuilderBuildFromReflect_SliceOfStructs(t *testing.T) {
	sb := NewSchemaBuilder()

	type User struct {
		ID   int
		Name string
	}

	var users []User
	schema := sb.BuildFromReflect(reflect.TypeOf(users))

	assert.NotNil(t, schema)
	assert.Equal(t, "array", schema.Type)
	assert.NotNil(t, schema.Items)
	assert.Equal(t, "object", schema.Items.Type)
}

func TestSchemaBuilderBuildFromReflect_AllIntTypes(t *testing.T) {
	tests := []struct {
		name     string
		typeFunc func() reflect.Type
		wantType string
		wantFmt  string
	}{
		{"int8", func() reflect.Type { var i int8; return reflect.TypeOf(i) }, "integer", "int64"},
		{"int16", func() reflect.Type { var i int16; return reflect.TypeOf(i) }, "integer", "int64"},
		{"int32", func() reflect.Type { var i int32; return reflect.TypeOf(i) }, "integer", "int64"},
		{"int64", func() reflect.Type { var i int64; return reflect.TypeOf(i) }, "integer", "int64"},
		{"uint", func() reflect.Type { var i uint; return reflect.TypeOf(i) }, "integer", "int64"},
		{"uint8", func() reflect.Type { var i uint8; return reflect.TypeOf(i) }, "integer", "int64"},
		{"uint16", func() reflect.Type { var i uint16; return reflect.TypeOf(i) }, "integer", "int64"},
		{"uint32", func() reflect.Type { var i uint32; return reflect.TypeOf(i) }, "integer", "int64"},
		{"uint64", func() reflect.Type { var i uint64; return reflect.TypeOf(i) }, "integer", "int64"},
		{"float32", func() reflect.Type { var f float32; return reflect.TypeOf(f) }, "number", "float"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewSchemaBuilder()
			schema := sb.BuildFromReflect(tt.typeFunc())

			assert.NotNil(t, schema)
			assert.Equal(t, tt.wantType, schema.Type)
			assert.Equal(t, tt.wantFmt, schema.Format)
		})
	}
}

func TestSchemaBuilderBuildSchema_ComplexNesting(t *testing.T) {
	sb := NewSchemaBuilder()

	// Test deeply nested array types
	schema := sb.BuildSchema("[][][]string")
	assert.NotNil(t, schema)
	assert.Equal(t, "array", schema.Type)
	assert.NotNil(t, schema.Items)
	assert.Equal(t, "array", schema.Items.Type)
	assert.NotNil(t, schema.Items.Items)
	assert.Equal(t, "array", schema.Items.Items.Type)
	assert.NotNil(t, schema.Items.Items.Items)
	assert.Equal(t, "string", schema.Items.Items.Items.Type)
}

func TestSchemaBuilderBuildStructSchema_MixedTypes(t *testing.T) {
	sb := NewSchemaBuilder()

	fields := map[string]string{
		"id":        "int64",
		"name":      "string",
		"active":    "bool",
		"score":     "float64",
		"tags":      "[]string",
		"metadata":  "map[string]interface{}",
		"createdAt": "time.Time",
		"updatedAt": "*time.Time",
	}

	schema := sb.BuildStructSchema("ComplexModel", fields)

	assert.NotNil(t, schema)
	assert.Equal(t, "object", schema.Type)
	assert.Len(t, schema.Properties, 8)

	// Check specific types
	assert.Equal(t, "integer", schema.Properties["id"].Type)
	assert.Equal(t, "string", schema.Properties["name"].Type)
	assert.Equal(t, "boolean", schema.Properties["active"].Type)
	assert.Equal(t, "number", schema.Properties["score"].Type)
	assert.Equal(t, "array", schema.Properties["tags"].Type)
	assert.Equal(t, "object", schema.Properties["metadata"].Type)
	assert.Equal(t, "string", schema.Properties["createdAt"].Type)
	assert.Equal(t, "date-time", schema.Properties["createdAt"].Format)

	// Check required fields (non-pointer types)
	assert.Len(t, schema.Required, 7)
	assert.NotContains(t, schema.Required, "updatedAt")
}

func TestSchemaBuilderBuildSchema_AllBasicTypes(t *testing.T) {
	tests := []struct {
		typeStr  string
		wantType string
		wantFmt  string
	}{
		{"string", "string", ""},
		{"int", "integer", "int64"},
		{"int8", "integer", "int64"},
		{"int16", "integer", "int64"},
		{"int32", "integer", "int64"},
		{"int64", "integer", "int64"},
		{"uint", "integer", "int64"},
		{"uint8", "integer", "int64"},
		{"uint16", "integer", "int64"},
		{"uint32", "integer", "int64"},
		{"uint64", "integer", "int64"},
		{"float32", "number", "float"},
		{"float64", "number", "double"},
		{"bool", "boolean", ""},
		{"byte", "string", "byte"},
		{"rune", "integer", "int32"},
		{"time.Time", "string", "date-time"},
		{"time.Duration", "string", ""},
		{"interface{}", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.typeStr, func(t *testing.T) {
			sb := NewSchemaBuilder()
			schema := sb.BuildSchema(tt.typeStr)

			assert.NotNil(t, schema)
			assert.Equal(t, tt.wantType, schema.Type)
			if tt.wantFmt != "" {
				assert.Equal(t, tt.wantFmt, schema.Format)
			}
		})
	}
}
