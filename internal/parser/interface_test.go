package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/YReshetko/go-annotation/internal/parser"
)

func TestParserSuccess(t *testing.T) {
	test := []struct {
		name   string
		text   string
		verify func(*testing.T, []parser.Annotation)
	}{
		{"empty", "",
			func(t *testing.T, annotations []parser.Annotation) {
				assert.Empty(t, annotations)
			},
		},
		{"no annotations", "Some text with no annotations",
			func(t *testing.T, annotations []parser.Annotation) {
				assert.Empty(t, annotations)
			},
		},
		{"false positive annotation syntax", "@Positive)",
			func(t *testing.T, annotations []parser.Annotation) {
				require.Len(t, annotations, 1)
				assert.Equal(t, "Positive", annotations[0].Name())
				assert.Empty(t, annotations[0].Params())
			},
		},
		{"single line annotation with no params", "// @SomeAnnotation",
			func(t *testing.T, annotations []parser.Annotation) {
				require.Len(t, annotations, 1)
				assert.Equal(t, "SomeAnnotation", annotations[0].Name())
				assert.Empty(t, annotations[0].Params())
			},
		},
		{"single line annotation with empty param", `// @SomeAnnotation(name="")`,
			func(t *testing.T, annotations []parser.Annotation) {
				require.Len(t, annotations, 1)
				assert.Equal(t, "SomeAnnotation", annotations[0].Name())
				assert.Len(t, annotations[0].Params(), 1)
				assert.Empty(t, annotations[0].Params()["name"])
			},
		},
		{"single line annotation with params", `// @RestAnnotation(method="GET", endpoint="/api/v1/rest")`,
			func(t *testing.T, annotations []parser.Annotation) {
				require.Len(t, annotations, 1)
				assert.Equal(t, "RestAnnotation", annotations[0].Name())
				require.Len(t, annotations[0].Params(), 2)

				v, ok := annotations[0].Params()["method"]
				require.True(t, ok)
				assert.Equal(t, "GET", v)

				v, ok = annotations[0].Params()["endpoint"]
				require.True(t, ok)
				assert.Equal(t, "/api/v1/rest", v)
			},
		},
		{"multi line annotation with params", `
					Some test
					@RestAnnotation(method="GET", endpoint="/api/v1/rest")
					Surrounds the annotation`,
			func(t *testing.T, annotations []parser.Annotation) {
				require.Len(t, annotations, 1)
				assert.Equal(t, "RestAnnotation", annotations[0].Name())
				require.Len(t, annotations[0].Params(), 2)

				v, ok := annotations[0].Params()["method"]
				require.True(t, ok)
				assert.Equal(t, "GET", v)

				v, ok = annotations[0].Params()["endpoint"]
				require.True(t, ok)
				assert.Equal(t, "/api/v1/rest", v)
			},
		},
		{"multi line multi annotation with params", `
					Some test
					@RestAnnotation(method="GET", endpoint="/api/v1/rest")
					@RestAnnotation
					(
							method="GET", 
							endpoint="/api/v1/rest"
					)
					@RestAnnotation(method="GET", endpoint="/api/v1/rest")
					Surrounds the annotation`,
			func(t *testing.T, annotations []parser.Annotation) {
				require.Len(t, annotations, 3)
				for _, annotation := range annotations {
					assert.Equal(t, "RestAnnotation", annotation.Name())
					require.Len(t, annotation.Params(), 2)

					v, ok := annotation.Params()["method"]
					require.True(t, ok)
					assert.Equal(t, "GET", v)

					v, ok = annotation.Params()["endpoint"]
					require.True(t, ok)
					assert.Equal(t, "/api/v1/rest", v)
				}

			},
		},
	}

	for _, s := range test {
		t.Run(s.name, func(t *testing.T) {
			annotations, err := parser.Parse(s.text)
			require.Nil(t, err)
			s.verify(t, annotations)
		})
	}
}

func TestParserFailure(t *testing.T) {
	test := []struct {
		name   string
		text   string
		errors []string
	}{
		{
			name: "No closed bracket", text: "@Invalid(",
			errors: []string{
				"expected token IDENT instead got EOF",
			},
		},
		{
			name: "invalid arguments syntax", text: "@Invalid(djksadfjkgsdkjfgdsj)",
			errors: []string{
				"expected token = instead got )",
			},
		},
		{
			name: "invalid syntax with no closed bracket", text: "@Invalid(djksadfjkgsdkjfgdsj",
			errors: []string{
				"expected token = instead got EOF",
			},
		},
		{
			name: "invalid non-closed quote", text: `@Invalid("name=val`,
			errors: []string{
				"expected token IDENT instead got STRING",
			},
		},
		{
			name: "invalid non-closed quote", text: `@Invalid("name"="val, "hello"="world")`,
			errors: []string{
				"expected token IDENT instead got STRING",
				"expected token = instead got STRING",
			},
		},
		{
			name: "multiline invalid syntax", text: `
					// Some text
					@Invalid(
						"djksadfjkgsdkjfgdsj" some text = "GET"
					) // Some text`,
			errors: []string{
				"expected token IDENT instead got STRING",
			},
		},
		{name: "multi line with invalid format params ", text: `
					Some test
					@RestAnnotation(method="GET", endpoint="/api/v1/rest")
					@RestAnnotation(method="GET",endpoint="/api/v1/rest)
					@RestAnnotation(method="GET", endpoint="/api/v1/rest")
					Surrounds the annotation`,
			errors: []string{
				"expected token = instead got STRING",
			},
		},
	}
	for _, s := range test {
		t.Run(s.name, func(t *testing.T) {
			annotations, err := parser.Parse(s.text)
			require.NotNil(t, err)
			for _, e := range s.errors {
				assert.Contains(t, err.Error(), e)
			}
			assert.Nil(t, annotations)
		})
	}
}
