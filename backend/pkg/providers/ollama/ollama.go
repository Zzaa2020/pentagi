package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"pentagi/pkg/config"
	"pentagi/pkg/providers/provider"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"gopkg.in/yaml.v3"
)

// AgentConfig represents the configuration for a single agent model
type AgentConfig struct {
	Model             string         `json:"model,omitempty" yaml:"model,omitempty"`
	MaxTokens         int            `json:"max_tokens,omitempty" yaml:"max_tokens,omitempty"`
	Temperature       float64        `json:"temperature,omitempty" yaml:"temperature,omitempty"`
	TopK              int            `json:"top_k,omitempty" yaml:"top_k,omitempty"`
	TopP              float64        `json:"top_p,omitempty" yaml:"top_p,omitempty"`
	N                 int            `json:"n,omitempty" yaml:"n,omitempty"`
	MinLength         int            `json:"min_length,omitempty" yaml:"min_length,omitempty"`
	MaxLength         int            `json:"max_length,omitempty" yaml:"max_length,omitempty"`
	RepetitionPenalty float64        `json:"repetition_penalty,omitempty" yaml:"repetition_penalty,omitempty"`
	FrequencyPenalty  float64        `json:"frequency_penalty,omitempty" yaml:"frequency_penalty,omitempty"`
	PresencePenalty   float64        `json:"presence_penalty,omitempty" yaml:"presence_penalty,omitempty"`
	JSON              bool           `json:"json,omitempty" yaml:"json,omitempty"`
	ResponseMIMEType  string         `json:"response_mime_type,omitempty" yaml:"response_mime_type,omitempty"`
	raw               map[string]any `json:"-" yaml:"-"`
}

// ProvidersConfig represents the configuration for all agent models for Ollama
type ProvidersConfig struct {
	Simple         *AgentConfig      `json:"simple,omitempty" yaml:"simple,omitempty"`
	SimpleJSON     *AgentConfig      `json:"simple_json,omitempty" yaml:"simple_json,omitempty"`
	Agent          *AgentConfig      `json:"agent,omitempty" yaml:"agent,omitempty"`
	Generator      *AgentConfig      `json:"generator,omitempty" yaml:"generator,omitempty"`
	Refiner        *AgentConfig      `json:"refiner,omitempty" yaml:"refiner,omitempty"`
	Adviser        *AgentConfig      `json:"adviser,omitempty" yaml:"adviser,omitempty"`
	Reflector      *AgentConfig      `json:"reflector,omitempty" yaml:"reflector,omitempty"`
	Searcher       *AgentConfig      `json:"searcher,omitempty" yaml:"searcher,omitempty"`
	Enricher       *AgentConfig      `json:"enricher,omitempty" yaml:"enricher,omitempty"`
	Coder          *AgentConfig      `json:"coder,omitempty" yaml:"coder,omitempty"`
	Installer      *AgentConfig      `json:"installer,omitempty" yaml:"installer,omitempty"`
	Pentester      *AgentConfig      `json:"pentester,omitempty" yaml:"pentester,omitempty"`
	defaultOptions []llms.CallOption `json:"-" yaml:"-"`
}

func LoadConfig(configPath string, defaultOptions []llms.CallOption) (*ProvidersConfig, error) {
	if configPath == "" {
		return nil, nil // No config path provided, return nil to use defaults
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ollama config file: %w", err)
	}

	var config ProvidersConfig
	ext := filepath.Ext(configPath)
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse ollama JSON config: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse ollama YAML config: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported ollama config file format: %s", ext)
	}

	config.defaultOptions = defaultOptions
	return &config, nil
}

func (ac *AgentConfig) UnmarshalJSON(data []byte) error {
	type embed AgentConfig
	var unmarshaler embed
	if err := json.Unmarshal(data, &unmarshaler); err != nil {
		return err
	}
	*ac = AgentConfig(unmarshaler)

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	ac.raw = raw
	return nil
}

func (ac *AgentConfig) UnmarshalYAML(value *yaml.Node) error {
	type embed AgentConfig
	var unmarshaler embed
	if err := value.Decode(&unmarshaler); err != nil {
		return err
	}
	*ac = AgentConfig(unmarshaler)

	var raw map[string]any
	if err := value.Decode(&raw); err != nil {
		return err
	}
	ac.raw = raw
	return nil
}

func (ac *AgentConfig) BuildOptions() []llms.CallOption {
	if ac == nil || ac.raw == nil {
		return nil
	}

	var options []llms.CallOption

	if _, ok := ac.raw["model"]; ok && ac.Model != "" {
		options = append(options, llms.WithModel(ac.Model))
	}
	if _, ok := ac.raw["max_tokens"]; ok {
		options = append(options, llms.WithMaxTokens(ac.MaxTokens))
	}
	if _, ok := ac.raw["temperature"]; ok {
		options = append(options, llms.WithTemperature(ac.Temperature))
	}
	if _, ok := ac.raw["top_k"]; ok {
		options = append(options, llms.WithTopK(ac.TopK))
	}
	if _, ok := ac.raw["top_p"]; ok {
		options = append(options, llms.WithTopP(ac.TopP))
	}
	if _, ok := ac.raw["n"]; ok {
		options = append(options, llms.WithN(ac.N))
	}
	if _, ok := ac.raw["min_length"]; ok {
		options = append(options, llms.WithMinLength(ac.MinLength))
	}
	if _, ok := ac.raw["max_length"]; ok {
		options = append(options, llms.WithMaxLength(ac.MaxLength))
	}
	if _, ok := ac.raw["repetition_penalty"]; ok {
		options = append(options, llms.WithRepetitionPenalty(ac.RepetitionPenalty))
	}
	if _, ok := ac.raw["frequency_penalty"]; ok {
		options = append(options, llms.WithFrequencyPenalty(ac.FrequencyPenalty))
	}
	if _, ok := ac.raw["presence_penalty"]; ok {
		options = append(options, llms.WithPresencePenalty(ac.PresencePenalty))
	}
	if _, ok := ac.raw["json"]; ok {
		options = append(options, llms.WithJSONMode())
	}
	if _, ok := ac.raw["response_mime_type"]; ok && ac.ResponseMIMEType != "" {
		options = append(options, llms.WithResponseMIMEType(ac.ResponseMIMEType))
	}

	return options
}

func (pc *ProvidersConfig) GetOptionsForType(optType provider.ProviderOptionsType) []llms.CallOption {
	if pc == nil {
		return pc.defaultOptions // Return default if no specific config
	}

	var agentConfig *AgentConfig
	switch optType {
	case provider.OptionsTypeSimple:
		agentConfig = pc.Simple
	case provider.OptionsTypeSimpleJSON:
		agentConfig = pc.SimpleJSON
	case provider.OptionsTypeAgent:
		agentConfig = pc.Agent
	case provider.OptionsTypeGenerator:
		agentConfig = pc.Generator
	case provider.OptionsTypeRefiner:
		agentConfig = pc.Refiner
	case provider.OptionsTypeAdviser:
		agentConfig = pc.Adviser
	case provider.OptionsTypeReflector:
		agentConfig = pc.Reflector
	case provider.OptionsTypeSearcher:
		agentConfig = pc.Searcher
	case provider.OptionsTypeEnricher:
		agentConfig = pc.Enricher
	case provider.OptionsTypeCoder:
		agentConfig = pc.Coder
	case provider.OptionsTypeInstaller:
		agentConfig = pc.Installer
	case provider.OptionsTypePentester:
		agentConfig = pc.Pentester
	default:
		return pc.defaultOptions // Unknown type, return default
	}

	if agentConfig != nil {
		if options := agentConfig.BuildOptions(); options != nil {
			return options
		}
	}

	return pc.defaultOptions // No specific config for type, return default
}

type ollamaProvider struct {
	llm     *openai.LLM
	model   string // Default model for the provider
	options map[provider.ProviderOptionsType][]llms.CallOption
}

func New(cfg *config.Config) (provider.Provider, error) {
	if cfg.OllamaServerURL == "" {
		// Ollama is not configured, so we don't initialize it.
		// This allows other providers to be used if Ollama is not set up.
		return nil, nil
	}

	httpClient := http.DefaultClient
	if cfg.ProxyURL != "" {
		proxyURL, err := url.Parse(cfg.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
		}
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
	}

	// For Ollama, the API key is not strictly required but the library might expect it.
	// We use a placeholder if no specific key is common for local Ollama setups.
	// The model is specified in the config file, not globally for the client.
	client, err := openai.New(
		openai.WithToken("ollama"), // Placeholder token
		openai.WithBaseURL(cfg.OllamaServerURL),
		openai.WithHTTPClient(httpClient),
		openai.WithModel(""), // Model will be set per-call based on config
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama client: %w", err)
	}

	// Default call options if no config file is provided or a type is missing
	defaultCallOptions := []llms.CallOption{
		llms.WithTemperature(0.7),
		llms.WithTopP(1.0),
		llms.WithN(1),
		llms.WithMaxTokens(4000), // A reasonable default
	}

	var providerOpts map[provider.ProviderOptionsType][]llms.CallOption
	if cfg.OllamaConfigPath != "" {
		loadedConfig, err := LoadConfig(cfg.OllamaConfigPath, defaultCallOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to load ollama provider config: %w", err)
		}
		if loadedConfig != nil {
			providerOpts = make(map[provider.ProviderOptionsType][]llms.CallOption)
			for _, optType := range provider.AllProviderOptionsTypes {
				if opts := loadedConfig.GetOptionsForType(optType); opts != nil {
					providerOpts[optType] = opts
				} else {
					providerOpts[optType] = defaultCallOptions // Fallback to defaults if type not in config
				}
			}
		}
	}

	// If no config file was loaded, or it was empty, populate with defaults for all types
	if providerOpts == nil {
		providerOpts = make(map[provider.ProviderOptionsType][]llms.CallOption)
		for _, optType := range provider.AllProviderOptionsTypes {
			providerOpts[optType] = defaultCallOptions
		}
	}

	// Determine a default model for the provider. This could be the 'Simple' model or any other.
	// This default model is used if a specific ProviderOptionsType doesn't specify one.
	defaultModel := "default" // Placeholder, actual model comes from config
	if simpleOpts, ok := providerOpts[provider.OptionsTypeSimple]; ok {
		tempOpts := llms.CallOptions{}
		for _, o := range simpleOpts {
			o(&tempOpts)
		}
		if tempOpts.Model != "" {
			defaultModel = tempOpts.Model
		}
	}


	return &ollamaProvider{
		llm:     client,
		model:   defaultModel,
		options: providerOpts,
	}, nil
}

func (p *ollamaProvider) Type() provider.ProviderType {
	return provider.ProviderOllama // Assuming ProviderOllama will be added to provider.ProviderType
}

func (p *ollamaProvider) Model(opt provider.ProviderOptionsType) string {
	options, ok := p.options[opt]
	if !ok || options == nil {
		// Fallback to the globally set default model for the provider if specific options are missing
		return p.model
	}

	callOpts := llms.CallOptions{Model: p.model} // Start with the provider's default model
	for _, option := range options {
		option(&callOpts)
	}

	return callOpts.Model
}

func (p *ollamaProvider) Call(
	ctx context.Context,
	opt provider.ProviderOptionsType,
	prompt string,
) (string, error) {
	options, ok := p.options[opt]
	if !ok || options == nil {
		return "", fmt.Errorf("%w: %s", provider.ErrInvalidProviderOptionsType, opt)
	}
	return provider.WrapGenerateFromSinglePrompt(ctx, p, opt, p.llm, prompt, options...)
}

func (p *ollamaProvider) CallEx(
	ctx context.Context,
	opt provider.ProviderOptionsType,
	chain []llms.MessageContent,
) (*llms.ContentResponse, error) {
	options, ok := p.options[opt]
	if !ok || options == nil {
		return nil, fmt.Errorf("%w: %s", provider.ErrInvalidProviderOptionsType, opt)
	}
	return provider.WrapGenerateContent(ctx, p, opt, p.llm.GenerateContent, chain, options...)
}

func (p *ollamaProvider) CallWithTools(
	ctx context.Context,
	opt provider.ProviderOptionsType,
	chain []llms.MessageContent,
	tools []llms.Tool,
) (*llms.ContentResponse, error) {
	options, ok := p.options[opt]
	if !ok || options == nil {
		return nil, fmt.Errorf("%w: %s", provider.ErrInvalidProviderOptionsType, opt)
	}

	finalOptions := make([]llms.CallOption, len(options))
	copy(finalOptions, options)
	finalOptions = append(finalOptions, llms.WithTools(tools))

	return provider.WrapGenerateContent(ctx, p, opt, p.llm.GenerateContent, chain, finalOptions...)
}

func (p *ollamaProvider) GetUsage(info map[string]any) (int64, int64) {
	var inputTokens, outputTokens int64
	if value, ok := info["PromptTokens"]; ok {
		if vInt, isInt := value.(int); isInt {
			inputTokens = int64(vInt)
		} else if vFloat, isFloat := value.(float64); isFloat { // Ollama might return float
			inputTokens = int64(vFloat)
		}
	}

	if value, ok := info["CompletionTokens"]; ok {
		if vInt, isInt := value.(int); isInt {
			outputTokens = int64(vInt)
		} else if vFloat, isFloat := value.(float64); isFloat {
			outputTokens = int64(vFloat)
		}
	}
	return inputTokens, outputTokens
}

// Ensure ollamaProvider implements the provider.Provider interface
var _ provider.Provider = &ollamaProvider{}
