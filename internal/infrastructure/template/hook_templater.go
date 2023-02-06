package template

import "github.com/bells307/gitlab-watcher/internal/model"

type HookTemplater interface {
	GetMergeRequestMessage(mr model.MergeRequest) (string, error)
	GetPipelineMessage(pipeline model.Pipeline) (string, error)
}
