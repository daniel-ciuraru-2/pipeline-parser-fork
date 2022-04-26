package github

import (
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/loaders/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
)

var (
	defaultTimeoutMS int = 360 * 60 * 1000
)

func parseWorkflowJobs(workflow *githubModels.Workflow) (*[]models.Job, error) {
	jobs, err := utils.MapToSliceErr(workflow.Jobs.NormalJobs, parseJob)
	if err != nil {
		return nil, err
	}
	return &jobs, nil
}

func parseJob(jobName string, job *githubModels.Job) (models.Job, error) {
	parsedJob := models.Job{
		ID:                   job.ID,
		Name:                 &job.Name,
		ContinueOnError:      &job.ContinueOnError,
		EnvironmentVariables: job.Env,
	}

	if job.Name == "" {
		parsedJob.Name = job.ID
	}

	if job.TimeoutMinutes != nil && *job.TimeoutMinutes == 0 {
		timeout := int(*job.TimeoutMinutes) * 60 * 1000
		parsedJob.TimeoutMS = &timeout
	} else {
		parsedJob.TimeoutMS = &defaultTimeoutMS
	}

	if job.If != "" {
		parsedJob.Conditions = &[]models.Condition{models.Condition(job.If)}
	}

	if job.Concurrency != nil {
		parsedJob.ConcurrencyGroup = job.Concurrency.Group
	}

	if job.Steps != nil {
		parsedJob.Steps = parseJobSteps(job.Steps)
	}

	if job.RunsOn != nil {
		parsedJob.Runner = parseRunsOnToRunner(job.RunsOn)
	}

	if job.Needs != nil {
		parsedJob.Dependencies = (*[]string)(job.Needs)
	}

	if job.Permissions != nil {
		permissions, err := parseTokenPermissions(job.Permissions)
		if err != nil {
			return models.Job{}, err
		}
		parsedJob.TokenPermissions = permissions
	}

	return parsedJob, nil
}
