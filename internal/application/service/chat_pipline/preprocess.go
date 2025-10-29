package chatpipline

import (
	"context"
	"regexp"
	"strings"
	"unicode"

	"github.com/Tencent/WeKnora/internal/config"
	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
	"github.com/yanyiwu/gojieba"
)

// PluginPreprocess Query preprocessing plugin
type PluginPreprocess struct {
	config    *config.Config
	jieba     *gojieba.Jieba
	stopwords map[string]struct{}
}

// Regular expressions for text cleaning
var (
	multiSpaceRegex = regexp.MustCompile(`\s+`)                                 // Multiple spaces
	urlRegex        = regexp.MustCompile(`https?://\S+`)                        // URLs
	emailRegex      = regexp.MustCompile(`\b[\w.%+-]+@[\w.-]+\.[a-zA-Z]{2,}\b`) // Email addresses
	punctRegex      = regexp.MustCompile(`[^\p{L}\p{N}\s]`)                     // Punctuation marks
)

// NewPluginPreprocess Creates a new query preprocessing plugin
func NewPluginPreprocess(
	eventManager *EventManager,
	config *config.Config,
	cleaner interfaces.ResourceCleaner,
) *PluginPreprocess {
	// Use default dictionary for Jieba tokenizer
	jieba := gojieba.NewJieba()

	// Load stopwords from built-in stopword library
	stopwords := loadStopwords()

	res := &PluginPreprocess{
		config:    config,
		jieba:     jieba,
		stopwords: stopwords,
	}

	// Register resource cleanup function
	if cleaner != nil {
		cleaner.RegisterWithName("JiebaPreprocessor", func() error {
			res.Close()
			return nil
		})
	}

	eventManager.Register(res)
	return res
}

// Load stopwords
func loadStopwords() map[string]struct{} {
	// Directly use some common stopwords built into Jieba
	commonStopwords := []string{
		"的", "了", "和", "是", "在", "我", "你", "他", "她", "它",
		"这", "那", "什么", "怎么", "如何", "为什么", "哪里", "什么时候",
		"the", "is", "are", "am", "I", "you", "he", "she", "it", "this",
		"that", "what", "how", "a", "an", "and", "or", "but", "if", "of",
		"to", "in", "on", "at", "by", "for", "with", "about", "from",
		"有", "无", "好", "来", "去", "说", "看", "想", "会", "可以",
		"吗", "呢", "啊", "吧", "的话", "就是", "只是", "因为", "所以",
	}

	result := make(map[string]struct{}, len(commonStopwords))
	for _, word := range commonStopwords {
		result[word] = struct{}{}
	}
	return result
}

// ActivationEvents Register activation events
func (p *PluginPreprocess) ActivationEvents() []types.EventType {
	return []types.EventType{types.PREPROCESS_QUERY}
}

// OnEvent Process events
func (p *PluginPreprocess) OnEvent(ctx context.Context, eventType types.EventType, chatManage *types.ChatManage, next func() *PluginError) *PluginError {
	if chatManage.RewriteQuery == "" {
		return next()
	}

	logger.GetLogger(ctx).Infof("Starting query preprocessing, original query: %s", chatManage.RewriteQuery)

	// 1. Basic text cleaning
	processed := p.cleanText(chatManage.RewriteQuery)

	// 2. Tokenization
	segments := p.segmentText(processed)

	// 3. Stopword filtering and reconstruction
	filteredSegments := p.filterStopwords(segments)

	// Update preprocessed query
	chatManage.ProcessedQuery = strings.Join(filteredSegments, " ")

	logger.GetLogger(ctx).Infof("Query preprocessing complete, processed query: %s", chatManage.ProcessedQuery)

	return next()
}

// cleanText Basic text cleaning
func (p *PluginPreprocess) cleanText(text string) string {
	// Remove URLs
	text = urlRegex.ReplaceAllString(text, " ")

	// Remove email addresses
	text = emailRegex.ReplaceAllString(text, " ")

	// Remove excessive spaces
	text = multiSpaceRegex.ReplaceAllString(text, " ")

	// Remove punctuation marks
	text = punctRegex.ReplaceAllString(text, " ")

	// Trim leading and trailing spaces
	text = strings.TrimSpace(text)

	return text
}

// segmentText Text tokenization
func (p *PluginPreprocess) segmentText(text string) []string {
	// Use Jieba tokenizer for tokenization, using search engine mode
	segments := p.jieba.CutForSearch(text, true)
	return segments
}

// filterStopwords Filter stopwords
func (p *PluginPreprocess) filterStopwords(segments []string) []string {
	var filtered []string

	for _, word := range segments {
		// If not a stopword and not blank, keep it
		if _, isStopword := p.stopwords[word]; !isStopword && !isBlank(word) {
			filtered = append(filtered, word)
		}
	}

	// If filtering results in empty list, return original tokenization results
	if len(filtered) == 0 {
		return segments
	}

	return filtered
}

// isBlank Check if a string is blank
func isBlank(str string) bool {
	for _, r := range str {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// Ensure resources are properly released
func (p *PluginPreprocess) Close() {
	if p.jieba != nil {
		p.jieba.Free()
		p.jieba = nil
	}
}

// ShutdownHandler Returns shutdown function
func (p *PluginPreprocess) ShutdownHandler() func() {
	return func() {
		p.Close()
	}
}
