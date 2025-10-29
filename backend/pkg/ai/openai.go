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

	prompt := fmt.Sprintf(`你是一位专业的舆情分析专家。请分析以下从短视频中提取的文本内容（包含封面文字和音频转录），进行全面的舆情监测分析。

**分析内容：**
%s

---

请从以下维度进行深度分析，并严格按照JSON格式返回结果（不要包含任何其他文字）：

{
  "sentiment_score": 0.75,
  "sentiment_label": "positive",
  "key_topics": ["话题1", "话题2", "话题3"],
  "risk_level": "low",
  "detailed_analysis": "详细的情感分析内容",
  "recommendations": ["建议1", "建议2", "建议3"]
}

**字段说明：**

1. **sentiment_score** (情感得分)
   - 类型：0.0-1.0之间的浮点数
   - 标准：0.0=极端负面，0.3=负面，0.5=中性，0.7=正面，1.0=极端正面
   - 考虑因素：用词情绪、论调态度、话题敏感度

2. **sentiment_label** (情感标签)
   - 可选值："positive"(正面)、"neutral"(中性)、"negative"(负面)
   - 判断依据：
     * positive: score > 0.6，表达积极、赞扬、支持的态度
     * neutral: 0.4 ≤ score ≤ 0.6，客观陈述或情感模糊
     * negative: score < 0.4，表达批评、质疑、不满的态度

3. **key_topics** (关键主题)
   - 提取3-5个核心话题关键词
   - 应包括：人物/事件名称、行业领域、争议焦点
   - 示例：["于书欣", "职场霸凌", "综艺黑幕", "公众人物责任"]

4. **risk_level** (舆情风险等级)
   - 可选值："low"(低风险)、"medium"(中等风险)、"high"(高风险)
   - 评估标准：
     * high: 涉及重大负面事件、公众人物丑闻、社会敏感话题，可能引发广泛关注和负面影响
     * medium: 存在争议性内容、部分负面声音，但影响范围有限
     * low: 正面或中性内容，无明显风险

5. **detailed_analysis** (详细分析)
   - 字数：150-300字
   - 必须包含：
     * 内容概述（简述视频主要信息）
     * 情感倾向分析（解释情感得分的依据）
     * 舆论走向判断（可能的公众反应）
     * 潜在影响评估（对相关方的影响）
   - 语言风格：专业、客观、有深度

6. **recommendations** (应对建议)
   - 提供2-4条具体可行的建议
   - 针对对象：内容创作者、相关当事人、舆情监测团队
   - 建议类型：
     * 内容优化建议
     * 风险应对措施
     * 公关沟通策略
     * 持续监测要点

**分析要求：**
- 综合考虑封面文字和音频内容，不要遗漏任何重要信息
- 关注事实陈述、情感倾向、价值观导向
- 识别潜在的敏感词、争议点、舆论风险
- 保持客观中立，基于内容本身进行分析
- 如果封面和音频传递的情感不一致，以整体综合判断为准`, coverText)

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
