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

	prompt := fmt.Sprintf(`你是一位资深的舆情监测分析师，具有丰富的网络舆情研判和危机应对经验。请对以下从短视频平台提取的内容（包含封面文字和音频转录）进行专业的舆情监测分析。

**待分析内容：**
%s

---

请基于舆情监测的专业框架进行深度分析，并严格按照JSON格式返回结果（不要包含markdown代码块标记或其他任何文字）：

{
  "sentiment_score": 0.75,
  "sentiment_label": "positive",
  "key_topics": ["话题1", "话题2", "话题3"],
  "risk_level": "low",
  "detailed_analysis": "完整的舆情分析报告内容...",
  "recommendations": ["策略1", "策略2", "策略3", "策略4"]
}

**字段详细说明：**

1. **sentiment_score** (舆情指数)
   - 数值范围：0.0-1.0的浮点数（保留两位小数）
   - 评分标准：
     * 0.0-0.2：强负面（严重批评、恶意攻击、负面情绪极强）
     * 0.2-0.4：负面（质疑、不满、批评为主）
     * 0.4-0.6：中性（客观陈述、情感中立、观点平衡）
     * 0.6-0.8：正面（支持、认可、积极态度）
     * 0.8-1.0：强正面（高度赞扬、热烈支持、正面情绪极强）
   - 综合考量：内容立场、情绪强度、用词倾向、价值导向

2. **sentiment_label** (舆情态度标签)
   - 可选值：
     * "positive"（正面）：score ≥ 0.6，内容持支持、赞扬、认可态度
     * "neutral"（中性）：0.4 ≤ score < 0.6，内容客观中立，无明显倾向
     * "negative"（负面）：score < 0.4，内容持批评、质疑、反对态度
   - 标签应与舆情指数保持一致

3. **key_topics** (核心舆情话题)
   - 提取3-6个关键主题标签
   - 涵盖维度：
     * 核心人物/机构名称
     * 事件/事项名称
     * 行业/领域分类
     * 争议焦点/关注点
     * 情感关键词
   - 示例：["某企业", "产品质量", "消费者维权", "品牌信誉", "食品安全"]

4. **risk_level** (舆情风险等级)
   - 风险分级：
     * "high"（高风险）：
       - 涉及重大公共安全、法律违规、道德失范
       - 公众人物负面事件、企业重大丑闻
       - 社会敏感话题、群体性争议
       - 可能引发舆论风暴、大规模传播、媒体跟进
       - 需要立即响应和危机公关
     
     * "medium"（中风险）：
       - 存在争议性观点、部分负面声音
       - 话题有一定传播潜力但影响可控
       - 需要持续关注和适度回应
       - 可能演变为高风险需提前预警
     
     * "low"（低风险）：
       - 正面内容或中性信息
       - 无明显争议点和负面影响
       - 常规监测即可，无需特殊应对

5. **detailed_analysis** (舆情分析报告)
   - 字数要求：500-800字
   - 报告结构（必须全部包含）：
   
     **【内容概述】**（100-150字）
     - 准确概括视频的核心内容和主要信息
     - 明确指出内容类型（资讯、评论、揭露、娱乐等）
     - 识别内容创作者立场和表达意图
     
     **【舆情态度分析】**（150-200字）
     - 详细解释舆情指数的评判依据
     - 分析内容的情感倾向和价值导向
     - 识别关键词、语气、论调的情感色彩
     - 评估内容对目标对象的态度（支持/中立/反对）
     
     **【传播趋势研判】**（100-150字）
     - 预测内容的传播潜力和扩散范围
     - 分析目标受众群体及其可能的反应
     - 评估话题在社交媒体的发酵可能性
     - 判断是否可能引发二次传播或媒体关注
     
     **【潜在影响评估】**（100-150字）
     - 分析对相关方（个人/企业/机构）的影响
     - 评估对品牌形象、公众信任的影响程度
     - 识别可能的连锁反应和衍生风险
     - 提出需要重点关注的风险点
     
     **【舆论引导建议】**（50-100字）
     - 建议舆情应对的基本方向
     - 提示关键的沟通策略和话语权把控
   
   - 语言要求：
     * 专业、严谨、客观
     * 避免主观臆断，基于内容事实
     * 使用舆情监测行业术语
     * 逻辑清晰，结构分明

6. **recommendations** (应对策略建议)
   - 提供3-6条分层次、可执行的应对策略
   - 策略分类：
   
     **即时应对**（针对高/中风险）：
     - 舆情监测加强措施
     - 紧急响应预案
     - 信息发布建议
     - 舆论引导方向
     
     **中期管理**：
     - 持续跟踪要点
     - 互动沟通策略
     - 内容优化方向
     - 风险预警机制
     
     **长期策略**：
     - 品牌形象建设
     - 舆情防御体系
     - 公众关系维护
     - 信任重建路径
   
   - 每条建议应：
     * 具体明确，可操作性强
     * 针对性强，符合实际情况
     * 提供清晰的执行方向
     * 考虑资源和可行性

**分析要求：**
- 必须综合分析封面文字和音频内容，确保信息完整性
- 高度关注敏感词汇、争议观点、价值导向
- 深入挖掘潜在的舆情风险和传播隐患
- 保持专业的第三方中立立场
- 如果封面与音频信息存在差异，以整体内容为准进行综合判断
- 特别注意识别可能的虚假信息、误导性内容、情绪煽动
- 报告必须达到500-800字，确保分析的深度和全面性`, coverText)

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
