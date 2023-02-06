package application

import (
	"fmt"
	"github.com/bells307/gitlab-watcher/internal/infrastructure/template"

	"github.com/bells307/gitlab-watcher/internal/infrastructure/sender"
	"github.com/bells307/gitlab-watcher/internal/model"
)

// Сервис обработки хуков гитлаба
type HookService struct {
	sender             sender.Sender
	templater          template.HookTemplater
	sendTemplateErrors bool
}

func NewHookService(sender sender.Sender, templater template.HookTemplater, sendTemplateErrors bool) *HookService {
	return &HookService{sender, templater, sendTemplateErrors}
}

func (s *HookService) ProcessMergeRequestHook(mr model.MergeRequest) error {
	msg, err := s.templater.GetMergeRequestMessage(mr)
	if err != nil {
		if s.sendTemplateErrors {
			if err := s.sender.SendMessage(fmt.Sprintf("error preparing merge request message: %v", err)); err != nil {
				return err
			}
		}
		return fmt.Errorf("error preparing merge request message: %v", err)
	}

	if msg == "" {
		return nil
	}

	if err := s.sender.SendMessage(msg); err != nil {
		return err
	}

	return nil
}

func (s *HookService) ProcessPipelineHook(pipeline model.Pipeline) error {
	msg, err := s.templater.GetPipelineMessage(pipeline)
	if err != nil {
		if s.sendTemplateErrors {
			if err := s.sender.SendMessage(fmt.Sprintf("error preparing pipeline message: %v", err)); err != nil {
				return err
			}
		}
		return fmt.Errorf("error preparing pipeline message: %v", err)
	}

	if msg == "" {
		return nil
	}

	if err := s.sender.SendMessage(msg); err != nil {
		return err
	}

	return nil
}
