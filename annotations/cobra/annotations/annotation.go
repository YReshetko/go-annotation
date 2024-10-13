package annotations

type CobraOutput struct{}

type Cobra struct {
	Build        string `annotation:"name=build,default=default"`
	Usage        string `annotation:"name=usage,required"`
	Example      string `annotation:"name=example"`
	Short        string `annotation:"name=short"`
	Long         string `annotation:"name=long"`
	SilenceUsage bool   `annotation:"name=silenceUsage,default=false"`
	SilenceError bool   `annotation:"name=silenceError,default=false"`
}

type CobraPersistPreRun struct{}
type CobraPreRun struct{}
type CobraRun struct{}
type CobraPostRun struct{}
type CobraPersistPostRun struct{}

func (a CobraPersistPreRun) IsPersistRun() bool { return true }
func (a CobraPersistPreRun) IsPreRun() bool     { return true }
func (a CobraPersistPreRun) IsPostRun() bool    { return false }

func (a CobraPreRun) IsPersistRun() bool { return false }
func (a CobraPreRun) IsPreRun() bool     { return true }
func (a CobraPreRun) IsPostRun() bool    { return false }

func (a CobraRun) IsPersistRun() bool { return false }
func (a CobraRun) IsPreRun() bool     { return false }
func (a CobraRun) IsPostRun() bool    { return false }

func (a CobraPostRun) IsPersistRun() bool { return false }
func (a CobraPostRun) IsPreRun() bool     { return false }
func (a CobraPostRun) IsPostRun() bool    { return true }

func (a CobraPersistPostRun) IsPersistRun() bool { return true }
func (a CobraPersistPostRun) IsPreRun() bool     { return false }
func (a CobraPersistPostRun) IsPostRun() bool    { return true }
