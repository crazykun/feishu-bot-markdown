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
