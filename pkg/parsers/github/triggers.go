package github

import (
	"errors"

	"github.com/argonsecurity/pipeline-parser/pkg/models"
	githubModels "github.com/argonsecurity/pipeline-parser/pkg/parsers/github/models"
	"github.com/argonsecurity/pipeline-parser/pkg/utils"
	"github.com/mitchellh/mapstructure"
)

const (
	pushEvent              = "push"
	forkEvent              = "fork"
	workflowDispatchEvent  = "workflow_dispatch"
	pullRequestEvent       = "pull_request"
	pullRequestTargetEvent = "pull_request_target"
)

var (
	githubEventToModelEvent = map[string]models.EventType{
		pushEvent:             models.PushEvent,
		forkEvent:             models.ForkEvent,
		workflowDispatchEvent: models.ManualEvent,
		pullRequestEvent:      models.PullRequestEvent,
	}
)

func parseWorkflowTriggers(workflow *githubModels.Workflow) ([]models.Trigger, error) {
	triggers := []models.Trigger{}
	if workflow.On == nil {
		return nil, nil
	}

	var on githubModels.On
	if events, ok := utils.ToSlice[string](workflow.On); ok {
		triggers = generateTriggersFromEvents(events)
	} else if err := mapstructure.Decode(workflow.On, &on); err == nil {
		events := utils.Filter(utils.GetMapKeys(on.Events), func(event string) bool {
			_, ok := githubEventToModelEvent[event]
			return !ok
		})
		triggers = append(triggers, generateTriggersFromEvents(events)...)
		if on.Push != nil {
			triggers = append(triggers, parseGitEvent(on.Push, models.PushEvent))
		}

		if on.PullRequest != nil {
			triggers = append(triggers, parseGitEvent(on.PullRequest, models.PullRequestEvent))
		}

		if on.PullRequestTarget != nil {
			triggers = append(triggers, parseGitEvent(on.PullRequestTarget, models.EventType(pullRequestTargetEvent)))
		}

		if on.WorkflowDispatch != nil {
			triggers = append(triggers, parseWorkflowDispatch(on.WorkflowDispatch))
		}

		if on.WorkflowCall != nil {
			triggers = append(triggers, parseWorkflowCall(on.WorkflowCall))
		}
	} else {
		return nil, errors.New("failed to parse workflow triggers")
	}
	return triggers, nil
}

func parseWorkflowCall(workflowCall *githubModels.WorkflowCall) models.Trigger {
	return models.Trigger{
		Event:     models.PipelineTriggerEvent,
		Paramters: parseInputs(workflowCall.Inputs),
	}
}

func parseInputs(inputs githubModels.Inputs) []models.Parameter {
	parameters := []models.Parameter{}
	if inputs != nil {
		for k, v := range inputs {
			parameters = append(parameters, models.Parameter{
				Name:        &k,
				Description: &v.Description,
				Default:     v.Default,
			})
		}
	}
	return parameters
}

func parseWorkflowDispatch(workflowDispatch *githubModels.WorkflowDispatch) models.Trigger {
	return models.Trigger{
		Event:     models.ManualEvent,
		Paramters: parseInputs(workflowDispatch.Inputs),
	}
}

func parseGitEvent(gitevent *githubModels.Gitevent, event models.EventType) models.Trigger {
	trigger := models.Trigger{
		Event: event,
		Paths: &models.Filter{
			AllowList: []string{},
			DenyList:  []string{},
		},
		Branches: &models.Filter{
			AllowList: []string{},
			DenyList:  []string{},
		},
	}

	for _, path := range gitevent.Paths {
		trigger.Paths.AllowList = append(trigger.Paths.AllowList, path)
	}
	for _, path := range gitevent.PathsIgnore {
		trigger.Paths.DenyList = append(trigger.Paths.DenyList, path)
	}
	for _, branch := range gitevent.Branches {
		trigger.Branches.AllowList = append(trigger.Branches.AllowList, branch)
	}
	for _, branch := range gitevent.BranchesIgnore {
		trigger.Branches.DenyList = append(trigger.Branches.DenyList, branch)
	}

	return trigger
}

func generateTriggersFromEvents(events []string) []models.Trigger {
	return utils.Map(events, func(event string) models.Trigger {
		modelEvent, ok := githubEventToModelEvent[event]
		if !ok {
			modelEvent = models.EventType(event)
		}
		return models.Trigger{
			Event: modelEvent,
		}
	})
}
