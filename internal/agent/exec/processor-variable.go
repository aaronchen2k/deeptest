package agentExec

import logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"

type ProcessorVariable struct {
	ID uint `json:"id" yaml:"id"`
	ProcessorEntity

	VariableName string `json:"variableName" yaml:"variableName"`
	RightValue   string `json:"rightValue" yaml:"rightValue"`
}

func (p ProcessorVariable) Run(s *Session) (ret *Log, err error) {
	logUtils.Infof("variable")
	return
}
