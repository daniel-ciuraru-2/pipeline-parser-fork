package github

import (
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/parsers/github/models"
	"gopkg.in/yaml.v3"
)

func Parse(data []byte) (*models.Pipeline, error) {
	workflow := &githubModels.Workflow{}
	if err := yaml.Unmarshal(data, workflow); err != nil {
		return nil, err
	}

	pipeline := &models.Pipeline{}
	triggers, err := parseWorkflowTriggers(workflow)
	if err != nil {
		return nil, err
	}

	jobs := parseWorkflowJobs(workflow)
	pipeline.Jobs = &jobs

	pipeline.Triggers = &triggers
	return pipeline, nil
}
