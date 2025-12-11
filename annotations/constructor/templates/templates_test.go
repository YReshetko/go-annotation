package templates_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/YReshetko/go-annotation/annotations/constructor/generators"
	"github.com/YReshetko/go-annotation/annotations/constructor/templates"
)

func TestBuildConstructor_PostConstruct_WithError(t *testing.T) {
	data := generators.ConstructorValues{
		FunctionName: "MyConstructor",
		Arguments:    []string{"a int", "b bool"},
		ReturnType:   "My",
		Fields: []struct {
			Name  string
			Value string
		}{
			{Name: "a", Value: "a"},
			{Name: "b", Value: "b"},
		},
		PostConstructs: []generators.PostConstructValues{
			{MethodName: "post1", ReturnsError: false},
			{MethodName: "post2", ReturnsError: true},
			{MethodName: "post3", ReturnsError: false},
		},
		ReturnsError: true,
	}

	res, err := templates.Execute(templates.ConstructorTpl, data)
	require.NoError(t, err)
	fmt.Println(string(res))
}

func TestBuildConstructor_PostConstruct_NoError(t *testing.T) {
	data := generators.ConstructorValues{
		FunctionName: "MyConstructor",
		Arguments:    []string{"a int", "b bool"},
		ReturnType:   "My",
		Fields: []struct {
			Name  string
			Value string
		}{
			{Name: "a", Value: "a"},
			{Name: "b", Value: "b"},
		},
		PostConstructs: []generators.PostConstructValues{
			{MethodName: "post1", ReturnsError: false},
			{MethodName: "post2", ReturnsError: false},
			{MethodName: "post3", ReturnsError: false},
		},
		ReturnsError: false,
	}

	res, err := templates.Execute(templates.ConstructorTpl, data)
	require.NoError(t, err)
	fmt.Println(string(res))
}

func TestBuildConstructor_NoPostConstruct(t *testing.T) {
	data := generators.ConstructorValues{
		FunctionName: "MyConstructor",
		Arguments:    []string{"a int", "b bool"},
		ReturnType:   "My",
		Fields: []struct {
			Name  string
			Value string
		}{
			{Name: "a", Value: "a"},
			{Name: "b", Value: "b"},
		},
	}

	res, err := templates.Execute(templates.ConstructorTpl, data)
	require.NoError(t, err)
	fmt.Println(string(res))
}

func TestBuildOptional_PostConstruct_WithError(t *testing.T) {
	data := generators.OptionalValues{
		FunctionName:     "MyOptional",
		OptionalTypeName: "MyOpt",
		ReturnType:       "My",
		PostConstructs: []generators.PostConstructValues{
			{MethodName: "post1", ReturnsError: false},
			{MethodName: "post2", ReturnsError: true},
			{MethodName: "post3", ReturnsError: false},
		},
		ReturnsError: true,
	}

	res, err := templates.Execute(templates.OptionalConstructorTpl, data)
	require.NoError(t, err)
	fmt.Println(string(res))
}

func TestBuildOptional_PostConstruct_NoError(t *testing.T) {
	data := generators.OptionalValues{
		FunctionName:     "MyOptional",
		OptionalTypeName: "MyOpt",
		ReturnType:       "My",
		PostConstructs: []generators.PostConstructValues{
			{MethodName: "post1", ReturnsError: false},
			{MethodName: "post2", ReturnsError: false},
			{MethodName: "post3", ReturnsError: false},
		},
	}

	res, err := templates.Execute(templates.OptionalConstructorTpl, data)
	require.NoError(t, err)
	fmt.Println(string(res))
}

func TestBuildOptional_NoPostConstruct(t *testing.T) {
	data := generators.OptionalValues{
		FunctionName:     "MyOptional",
		OptionalTypeName: "MyOpt",
		ReturnType:       "My",
	}

	res, err := templates.Execute(templates.OptionalConstructorTpl, data)
	require.NoError(t, err)
	fmt.Println(string(res))
}

func TestBuildBuilder_PostConstruct_WithError(t *testing.T) {
	data := generators.BuilderValues{
		BuilderTypeName: "MyBuilder",
		ConstructorName: "MyBuilderConstructor",
		BuildMethodName: "Build",
		ReturnType:      "My",
		PostConstructs: []generators.PostConstructValues{
			{MethodName: "post1", ReturnsError: false},
			{MethodName: "post2", ReturnsError: true},
			{MethodName: "post3", ReturnsError: false},
		},
		ReturnsError: true,
	}

	res, err := templates.Execute(templates.BuilderBuildMethodTpl, data)
	require.NoError(t, err)
	fmt.Println(string(res))
}

func TestBuildBuilder_PostConstruct_NoError(t *testing.T) {
	data := generators.BuilderValues{
		BuilderTypeName: "MyBuilder",
		ConstructorName: "MyBuilderConstructor",
		BuildMethodName: "Build",
		ReturnType:      "My",
		PostConstructs: []generators.PostConstructValues{
			{MethodName: "post1", ReturnsError: false},
			{MethodName: "post2", ReturnsError: false},
			{MethodName: "post3", ReturnsError: false},
		},
	}

	res, err := templates.Execute(templates.BuilderBuildMethodTpl, data)
	require.NoError(t, err)
	fmt.Println(string(res))
}

func TestBuildBuilder_NoPostConstruct(t *testing.T) {
	data := generators.BuilderValues{
		BuilderTypeName: "MyBuilder",
		ConstructorName: "MyBuilderConstructor",
		BuildMethodName: "Build",
		ReturnType:      "My",
	}

	res, err := templates.Execute(templates.BuilderBuildMethodTpl, data)
	require.NoError(t, err)
	fmt.Println(string(res))
}
