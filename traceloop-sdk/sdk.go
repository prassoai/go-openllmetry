package traceloop

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	apitrace "go.opentelemetry.io/otel/trace"

	semconvai "github.com/prassoai/go-openllmetry/semconv-ai"
)

type Traceloop struct {
	Tracer apitrace.Tracer
}

type LLMSpan struct {
	span apitrace.Span
}

func NewClient(tracer apitrace.Tracer) (*Traceloop, error) {
	instance := Traceloop{
		Tracer: tracer,
	}
	return &instance, nil
}

func setMessagesAttribute(span apitrace.Span, prefix string, messages []Message) {
	for _, message := range messages {
		attrsPrefix := fmt.Sprintf("%s.%d", prefix, message.Index)
		span.SetAttributes(
			attribute.String(attrsPrefix+".content", message.Content),
			attribute.String(attrsPrefix+".role", message.Role),
		)
	}
}

func (instance *Traceloop) LogPrompt(ctx context.Context, prompt Prompt, workflowAttrs WorkflowAttributes) (LLMSpan, error) {
	spanName := fmt.Sprintf("%s.%s", prompt.Vendor, prompt.Mode)
	_, span := instance.Tracer.Start(ctx, spanName)

	span.SetAttributes(
		semconvai.LLMVendor.String(prompt.Vendor),
		semconvai.LLMRequestModel.String(prompt.Model),
		semconvai.LLMRequestType.String(prompt.Mode),
		semconvai.TraceloopWorkflowName.String(workflowAttrs.Name),
	)

	setMessagesAttribute(span, "llm.prompts", prompt.Messages)

	return LLMSpan{
		span: span,
	}, nil
}

func (llmSpan *LLMSpan) LogCompletion(ctx context.Context, completion Completion, usage Usage) error {
	llmSpan.span.SetAttributes(
		semconvai.LLMResponseModel.String(completion.Model),
		semconvai.LLMUsageTotalTokens.Int(usage.TotalTokens),
		semconvai.LLMUsageCompletionTokens.Int(usage.CompletionTokens),
		semconvai.LLMUsagePromptTokens.Int(usage.PromptTokens),
	)

	setMessagesAttribute(llmSpan.span, "llm.completions", completion.Messages)

	defer llmSpan.span.End()

	return nil
}
