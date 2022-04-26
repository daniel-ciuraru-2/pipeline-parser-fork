package models

type Job struct {
	ID                   *string
	Name                 *string
	Steps                *[]Step
	ContinueOnError      *bool
	PreSteps             *[]Step
	PostSteps            *[]Step
	EnvironmentVariables *EnvironmentVariables
	Runner               *Runner
	Conditions           *[]Condition
	ConcurrencyGroup     *string
	Inputs               *[]Parameter
	TimeoutMS            *int
	Tags                 *[]string
	TokenPermissions     *map[string]Permission
	Dependencies         *[]string
	Metadata             Metadata
}
