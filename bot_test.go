package bot

import (
	"strings"
	"testing"
)

// 飞书token信息
var Hook = "https://open.feishu.cn/open-apis/bot/v2/hook/xxxx-xxxx-xxxx-xxxx"

var TmplTestFeishu = &FeishuMsg{
	Title: "飞书markdown消息测试！",
	Markdown: map[string]any{
		"标题": "标题测试",
		"内容": "这是内容**粗体**, *斜体*, ~~删除线~~",
		"状态": "<font color='green'>成功</font> <font color='red'>失败</font> <font color='grey'>灰色</font>",
		"列表": "\n- 列表1\n- 列表2\n- 列表3",
		"链接": "[飞书机器人助手](https://www.feishu.cn/hc/zh-CN/articles/236028437163-%E6%9C%BA%E5%99%A8%E4%BA%BA%E6%B6%88%E6%81%AF%E5%86%85%E5%AE%B9%E6%94%AF%E6%8C%81%E7%9A%84%E6%96%87%E6%9C%AC%E6%A0%B7%E5%BC%8F)",
		"时间": `2021-08-12 12:00:00快速`,
	},
	Note:        "这是备注",
	NoteEmoji:   true,
	Link:        "http://www.baidu.com",
	HeaderColor: ColorWathet,
}

func TestSendTxt(t *testing.T) {
	// 如果Hook结尾是xxxx-xxxx-xxxx-xxxx，提示让输入正确的Hook
	if strings.HasSuffix(Hook, "xxxx-xxxx-xxxx-xxxx") {
		t.Fatal("请设置正确的Hook")
	}
	SendFeishuMsg(Hook, TmplTestFeishu)
	t.Log(TmplTestFeishu.Response)
}

func TestMarkdownItems(t *testing.T) {
	// 测试使用 MarkdownItems 的情况
	msg := &FeishuMsg{
		Title: "测试 MarkdownItems",
		MarkdownItems: []Text{
			{Tag: "第一项", Content: "这是第一个有序项目"},
			{Tag: "第二项", Content: "这是第二个有序项目"},
			{Tag: "第三项", Content: "这是第三个有序项目"},
		},
		Note:        "使用 MarkdownItems 保持顺序",
		HeaderColor: ColorGreen,
	}

	// 测试构建内容
	content := msg.buildMarkdownContent()
	expected := "**第一项**：这是第一个有序项目\n**第二项**：这是第二个有序项目\n**第三项**：这是第三个有序项目\n"
	if content != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, content)
	}

	t.Log("MarkdownItems 测试通过")
}

func TestMarkdownFallback(t *testing.T) {
	// 测试当 Markdown 为空时使用 MarkdownItems 的情况
	msg := &FeishuMsg{
		Title:    "测试回退机制",
		Markdown: map[string]any{}, // 空的 map
		MarkdownItems: []Text{
			{Tag: "回退内容", Content: "当 Markdown 为空时使用此内容"},
		},
		Note:        "测试回退机制",
		HeaderColor: ColorYellow,
	}

	content := msg.buildMarkdownContent()
	expected := "**回退内容**：当 Markdown 为空时使用此内容\n"
	if content != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, content)
	}

	t.Log("回退机制测试通过")
}

func TestMarkdownArray(t *testing.T) {
	// 测试使用 MarkdownArray 的情况
	msg := &FeishuMsg{
		Title: "测试 MarkdownArray",
		MarkdownArray: [][2]string{
			{"第一项", "这是第一个有序项目"},
			{"第二项", "这是第二个有序项目"},
			{"第三项", "这是第三个有序项目"},
		},
		Note:        "使用 MarkdownArray 最简洁的方式",
		HeaderColor: ColorBlue,
	}

	// 测试构建内容
	content := msg.buildMarkdownContent()
	expected := "**第一项**：这是第一个有序项目\n**第二项**：这是第二个有序项目\n**第三项**：这是第三个有序项目\n"
	if content != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, content)
	}

	t.Log("MarkdownArray 测试通过")
}

func TestMarkdownPriority(t *testing.T) {
	// 测试当前实现：同时支持多种格式
	msg := &FeishuMsg{
		Title: "测试多种格式同时使用",
		Markdown: map[string]any{
			"Map格式": "来自 Markdown map",
		},
		MarkdownItems: []Text{
			{Tag: "Items格式", Content: "来自 MarkdownItems"},
		},
		MarkdownArray: [][2]string{
			{"Array格式", "来自 MarkdownArray"},
		},
		Note:        "测试多种格式同时使用",
		HeaderColor: ColorRed,
	}

	content := msg.buildMarkdownContent()
	t.Logf("生成的内容:\n%s", content)

	// 现在应该包含所有三种格式的内容
	if !strings.Contains(content, "来自 Markdown map") {
		t.Errorf("应该包含 Markdown 的内容")
	}
	if !strings.Contains(content, "来自 MarkdownItems") {
		t.Errorf("应该包含 MarkdownItems 的内容")
	}
	if !strings.Contains(content, "来自 MarkdownArray") {
		t.Errorf("应该包含 MarkdownArray 的内容")
	}

	t.Log("多格式同时使用测试通过")
}

// 测试卡片2.0 - 宽屏模式
func TestCardV2WideScreen(t *testing.T) {
	msg := &FeishuMsg{
		Title:         "测试宽屏模式",
		MarkdownArray: [][2]string{{"测试", "宽屏模式"}},
		WideScreen:    true,
		EnableForward: true,
		HeaderColor:   ColorBlue,
	}

	card := FormatMsg(msg)

	if card.Card.Config == nil {
		t.Error("Config 不应该为 nil")
	}

	if !card.Card.Config.WideScreenMode {
		t.Error("WideScreenMode 应该为 true")
	}

	if !card.Card.Config.EnableForward {
		t.Error("EnableForward 应该为 true")
	}

	t.Log("卡片2.0宽屏模式测试通过")
}

// 测试卡片2.0 - 交互按钮
func TestCardV2Actions(t *testing.T) {
	msg := &FeishuMsg{
		Title:         "测试交互按钮",
		MarkdownArray: [][2]string{{"测试", "交互按钮"}},
		Actions: []Action{
			CreatePrimaryButtonElement("确认", "https://example.com"),
			CreateButtonElement("取消", ""),
		},
		HeaderColor: ColorGreen,
	}

	card := FormatMsg(msg)

	// 查找 action 元素
	hasAction := false
	for _, elem := range card.Card.Elements {
		if elem.Tag == "action" && len(elem.Actions) > 0 {
			hasAction = true
			if len(elem.Actions) != 2 {
				t.Errorf("应该有2个按钮，实际有 %d 个", len(elem.Actions))
			}

			// 检查第一个按钮类型
			if elem.Actions[0].Type != "primary" {
				t.Error("第一个按钮应该是 primary 类型")
			}

			// 检查第二个按钮类型
			if elem.Actions[1].Type != "default" {
				t.Error("第二个按钮应该是 default 类型")
			}
			break
		}
	}

	if !hasAction {
		t.Error("应该包含 action 元素")
	}

	t.Log("卡片2.0交互按钮测试通过")
}

// 测试卡片2.0 - 自定义图标
func TestCardV2CustomIcon(t *testing.T) {
	msg := &FeishuMsg{
		Title:         "测试自定义图标",
		MarkdownArray: [][2]string{{"测试", "自定义图标"}},
		CustomIcon: &Icon{
			Tag:   "standard_icon",
			Token: "bell_outlined",
		},
		HeaderColor: ColorRed,
	}

	card := FormatMsg(msg)

	if card.Card.Header.UdIcon == nil {
		t.Error("UdIcon 不应该为 nil")
	}

	if card.Card.Header.UdIcon.Token != "bell_outlined" {
		t.Errorf("图标Token应该是 bell_outlined，实际是 %s", card.Card.Header.UdIcon.Token)
	}

	t.Log("卡片2.0自定义图标测试通过")
}

// 测试卡片2.0 - 卡片链接
func TestCardV2CardLink(t *testing.T) {
	msg := &FeishuMsg{
		Title:         "测试卡片链接",
		MarkdownArray: [][2]string{{"测试", "卡片链接"}},
		Link:          "https://www.feishu.cn",
		HeaderColor:   ColorYellow,
	}

	card := FormatMsg(msg)

	if card.Card.CardLink == nil {
		t.Error("CardLink 不应该为 nil")
	}

	if card.Card.CardLink.Url != "https://www.feishu.cn" {
		t.Errorf("卡片链接URL不正确，实际是 %s", card.Card.CardLink.Url)
	}

	t.Log("卡片2.0卡片链接测试通过")
}

// 测试创建图片元素
func TestCreateImageElement(t *testing.T) {
	imgKey := "img_v3_025h_xxxx"
	elem := CreateImageElement(imgKey, "测试图片")

	if elem.Tag != "img" {
		t.Errorf("元素标签应该是 img，实际是 %s", elem.Tag)
	}

	if elem.ImgKey != imgKey {
		t.Errorf("图片key不正确，实际是 %s", elem.ImgKey)
	}

	if elem.Alt == nil || elem.Alt.Content != "测试图片" {
		t.Error("图片alt文本不正确")
	}

	t.Log("创建图片元素测试通过")
}

// 测试创建多列布局
func TestCreateColumnSetElement(t *testing.T) {
	columns := []Column{
		CreateColumn("top", "第一列内容"),
		CreateColumn("top", "第二列内容"),
	}

	elem := CreateColumnSetElement(columns, "bisect")

	if elem.Tag != "column_set" {
		t.Errorf("元素标签应该是 column_set，实际是 %s", elem.Tag)
	}

	if len(elem.Columns) != 2 {
		t.Errorf("应该有2列，实际有 %d 列", len(elem.Columns))
	}

	if elem.FlexMode != "bisect" {
		t.Errorf("FlexMode应该是 bisect，实际是 %s", elem.FlexMode)
	}

	t.Log("创建多列布局测试通过")
}

// 测试带确认弹窗的按钮
func TestButtonWithConfirm(t *testing.T) {
	action := Action{
		Tag: "button",
		Text: &Text{
			Content: "删除",
			Tag:     "plain_text",
		},
		Type: "danger",
		Confirm: &Confirm{
			Title: Text{
				Content: "确认删除？",
				Tag:     "plain_text",
			},
			Text: Text{
				Content: "删除后不可恢复",
				Tag:     "plain_text",
			},
		},
	}

	if action.Confirm == nil {
		t.Error("Confirm 不应该为 nil")
	}

	if action.Confirm.Title.Content != "确认删除？" {
		t.Error("确认标题不正确")
	}

	t.Log("带确认弹窗的按钮测试通过")
}
