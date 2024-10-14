# Cobra annotations (*experimental*)

---
GO 1.18+
---

## Usage

To start using the annotations developer needs to add corresponding annotation processor to the gen tool:
```go
import (
	_ "github.com/YReshetko/go-annotation/annotations/cobra"
	"github.com/YReshetko/go-annotation/pkg"
)

func main() {
	annotation.Process()
}
```

When developer applies one of annotations described below the annotation processor puts a number of files within main package that is marked with `@CobraOutput` annotation.
The CLI program entry point is also generated within the annotation processor.


## Exposed annotations
- [CobraOutput](#output)
- [Cobra](#cobra)
- [CobraRun](#run)
- [CobraPersistPreRun](#persistprerun)
- [CobraPreRun](#prerun)
- [CobraPostRun](#postrun)
- [CobraPersistPostRun](#persistpostrun)
  
  
  
  
### Output

Can be declared on any package, but the output code is printed with `package main` 

**Example:**
```go
// @CobraOutput
package main

```
The file can be used to add go doc for your CLI program

### Cobra
The `@Cobra` annotation is designed to provide all required information to setup cobra command. It should mark a structure declaration only, otherwise the processor returns an error.
The annotation has next parameters:

**build** - parameter defines a number of build tags where output can be placed. With the functionality we can build different set of commands within the same codebase.
For example, if build parameter is omitted, all commands are printed into the same `default.gen.go` file with no build tags.
If the parameter has one or more values, separated by comma (eg. `build="private,public"`), processor generates 2 files: `private.gen.go` and `public.gen.go`. 
So, you can have different set of commands, that a marked with `public` go to `public.gen.go`, `private` to `private.gen.go`. To build desired setup you would need to use `--tags=public` or `--tags==private`

**usage** - similar to `cmd.Use`, but instead the parameter should have a sequence of command words, separated by whitespace. For example, root command must have one word (in most cases that should be a binary file name), child for a root command must have 2 words and so on, for more details see examples folder.
The parameter is extremely important, because it is used to build commands tree.

**example** - used to fill `cmd.Example` (can be omitted)

**short** - used to fill `cmd.Short` (can be omitted)

**long** - used to fill `cmd.Long` (can be omitted)

**silenceUsage** - used to fill `cmd.SilenceUsage` (can be omitted)

**silenceError** - used to fill `cmd.SilenceError` (can be omitted)

The structure that is marked by `@Cobra` annotation can have non-pointer fields that are marked by flag tags:

**flag** - contains comma separated parameters that define a flag name on the first position. 
Then if there is:
- `required` - the flag is marked as required for cobra command
- `persist` - the flag is registered for as a persistent
- `inherited` - the flag will be taken from parent command

**short** - the one letter flag shortcut

**default** - flag default value

**description** - flag description

Example:
```go
// @Cobra(
//  build = "default",
//  usage = "cli",
//  example = "cli [-F file | -D dir] ... [-f format] profile",
//  short = "Root command of the application (short)",
//  long = "Root command of the application (long)",
//
// )
type RootCommand struct {
    Output string `flag:"output,required" short:"o" description:"output file name"`
    Num    int    `flag:"num" default:"42" description:"some number for command"`
    IsOK   bool   `flag:"is-ok,persist" short:"i" default:"true" description:"some persistent flag"`
}

// @Cobra(
//  build = "default",
//  usage = "cli get",
//  example = "cli get [-F format] resource",
//  short = "Get command of the application (short)",
//  long = "Get command of the application (long)",
//
// )
type GetCommand struct {
    Dur   time.Duration `flag:"dur,persist" short:"d" default:"12s" description:"Duration flag description"`
}
```

The flags are parsed for primitives, apart from complex and uintptr. Strings and time.Duration are supported put of box as well.
If you need to support any other type of flags you can implement following interface:
```go
type flagsMarshaller interface {
    MarshalFlag(string) error
}
```
The functionality can help to support enums and many other types. Some examples you can see into examples folder

### Run
The `@CobraRun` is used to mark an exported method of the structure described above. So, the method is called in `cmd.RunE` just after flags parsing.
The method receiver should be non-pointer. The method should have one of following signatures:
```go
Run(cmd *cobra.Command, agrs []string) error
Run(cmd *cobra.Command, agrs []string)
```
Method name can be different from `Run`

### PersistPreRun
The `@CobraPersistPreRun` annotation is used similar as above but for `cmd.PersistentPreRunE`

### PreRun
The `@CobraPreRun` annotation is used similar as above but for `cmd.PreRunE`

### PostRun
The `@CobraPostRun` annotation is used similar as above but for `cmd.PostRunE`

### PersistPostRun
The `@CobraPersistPostRun` annotation is used similar as above but for `cmd.PersistentPostRunE`

# TODO:
- [ ] Support flag groups (all required, one of, any of) 
- [ ] Inject CobraContext into commands structures, to linc parent and child command and use it to prefill some data into PreRun functions
