package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

/**
 * @Description: 生成消息卡片
 * 飞书消息卡片结构文档 https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/card-structure/card-content
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
	Elemens           []Element `json:"elements,omitempty"`
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
		Elemens: []Element{
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
	Title       string         `json:"title"`          // 标题
	Markdown    map[string]any `json:"contents"`       // 内容
	Note        string         `json:"note"`           // 备注
	Link        string         `json:"link,omitempty"` // 链接
	HeaderColor FeishuColor    `json:"-"`              // 标题颜色
}

// NewStatMsg 构造一个统计消息卡片
func FormatMsg(f *FeishuMsg) *Msg {
	var elements []Element
	md := ""
	for k, v := range f.Markdown {
		md += fmt.Sprintf("**%s**：%s\n", k, v)
	}
	elements = append(elements, CreateMarkdownElement(md))

	// 添加备注
	if f.Note != "" {
		// 随机生成一个emoji表情
		emoji := []string{"👍", "👏", "👌", "👊", "✌", "👋", "👆", "👇", "👈", "👉", "👎", "👓", "👔", "👕", "👖", "👗", "👘", "👙", "👚", "👛", "👜", "👝", "👞", "👟", "👠", "👡", "👢", "👣", "👤", "👥", "👦", "👧", "👨", "👩", "👪", "👫", "👬", "👭", "👮", "👯", "👰", "👱", "👲", "👳", "👴", "👵", "👶", "👷", "👸", "👹", "👺", "👻", "👼", "👽", "👾", "👿", "💀", "💁", "💂", "💃", "💄", "💅", "💆", "💇", "💈", "💉", "💊", "💋", "💌", "💍", "💎", "💏", "💐", "💑", "💒", "💓", "💔", "💕", "💖", "💗", "💘", "💙", "💚", "💛", "💜", "💝", "💞", "💟", "💠", "💡", "💢", "💣", "💤", "💥", "💦", "💧", "💨", "💩", "💪", "💫", "💬", "💭", "💮", "💯", "💰", "💱", "💲", "💳", "💴", "💵"}
		emojiIndex := rand.Intn(len(emoji))
		elements = append(elements, CreateNoteElement(emoji[emojiIndex]+f.Note+emoji[emojiIndex]))
	}

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

// 发送消息
func SendFeishuMsg(hook string, f *FeishuMsg) {
	if hook == "" {
		return
	}

	// 将消息内容转换为JSON格式
	data, _ := json.Marshal(FormatMsg(f))

	// 创建HTTP POST请求
	req, _ := http.NewRequest("POST", hook, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求并打印响应结果
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		defer resp.Body.Close()
	}
}
