package http

import (
	"errors"

	"github.com/bells307/gitlab-watcher/internal/model"
	"github.com/bells307/gitlab-watcher/pkg/gin/err_resp"
	"github.com/gin-gonic/gin"
)

const GitlabEventHeaderName string = "X-Gitlab-Event"

// Строковые значения, передаваемые в хедере
const (
	MergeRequestEventHeader string = "Merge Request Hook"
	PipelineEventHeader     string = "Pipeline Hook"
)

type HookHandler struct {
	hookServices []HookService
}

// Сервис обработки хуков от гитлаба
type HookService interface {
	ProcessMergeRequestHook(model.MergeRequest) error
	ProcessPipelineHook(model.Pipeline) error
}

func NewHookHandler(hookServices []HookService) *HookHandler {
	return &HookHandler{hookServices}
}

func (h *HookHandler) Register(router *gin.Engine) {
	router.POST("/api/hook", h.processHook)
}

func (h *HookHandler) processHook(c *gin.Context) {
	header, ok := c.Request.Header[GitlabEventHeaderName]
	if !ok {
		err_resp.NewErrorResponse(c, errors.New("gitlab event header not provided"))
		return
	}

	switch header[0] {
	case MergeRequestEventHeader:
		var input model.MergeRequest
		if err := c.Bind(&input); err != nil {
			return
		}

		for _, s := range h.hookServices {
			if err := s.ProcessMergeRequestHook(input); err != nil {
				err_resp.NewErrorResponse(c, err)
				return
			}
		}

	case PipelineEventHeader:
		var input model.Pipeline
		if err := c.Bind(&input); err != nil {
			return
		}

		for _, s := range h.hookServices {
			if err := s.ProcessPipelineHook(input); err != nil {
				err_resp.NewErrorResponse(c, err)
				return
			}
		}

	default:
		err_resp.NewErrorResponse(c, errors.New("unknown event type header"))
		return
	}
}
