package template

import (
	"bytes"
	"fmt"
	"github.com/bells307/gitlab-watcher/internal/model"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	mergeRequestTemplateFile = "merge_request.txt"
	pipelineTemplateFile     = "pipeline.txt"
)

type HookTextTemplater struct {
	// Директория с шаблонами
	templatesDir string
	// Вспомогательные функции для шаблонов
	funcs map[string]any
}

func NewHookTextTemplater(templatesDir string) *HookTextTemplater {
	// Базовые функции работы со строками
	funcs := map[string]any{
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
	}

	return &HookTextTemplater{templatesDir, funcs}
}

// AddFunc Добавить функцию, которую можно будет использовать в текстовом шаблоне
func (t *HookTextTemplater) AddFunc(name string, fn any) {
	t.funcs[name] = fn
}

func (t *HookTextTemplater) GetMergeRequestMessage(mr model.MergeRequest) (string, error) {
	path := filepath.Join(t.templatesDir, mergeRequestTemplateFile)
	tpl, err := template.New(mergeRequestTemplateFile).Funcs(t.funcs).ParseFiles(path)
	if err != nil {
		return "", fmt.Errorf("error parsing merge request template file: %v", err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, mr); err != nil {
		return "", fmt.Errorf("error executing merge request template: %v", err)
	}

	return buf.String(), nil
}

func (t *HookTextTemplater) GetPipelineMessage(pipeline model.Pipeline) (string, error) {
	path := filepath.Join(t.templatesDir, pipelineTemplateFile)
	tpl, err := template.New(pipelineTemplateFile).Funcs(t.funcs).ParseFiles(path)
	if err != nil {
		return "", fmt.Errorf("error parsing pipeline template file: %v", err)
	}

	var buf bytes.Buffer
	if err := tpl.Funcs(t.funcs).Execute(&buf, pipeline); err != nil {
		return "", fmt.Errorf("error executing pipeline template: %v", err)
	}

	return buf.String(), nil
}
