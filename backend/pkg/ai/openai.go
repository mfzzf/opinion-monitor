package ai

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenAIClient struct {
	client      openai.Client
	ModelVision string
	ModelChat   string
}

type SentimentReport struct {
	SentimentScore   float64  `json:"sentiment_score"`
	SentimentLabel   string   `json:"sentiment_label"`
	KeyTopics        []string `json:"key_topics"`
	RiskLevel        string   `json:"risk_level"`
	DetailedAnalysis string   `json:"detailed_analysis"`
	Recommendations  []string `json:"recommendations"`
}

func NewOpenAIClient(apiBase, apiKey, modelVision, modelChat string) *OpenAIClient {
	// 创建客户端选项
	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}

	// 如果提供了自定义 API Base URL，则使用它
	fmt.Println("apiBase", apiBase)
	if apiBase != "" {
		// 移除尾部斜杠（如果有）
		apiBase = strings.TrimSuffix(apiBase, "/")
		opts = append(opts, option.WithBaseURL(apiBase))
	}

	client := openai.NewClient(opts...)

	return &OpenAIClient{
		client:      client,
		ModelVision: modelVision,
		ModelChat:   modelChat,
	}
}

func (c *OpenAIClient) ExtractTextFromImage(imagePath string) (string, error) {
	ctx := context.Background()

	// 读取图片文件
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	// 编码为 base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)
	imageURL := fmt.Sprintf("data:image/jpeg;base64,%s", base64Image)

	// 创建聊天完成请求
	chatCompletion, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: c.ModelVision,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
				openai.TextContentPart("convert to markdown"),
				openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
					URL: imageURL,
				}),
			}),
		},
		MaxTokens: openai.Int(1000),
	})

	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) AnalyzeSentiment(coverText string) (*SentimentReport, error) {
	ctx := context.Background()

	prompt := fmt.Sprintf(`分析以下从短视频封面中提取的文字内容，进行舆情监测分析。

文字内容：%s

请提供详细的舆情分析报告，严格按照以下JSON格式返回（不要包含任何其他文字）：
{
  "sentiment_score": 0.75,
  "sentiment_label": "positive",
  "key_topics": ["话题1", "话题2", "话题3"],
  "risk_level": "low",
  "detailed_analysis": "详细的情感分析，包括情感倾向、主要观点、可能的影响等",
  "recommendations": ["建议1", "建议2", "建议3"]
}

注意：
- sentiment_score: 0-1之间的浮点数，表示情感倾向（0=极端负面，0.5=中性，1=极端正面）
- sentiment_label: "positive"(正面)、"neutral"(中性)或"negative"(负面)
- key_topics: 3-5个关键主题
- risk_level: "low"(低风险)、"medium"(中等风险)或"high"(高风险)
- detailed_analysis: 100-200字的详细分析
- recommendations: 2-3条可操作的建议`, coverText)

	// 创建聊天完成请求
	chatCompletion, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: c.ModelChat,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return nil, fmt.Errorf("no response from API")
	}

	// 解析 JSON 响应
	content := chatCompletion.Choices[0].Message.Content

	// 清理 markdown 代码块（如果存在）
	content = cleanJSONResponse(content)

	var report SentimentReport
	if err := json.Unmarshal([]byte(content), &report); err != nil {
		return nil, fmt.Errorf("failed to parse sentiment report: %w, content: %s", err, content)
	}

	return &report, nil
}

func cleanJSONResponse(content string) string {
	// 移除 markdown 代码块标记
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	return strings.TrimSpace(content)
}
