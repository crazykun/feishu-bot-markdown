package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

/**
 * @Description: 生成消息卡片（飞书卡片2.0）
 * 飞书消息卡片结构文档 https://open.feishu.cn/document/uAjLw4CM/ukzMukzMukzM/feishu-cards/card-json-structure
 * 卡片搭建工具 https://open.feishu.cn/cardkit
 * 卡片的正文内容，支持配置多语言。卡片的正文内容支持多种模块，包括多列布局、内容模块、分割线、图片、备注、交互模块等。
 * 在卡片的正文内容中，支持添加以下属性：
 * column_set：多列布局，可以横向排布多个列容器，在列内纵向自由组合图文内容，解决多列内容对齐问题，并实现了灵活的图文混排。
 * div：内容模块，以格式化的文本为主体，支持混合图片、交互组件的富文本内容。
 * markdown：使用 Markdown 标签构造富文本内容。
 * hr：模块之间的分割线。
 * img：用于展示图片的组件。
 * note：备注组件，用于展示卡片内的次要信息。
 * actions：交互模块, 可以添加按钮。使用交互组件可以实现消息卡片与用户之间的信息交互。
 */

type any = interface{}

// Msg 飞书消息结构
type Msg struct {
	MsgType string `json:"msg_type"`
	Card    Card   `json:"card"`
}

// Text 文本对象
type Text struct {
	Content string `json:"content,omitempty"`
	Tag     string `json:"tag"`
}

// Header 卡片头部
type Header struct {
	Title    Text   `json:"title"`
	Template string `json:"template,omitempty"`
	UdIcon   *Icon  `json:"ud_icon,omitempty"` // 自定义图标（卡片2.0新增）
}

// Icon 图标对象
type Icon struct {
	Tag   string `json:"tag"`
	Token string `json:"token"`
	Color string `json:"color,omitempty"`
}

// Config 卡片全局配置（卡片2.0新增）
type Config struct {
	WideScreenMode bool `json:"wide_screen_mode,omitempty"` // 是否启用宽屏模式
	EnableForward  bool `json:"enable_forward,omitempty"`   // 是否允许转发
}

// CardLink 卡片链接
type CardLink struct {
	Url        string `json:"url"`
	AndroidUrl string `json:"android_url,omitempty"`
	IosUrl     string `json:"ios_url,omitempty"`
	PCUrl      string `json:"pc_url,omitempty"`
}

// Card 卡片主体
type Card struct {
	Config       *Config   `json:"config,omitempty"`         // 全局配置
	Header       Header    `json:"header"`
	Elements     []Element `json:"elements"`
	CardLink     *CardLink `json:"card_link,omitempty"`      // 卡片链接
	I18nElements *I18nElements `json:"i18n_elements,omitempty"` // 国际化元素
}

// I18nElements 国际化元素配置
type I18nElements struct {
	ZhCn *I18nElement `json:"zh_cn,omitempty"` // 中文
	EnUs *I18nElement `json:"en_us,omitempty"` // 英文
	JaJp *I18nElement `json:"ja_jp,omitempty"` // 日文
}

// I18nElement 单个语言的元素配置
type I18nElement struct {
	Elements []Element `json:"elements"`
	Header   *Header   `json:"header,omitempty"`
}

// Column 表示卡片中多列布局中的一列。可以包含多个元素，例如文本、图片等
type Column struct {
	Tag           string    `json:"tag"`
	Width         string    `json:"width,omitempty"`
	Weight        int       `json:"weight,omitempty"`
	Elements      []Element `json:"elements"`
	VerticalAlign string    `json:"vertical_align,omitempty"`
	Padding       *Padding  `json:"padding,omitempty"` // 内边距
}

// Padding 内边距配置
type Padding struct {
	Top    string `json:"top,omitempty"`
	Right  string `json:"right,omitempty"`
	Bottom string `json:"bottom,omitempty"`
	Left   string `json:"left,omitempty"`
}

// Action 表示卡片中的一个交互组件，例如按钮、选择器等
type Action struct {
	Tag     string      `json:"tag"`
	Text    *Text       `json:"text,omitempty"`
	Url     string      `json:"url,omitempty"`
	Type    string      `json:"type,omitempty"`
	Value   any         `json:"value,omitempty"`
	Confirm *Confirm    `json:"confirm,omitempty"` // 二次确认弹窗
	Options []*Option   `json:"options,omitempty"` // 下拉选项
}

// Confirm 二次确认弹窗配置
type Confirm struct {
	Title Text    `json:"title"`
	Text  Text    `json:"text"`
}

// Option 下拉选项
type Option struct {
	Text  Text   `json:"text"`
	Value string `json:"value"`
	Url   string `json:"url,omitempty"`
}

// Element 表示卡片中的一个元素，可以是多种类型，例如文本、图片、按钮等
type Element struct {
	Tag               string    `json:"tag"`
	TextAlign         string    `json:"text_align,omitempty"`
	Content           string    `json:"content,omitempty"`
	FlexMode          string    `json:"flex_mode,omitempty"`
	BackgroundStyle   string    `json:"background_style,omitempty"`
	HorizontalSpacing string    `json:"horizontal_spacing,omitempty"`
	Columns           []Column  `json:"columns,omitempty"`
	Actions           []Action  `json:"actions,omitempty"`
	Elements          []Element `json:"elements,omitempty"`
	
	// 图片相关字段
	ImgKey            string    `json:"img_key,omitempty"`
	Alt               *Text     `json:"alt,omitempty"`
	Title             *Text     `json:"title,omitempty"`
	CustomWidth       string    `json:"custom_width,omitempty"`
	CompactWidth      bool      `json:"compact_width,omitempty"`
	Mode              string    `json:"mode,omitempty"`
	
	// 文本样式
	Text              *Text     `json:"text,omitempty"`
	Extra             *Element  `json:"extra,omitempty"`
	Field             *Field    `json:"field,omitempty"`
	IsShort           bool      `json:"is_short,omitempty"`
}

// Field 字段对象（用于div模块）
type Field struct {
	IsShort bool  `json:"is_short"`
	Text    *Text `json:"text"`
}

// CreateMarkdownElement 构建一个 Markdown 元素，用于显示富文本内容
func CreateMarkdownElement(content string) Element {
	return Element{
		Tag:       "markdown",
		TextAlign: "left",
		Content:   content,
	}
}

// CreateMarkdownCenterElement 构建一个 Markdown 元素, 并居中，用于显示富文本内容
func CreateMarkdownCenterElement(content string) Element {
	return Element{
		Tag:       "markdown",
		TextAlign: "center",
		Content:   content,
	}
}

// CreateTextElement 构建一个纯文本元素
func CreateTextElement(content string) Element {
	return Element{
		Tag:     "plain_text",
		Content: content,
	}
}

// CreateNoteElement 构建一个备注元素，用于显示卡片内的次要信息
func CreateNoteElement(content string) Element {
	return Element{
		Tag: "note",
		Elements: []Element{
			CreateTextElement(content),
		},
	}
}

// CreateColumn 构建一个多列布局中的一列
func CreateColumn(align, content string) Column {
	return Column{
		Tag:           "column",
		Width:         "weighted",
		Weight:        1,
		VerticalAlign: align,
		Elements: []Element{
			CreateMarkdownElement(content),
		},
	}
}

// CreateCenterColumn 构建一个多列布局中的一列并居中
func CreateCenterColumn(align, content string) Column {
	return Column{
		Tag:           "column",
		Width:         "weighted",
		Weight:        1,
		VerticalAlign: align,
		Elements: []Element{
			CreateMarkdownCenterElement(content),
		},
	}
}

// CreateImageElement 构建一个图片元素
func CreateImageElement(imgKey, alt string) Element {
	return Element{
		Tag:    "img",
		ImgKey: imgKey,
		Alt: &Text{
			Content: alt,
			Tag:     "plain_text",
		},
	}
}

// CreateButtonElement 构建一个按钮元素
func CreateButtonElement(text, url string) Action {
	return Action{
		Tag: "button",
		Text: &Text{
			Content: text,
			Tag:     "plain_text",
		},
		Url:  url,
		Type: "default",
	}
}

// CreatePrimaryButtonElement 构建一个主按钮元素
func CreatePrimaryButtonElement(text, url string) Action {
	return Action{
		Tag: "button",
		Text: &Text{
			Content: text,
			Tag:     "plain_text",
		},
		Url:  url,
		Type: "primary",
	}
}

// Hr 构建一个模块之间的分割线
func Hr() Element {
	return Element{
		Tag: "hr",
	}
}

// CreateColumnSetElement 构建一个多列布局元素
func CreateColumnSetElement(columns []Column, flexMode string) Element {
	return Element{
		Tag:               "column_set",
		FlexMode:          flexMode,
		HorizontalSpacing: "default",
		Columns:           columns,
	}
}

type FeishuColor string

const (
	ColorBlue      FeishuColor = "blue"
	ColorWathet    FeishuColor = "wathet"
	ColorTurquoise FeishuColor = "turquoise"
	ColorGreen     FeishuColor = "green"
	ColorYellow    FeishuColor = "yellow"
	ColorOrange    FeishuColor = "orange"
	ColorRed       FeishuColor = "red"
	ColorCarmine   FeishuColor = "carmine"
	ColorViolet    FeishuColor = "violet"
	ColorGrey      FeishuColor = "grey"
	ColorDefault   FeishuColor = "default"
)

type FeishuMsg struct {
	Title         string         `json:"title"`                    // 标题
	Markdown      map[string]any `json:"markdown,omitempty"`       // 内容 (map形式，可能无序)
	MarkdownItems []Text         `json:"markdown_items,omitempty"` // 内容 (切片形式，保持顺序)
	MarkdownArray [][2]string    `json:"markdown_array,omitempty"` // 内容 (键值对数组形式，最简洁)
	Note          string         `json:"note"`                     // 备注
	NoteEmoji     bool           `json:"note_emoji"`               // 是否备注附带随机emoji表情
	Link          string         `json:"link,omitempty"`           // 链接
	HeaderColor   FeishuColor    `json:"-"`                        // 标题颜色
	Response      any            `json:"response,omitempty"`       // 响应内容
	
	// 卡片2.0新增字段
	WideScreen    bool           `json:"-"`                        // 是否启用宽屏模式
	EnableForward bool           `json:"-"`                        // 是否允许转发
	CustomIcon    *Icon          `json:"-"`                        // 自定义图标
	Actions       []Action       `json:"-"`                        // 交互组件（按钮等）
	Images        []string       `json:"-"`                        // 图片列表（img_key）
}

// buildMarkdownContent 构建markdown内容字符串
// 支持同时使用多种格式，按顺序输出：Markdown -> MarkdownItems -> MarkdownArray
func (f *FeishuMsg) buildMarkdownContent() string {
	var md strings.Builder

	// 1. 处理 Markdown map（可能无序）
	if len(f.Markdown) > 0 {
		for k, v := range f.Markdown {
			md.WriteString(fmt.Sprintf("**%s**：%s\n", k, v))
		}
	}

	// 2. 处理 MarkdownItems（保持顺序，支持混合内容）
	if len(f.MarkdownItems) > 0 {
		for _, item := range f.MarkdownItems {
			if item.Tag != "" {
				// 如果有 Tag，则格式化为键值对形式
				md.WriteString(fmt.Sprintf("**%s**：%s\n", item.Tag, item.Content))
			} else {
				// 如果没有 Tag，直接使用 Content
				md.WriteString(item.Content)
				md.WriteString("\n")
			}
		}
	}

	// 3. 处理 MarkdownArray（最简洁的键值对）
	if len(f.MarkdownArray) > 0 {
		for _, arr := range f.MarkdownArray {
			md.WriteString(fmt.Sprintf("**%s**：%s\n", arr[0], arr[1]))
		}
	}

	return md.String()
}

// buildNoteContent 构建备注内容
func (f *FeishuMsg) buildNoteContent() string {
	note := f.Note
	if note == "" {
		note = time.Now().Format("2006-01-02 15:04:05")
	}

	if f.NoteEmoji {
		// 随机生成一个emoji表情
		emoji := []string{"👍", "👏", "👌", "👊", "✌", "👋", "👆", "👇", "👈", "👉", "👎", "👓", "👔", "👕", "👖", "👗", "👘", "👙", "👚", "👛", "👜", "👝", "👞", "👟", "👠", "👡", "👢", "👣", "👤", "👥", "👦", "👧", "👨", "👩", "👪", "👫", "👬", "👭", "👮", "👯", "👰", "👱", "👲", "👳", "👴", "👵", "👶", "👷", "👸", "👹", "👺", "👻", "👼", "👽", "👾", "👿", "💀", "💁", "💂", "💃", "💄", "💅", "💆", "💇", "💈", "💉", "💊", "💋", "💌", "💍", "💎", "💏", "💐", "💑", "💒", "💓", "💔", "💕", "💖", "💗", "💘", "💙", "💚", "💛", "💜", "💝", "💞", "💟", "💠", "💡", "💢", "💣", "💤", "💥", "💦", "💧", "💨", "💩", "💪", "💫", "💬", "💭", "💮", "💯", "💰", "💱", "💲", "💳", "💴", "💵"}
		emojiIndex := rand.Intn(len(emoji))
		note = emoji[emojiIndex] + note + emoji[emojiIndex]
	}
	return note
}

// FormatMsg 构造一个统计消息卡片（飞书卡片2.0）
func FormatMsg(f *FeishuMsg) *Msg {
	elements := make([]Element, 0)

	// 添加markdown内容
	mdContent := f.buildMarkdownContent()
	if mdContent != "" {
		elements = append(elements, CreateMarkdownElement(mdContent))
	}

	// 添加图片（如果有）
	if len(f.Images) > 0 {
		for _, imgKey := range f.Images {
			elements = append(elements, CreateImageElement(imgKey, "图片"))
		}
	}

	// 添加交互组件（如果有）
	if len(f.Actions) > 0 {
		elements = append(elements, Element{
			Tag:     "action",
			Actions: f.Actions,
		})
	}

	// 添加备注
	noteContent := f.buildNoteContent()
	elements = append(elements, CreateNoteElement(noteContent))

	// 构建卡片链接
	var cardLink *CardLink
	if f.Link != "" {
		cardLink = &CardLink{
			Url: f.Link,
		}
	}

	// 构建全局配置
	var config *Config
	if f.WideScreen || f.EnableForward {
		config = &Config{
			WideScreenMode: f.WideScreen,
			EnableForward:  f.EnableForward,
		}
	}

	// 构建头部
	header := Header{
		Title: Text{
			Content: f.Title,
			Tag:     "plain_text",
		},
		Template: string(f.HeaderColor),
	}
	
	// 添加自定义图标（如果有）
	if f.CustomIcon != nil {
		header.UdIcon = f.CustomIcon
	}

	return &Msg{
		MsgType: "interactive",
		Card: Card{
			Config:   config,
			Header:   header,
			Elements: elements,
			CardLink: cardLink,
		},
	}
}

// SendFeishuMsg 发送消息到飞书
func SendFeishuMsg(hook string, f *FeishuMsg) error {
	if hook == "" {
		return fmt.Errorf("hook url is empty")
	}

	// 将消息内容转换为JSON格式
	msg := FormatMsg(f)
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// 创建HTTP POST请求
	req, err := http.NewRequest("POST", hook, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	f.Response = buf.String()
	return nil
}
