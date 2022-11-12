# Constructor annotations

# TODO:
- [X] Support type embedding.
- [ ] Builder should not mutate internal structure and prepare new one on method build().
- [ ] At the moment constructor arguments changes arguments order. Must be fixed.
- [ ] Introduce a parameter that forces printing `TypeName`/`FieldName` in lower or upper case for templates.

## Usage

To start using the annotations developer needs to add corresponding annotation processor to the gen tool:
```go
import (
	_ "github.com/YReshetko/go-annotation/annotations/constructor"
	"github.com/YReshetko/go-annotation/pkg"
)

func main() {
	annotation.Process()
}
```

When developer applies one of annotations described below the annotation processor puts `constructor.gen.go` file in the same package where the annotations are used. All generated code will be exposed in single file for a package.


## Exposed annotations
- [Constructor](#constructor)
- [Optional](#optional)
- [Builder](#builder)
- [PostConstruct](#postconstruct)
- [Init](#init)
- [Exclude](#exclude)

### Constructor

Can be declared on structures only. When annotation processor receives the annotation with structure it generates a function with arguments of structure fields and returns new prefilled structure.

**Example:**
```go
// @Constructor
type Some struct {
    fieldA int
    fieldB string
}
```
The next function will be generated:
```go
func NewSome (fieldA int, fieldB string) Some {
    returnValue := Some{
        fieldA: fieldA,
        fieldB: fieldB,
    }
    return returnValue
}
```

Developer can override function name and return type (pointer/structure) by annotation parameters:

**name** - overrides function name. Here can be added explicit name of function for example `SomeConstructor` or template with static argument `TypeName`: `{{ .TypeName }}Constructor`. In both cases the generated function name will be `SomeConstructor`.

**type** - receives one of two constants: `pointer`, `struct`. By default (if developer doesn't write the parameter) the constructor will use `struct` value. If `pointer` is used the function is generated with pointer return type: `func NewSome (fieldA int, fieldB string) *Some...`

**Full annotation syntax:**
```go
// @Constructor(name="{{ .TypeName }}Constructor", type="pointer")
type Some struct {
    fieldA int
    fieldB int
}
```

### Optional

Can be declared on structures only. When annotation processor receives the annotation with structure it generates an optional function type, batch of functions that receive structure field type and a function that receives options and returns structure type.

**Example:**
```go
// @Optional
type Some struct {
    fieldA int
    fieldB int
}
```
The next code will be generated:
```go
type SomeOption func(*Some)

func NewSome(opts ...SomeOption) Some {
    rt := &Some{}
    for _, o := range opts {
        o(rt)
    }
    return *rt
}

func WithFieldA(v int) SomeOption {
    return func(rt *Some) {
        rt.fieldA = v
    }
}

func WithFieldB(v string) SomeOption {
    return func(rt *Some) {
        rt.fieldB = v
    }
}
```

Developer can override optional type name, `with...` function names, function name that builds a structure and return type (pointer/structure) by annotation parameters:

**name** - optional type name. Template with `TypeName` can be used here. For example: `name=My{{ .TypeName }}Option` then next optional type is generated: `type MySomeOption func(*Some)`.

**constructor** - overrides constructor name. Here can be added explicit name of function for example `SomeConstructor` or template with static argument `TypeName`: `{{ .TypeName }}Constructor`. In both cases the generated function name will be `SomeConstructor`.

**with** - overrides with functions pattern name. **_Developer should not use explicit function name as it caused a panic._** `FieldName` must be used as template parameter. For example: `SomeWith{{.FieldName}}` then the next functions are generated for the above example: `func SomeWithFieldA(v int) SomeOption` and `func SomeWithFieldB(v string) SomeOption`

**type** - receives one of two constants: `pointer`, `struct`. By default (if developer doesn't write the parameter) the constructor will use `struct` value. If `pointer` is used the function is generated with pointer return type: `func NewSome (fieldA int, fieldB string) *Some...`

**Full annotation syntax:**
```go
// @Optional(name="My{{ .TypeName }}Option", constructor="{{ .TypeName }}Constructor", with="SomeWith{{ .FieldName }}", type="pointer")
type Some struct {
    fieldA int
    fieldB string
}
```

### Builder

Can be declared on structures only. When annotation processor receives the annotation with structure it generates a builder type with batch of methods that sets structure parameters and method `build()` that returns new prefilled structure.

**Example:**
```go
// @Builder
type Some struct {
    fieldA int
    fieldB string
}
```
The next code will be generated:
```go

type SomeBuilder struct {
    value Some
}

func NewSomeBuilder() *SomeBuilder {
    return &SomeBuilder{}
}

func (b *SomeBuilder) FieldA(v int) *SomeBuilder {
    b.value.fieldA = v
    return b
}

func (b *SomeBuilder) FieldB(v string) *SomeBuilder {
    b.value.fieldB = v
    return b
}

func (b *SomeBuilder) Build() Some {
    return b.value
}
```

Developer can override builder type name, setter method names, builder constructor name, method name that builds a structure and return type (pointer/structure) by annotation parameters:

**name** - builder type name. Template with `TypeName` can be used here. For example: `name=My{{ .TypeName }}Builder` then next builder type is generated: `type MySomeBuilder struct `.

**constructor** - overrides builder constructor name. Here can be added explicit name of function for example `SomeBuilderConstructor` or template with static argument `TypeName`: `{{ .TypeName }}BuilderConstructor`. In both cases the generated function name will be `SomeBuilderConstructor`.

**build** - overrides setter methods pattern name. **_Developer should not use explicit method name as it caused a panic_**. `FieldName` must be used as template parameter. For example: `Set{{.FieldName}}` then the next functions are generated for the above example: `func (b *SomeBuilder) SetFieldA(v int) *SomeBuilder` and `func (b *SomeBuilder) SetFieldB(v string) *SomeBuilder`

**terminator** - the `Build()` method name can be overridden by the parameter. Can not be used any templates here the method name must be specified explicitly.

**type** - receives one of two constants: `pointer`, `struct`. By default (if developer doesn't write the parameter) the constructor will use `struct` value. If `pointer` is used the function is generated with pointer return type: `func NewSome (fieldA int, fieldB string) *Some...`

**Full annotation syntax:**
```go
// @Builder(name="{{.TypeName}}Builder", constructor="{{.TypeName}}BuilderConstructor", build="Set{{.FieldName}}", terminator="MyBuildMethod", type="pointer")
type Some struct {
    fieldA int
    fieldB string
}
```

### PostConstruct

The annotation can be used on methods without arguments of a structure that has any of above annotations (`@Constructor`, `@Optional`, `@Builder`). The methods that marked by `PostConstruct` annotation are called before returning new structure.

**Example:**
```go
// @Constructor, @Optional(constructor="newSome"), @Builder
type Some struct {
    fieldA int
    fieldB string
}

// @PostConstruct
func (s Some)print()  {
    fmt.Println(s.fieldA, s.fieldB)
}
```
The next code will be generated:
```go

func NewSome(fieldA int, fieldB string) Some {
    ...
    returnValue.print()
    return returnValue
}
...
func newSome(opts ...SomeOption) Some {
    ...
    rt.print()
    return *rt
}
...
func (b *SomeBuilder) Build() Some {
    b.value.print()
    return b.value
}
```

The annotation has only one parameter **priority** that defines an order of the methods call. The first method will be called with the lowest priority value. By default `priority=1`.
**Full annotation syntax:**
```go
// @PostConstruct(priority="10")
func (s Some) print() {
    fmt.Println(s.fieldA, s.fieldB)
}
```

### Init

Sometimes developer doesn't need to explicitly set slices, maps and channels, but it should be initialised on structure creation, so the annotation helps to avoid nil references on slices, maps and channels.


**Example:**
```go
// @Constructor, @Optional(constructor="newSome"), @Builder
type Some struct {
    field []string //@Init
}
```
The next code will be generated:
```go

func NewSome() Some {
    returnValue := Some{
        field: []string{},
    }
    return returnValue
}
...
func newSome(opts ...SomeOption) Some {
    rt := &Some{
        field: []string{},
    }
    for _, o := range opts {
        o(rt)
    }
    return *rt
}
...
func (b *SomeBuilder) Build() Some {
    if b.value.field == nil {
        b.value.field = []string{}
    }
    return b.value
}
```

As you can see the initialisation is a bit different in all constructors, so that have next effects:
- `@Constructor` - will not be generated with parameters that declared with fields that marked by `@Init` annotation
- `@Optional` - the initialisation is done prior options are applied, so it cn be overridden by corresponding option.
- `@Builder` - checks if the field is already initialised and if not it will do that work.

The `@Init` annotation has two parameters that helps to set up initial length and capacity for the slices, map and channels:

**len** - applicable for slices only. Sets initial length

**cap** - applicable for slices, maps and channels, Sets initial capacity

**Full annotation syntax:**
```go
// @Constructor, @Optional(constructor="newSome"), @Builder
type Some struct {
    field []string //@Init(len="10", cap="20")
}
```

### Exclude

This annotation should be declared on those fields that should not be set externally.

**Example:**
```go
// @Constructor, @Optional(constructor="newSome"), @Builder
type Some struct {
    field []string //@Exclude
}
```
The next code will be generated:
```go
func NewSome() Some {
    returnValue := Some{}
    return returnValue
}

type SomeOption func(*Some)

func newSome(opts ...SomeOption) Some {
    rt := &Some{}
    for _, o := range opts {
        o(rt)
	}
    return *rt
}

type SomeBuilder struct {
    value Some
}

func NewSomeBuilder() *SomeBuilder {
    return &SomeBuilder{}
}

func (b *SomeBuilder) Build() Some {
    return b.value
}
```

As you can see there is no possibility to set `field`. This annotation doesn't have any parameters.