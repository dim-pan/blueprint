package parser

import "fmt"

type Priority string

const (
	PriorityMustHave   Priority = "must-have"
	PriorityShouldHave Priority = "should-have"
	PriorityNiceToHave Priority = "nice-to-have"
)

type Pattern string

const (
	PatternUbiquitous  Pattern = "ubiquitous"
	PatternEventDriven Pattern = "event-driven"
	PatternUnwanted    Pattern = "unwanted"
	PatternStateDriven Pattern = "state-driven"
	PatternOptional    Pattern = "optional"
)

type AcceptanceCriterion struct {
	Description string
}

type Requirement struct {
	ID           string
	Title        string
	Priority     Priority
	Pattern      Pattern
	Statement    string
	Accept       []AcceptanceCriterion
	DerivedFrom  string
	AllocatedTo  string
	SourceFile   string
	LineNumber   int
}

type SystemModel struct {
	Requirements []Requirement
	Components   []Component
	Interfaces   []InterfaceDefinition
	TestSpecs    []TestSpec
}

type TestSpec struct {
	ID         string
	Title      string
	Verifies   []string
	Given      string
	When       string
	Expect     string
	SourceFile string
	LineNumber int
}

type Component struct {
	ID             string
	Name           string
	Responsibility string
	DependsOn      []string
	Satisfies      []string
	SourceFile     string
	LineNumber     int
}

type Field struct {
	Name     string
	Type     string
	Optional bool
}

type InterfaceDefinition struct {
	Name       string
	Fields     []Field
	SourceFile string
	LineNumber int
}

type ParseError struct {
	File    string
	Line    int
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Message)
}
