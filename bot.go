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
 * @Description: 生成消息卡片
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

type Msg struct {
	MsgType string `json:"msg_type"`
	Card    Card   `json:"card"`
}

type Text struct {
	Content string `json:"content,omitempty"`
	Tag     string `json:"tag"`
}

type Header struct {
	Title    Text   `json:"title"`
	Template string `json:"template,omitempty"`
}

type Card struct {
	Elements []Element `json:"elements"`
	Header   Header    `json:"header"`
	CardLink CardLink  `json:"card_link,omitempty"`
}

type CardLink struct {
	Url        string `json:"url"`
	AndroidUrl string `json:"android_url,omitempty"`
	IosUrl     string `json:"ios_url,omitempty"`
	PCUrl      string `json:"pc_url,omitempty"`
}

// Column 表示卡片中多列布局中的一列。可以包含多个元素，例如文本、图片等
type Column struct {
	Tag           string    `json:"tag"`
	Width         string    `json:"width"`
	Weight        int       `json:"weight"`
	Elements      []Element `json:"elements"`
	VerticalAlign string    `json:"vertical_align,omitempty"`
}

// Action 表示卡片中的一个交互组件，例如按钮、选择器等
type Action struct {
	Tag   string `json:"tag"`
	Text  Text   `json:"text"`
	Url   string `json:"url"`
	Type  string `json:"type"`
	Value any    `json:"value"`
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

// Hr 构建一个模块之间的分割线
func Hr() Element {
	return Element{
		Tag: "hr",
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

// FormatMsg 构造一个统计消息卡片
func FormatMsg(f *FeishuMsg) *Msg {
	elements := make([]Element, 0)

	// 添加markdown内容
	mdContent := f.buildMarkdownContent()
	if mdContent != "" {
		elements = append(elements, CreateMarkdownElement(mdContent))
	}

	// 添加备注
	noteContent := f.buildNoteContent()
	elements = append(elements, CreateNoteElement(noteContent))

	return &Msg{
		MsgType: "interactive",
		Card: Card{
			Elements: elements,
			Header: Header{
				Title: Text{
					Content: f.Title,
					Tag:     "plain_text",
				},
				Template: string(f.HeaderColor),
			},
			CardLink: CardLink{
				Url: f.Link,
			},
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
