# Go飞书机器人Markdown消息工具

## 项目简介

一个基于Go语言封装的飞书群机器人消息发送工具，支持**飞书卡片2.0**格式，提供Markdown格式消息的快速构建与发送，简洁的API和丰富的消息格式支持。

---

## ✨ 新特性 - 飞书卡片2.0

本项目已全面升级至**飞书卡片2.0**，新增以下特性：

- 🎨 **宽屏模式**：支持更宽的展示区域，适合复杂内容
- 🔄 **转发控制**：可配置是否允许消息转发
- 🎯 **交互组件**：支持按钮、下拉选择等交互元素
- 🖼️ **自定义图标**：支持设置卡片头部自定义图标
- 🌍 **国际化支持**：支持多语言卡片内容
- 📐 **增强布局**：更灵活的多列布局和样式控制
- 🖼️ **图片支持**：直接嵌入图片到卡片中

---

## 安装

```bash
go get github.com/crazykun/feishu-bot-markdown
```

---

## 快速开始

### 1. 配置机器人Hook地址

```go
import (
	bot "github.com/crazykun/feishu-bot-markdown"
)

var Hook = "https://open.feishu.cn/open-apis/bot/v2/hook/xxxx-xxxx-xxxx-xxxx"
```

### 2. 构建消息模板

#### 方式一：使用 Map（传统方式）
```go
msg := &bot.FeishuMsg{
	Title: "任务完成通知",
	Markdown: map[string]any{
		"状态": "<font color='green'>成功</font>",
		"进度": "已完成 100%",
		"日志": "```日志内容```",
	},
	Note: "点击查看详细日志",
	NoteEmoji: true,
	Link: "http://example.com/logs",
}
```

#### 方式二：使用 MarkdownItems（灵活，保持顺序）
```go
msg := &bot.FeishuMsg{
	Title: "任务完成通知",
	MarkdownItems: []bot.Text{
		{Tag: "状态", Content: "<font color='green'>成功</font>"},
		{Tag: "进度", Content: "已完成 100%"},
		{Tag: "日志", Content: "```日志内容```"},
	},
	Note: "点击查看详细日志",
	NoteEmoji: true,
	Link: "http://example.com/logs",
}
```

#### 方式三：使用 MarkdownArray（推荐，最简洁）
```go
msg := &bot.FeishuMsg{
	Title: "任务完成通知",
	MarkdownArray: [][2]string{
		{"状态", "<font color='green'>成功</font>"},
		{"进度", "已完成 100%"},
		{"日志", "```日志内容```"},
	},
	Note: "点击查看详细日志",
	NoteEmoji: true,
	Link: "http://example.com/logs",
}
```

### 3. 发送消息

```go
bot.SendFeishuMsg(Hook, msg)
```

---

## 核心功能

- [x] 支持完整Markdown语法子集
- [x] 自定义消息卡片样式
- [x] 颜色配置（支持绿色/红色/灰色等）
- [x] 超链接与@用户功能
- [x] 代码块与列表渲染
- [x] 支持Note备注和随机Emoji
- [x] **新增**：MarkdownItems 切片支持，解决 map 转 JSON 无序问题
- [x] **新增**：MarkdownArray 键值对数组支持，提供最简洁的使用方式
- [x] **新增**：智能优先级机制，支持多种内容格式自动选择
- [x] **新增**：飞书卡片2.0完整支持


---

## 使用示例

![消息展示效果](https://raw.githubusercontent.com/crazykun/feishu-bot-markdown/main/src/screenshot.jpg)
> 图1：消息发送效果示例（包含标题、内容、状态提示和操作链接）

---

## Markdown语法支持

| 功能        | 语法示例                          | 展示效果                     | 说明                     |
|-------------|----------------------------------|----------------------------|-------------------------|
| 斜体        | *斜体文字*                       | *斜体文字*                  | 使用星号包裹            |
| 颜色        | <font color='green'>成功</font>  | <font color='green'>成功</font> | 支持green/red/grey      |
| 链接        | [飞书官网](https://feishu.cn)    | [飞书官网](https://feishu.cn) | 需要完整URL             |
| 代码块      | ```go\nfmt.Println("Hello")\n``` | ```go\nfmt.Println("Hello")\n``` | 支持语言高亮            |

---

## 高级用法

### 自定义卡片样式

```go
// 设置卡片头部颜色
msg.HeaderColor = bot.ColorWathet // 蓝色主题
msg.HeaderColor = bot.ColorGreen  // 绿色主题
msg.HeaderColor = bot.ColorRed    // 红色主题
msg.HeaderColor = bot.ColorGrey    // 灰色主题
msg.HeaderColor = bot.ColorDefault // 默认主题
```

### @指定用户

```go
msg.Markdown["负责人"] = `<at id=user_123>张三</at>`
```

### 解决 JSON 序列化顺序问题

提供三种解决方案来保证内容顺序：

```go
// 问题：map 在 JSON 序列化时顺序不固定
msg := &bot.FeishuMsg{
	Markdown: map[string]any{
		"第三步": "完成",
		"第一步": "开始", 
		"第二步": "进行中",
	},
}

// 解决方案1：使用 MarkdownItems（灵活）
msg := &bot.FeishuMsg{
	MarkdownItems: []bot.Text{
		{Tag: "第一步", Content: "开始"},
		{Tag: "第二步", Content: "进行中"},
		{Tag: "第三步", Content: "完成"},
	},
}

// 解决方案2：使用 MarkdownArray（最简洁，推荐）
msg := &bot.FeishuMsg{
	MarkdownArray: [][2]string{
		{"第一步", "开始"},
		{"第二步", "进行中"},
		{"第三步", "完成"},
	},
}
```

### 智能优先级机制

系统按以下优先级自动选择内容格式：**Markdown** > **MarkdownItems** > **MarkdownArray**

```go
msg := &bot.FeishuMsg{
	Title: "智能选择示例",
	Markdown: map[string]any{}, // 空的 map
	MarkdownItems: []bot.Text{}, // 空的切片
	MarkdownArray: [][2]string{
		{"自动使用", "MarkdownArray 内容"},
	},
}
// 系统会自动使用 MarkdownArray 的内容
```

---

## 🚀 飞书卡片2.0 新特性详解

### 1. 宽屏模式与转发控制

```go
msg := &bot.FeishuMsg{
	Title: "宽屏模式示例",
	MarkdownArray: [][2]string{
		{"特性", "宽屏模式展示"},
		{"说明", "适合展示更多内容"},
	},
	WideScreen:    true,  // 启用宽屏模式
	EnableForward: true,  // 允许转发
	HeaderColor:   bot.ColorTurquoise,
}
```

### 2. 交互按钮

```go
msg := &bot.FeishuMsg{
	Title: "带按钮的消息",
	MarkdownArray: [][2]string{
		{"通知类型", "重要提醒"},
		{"内容", "请点击下方按钮查看详情"},
	},
	Actions: []bot.Action{
		bot.CreatePrimaryButtonElement("查看详情", "https://www.feishu.cn"),
		bot.CreateButtonElement("取消", ""),
	},
	HeaderColor: bot.ColorOrange,
}
```

### 3. 自定义图标

```go
msg := &bot.FeishuMsg{
	Title: "自定义图标示例",
	MarkdownArray: [][2]string{
		{"特性", "支持自定义头部图标"},
		{"说明", "使用图片token配置"},
	},
	CustomIcon: &bot.Icon{
		Tag:   "standard_icon",
		Token: "bell_outlined", // 飞书内置图标token
	},
	HeaderColor: bot.ColorRed,
}
```

### 4. 图片嵌入

```go
msg := &bot.FeishuMsg{
	Title: "带图片的消息",
	MarkdownArray: [][2]string{
		{"说明", "查看下方图片"},
	},
	Images: []string{
		"img_v3_025h_xxxx", // 图片token，需要先上传图片获取
	},
	HeaderColor: bot.ColorBlue,
}
```

### 5. 多列布局

```go
// 创建多列布局元素
columns := []bot.Column{
	bot.CreateColumn("top", "第一列内容"),
	bot.CreateColumn("top", "第二列内容"),
}
columnSet := bot.CreateColumnSetElement(columns, "bisect")

// 可以在 Elements 中使用
```

### 6. 带确认弹窗的按钮

```go
action := bot.Action{
	Tag: "button",
	Text: &bot.Text{
		Content: "删除",
		Tag:     "plain_text",
	},
	Type: "danger",
	Confirm: &bot.Confirm{
		Title: bot.Text{
			Content: "确认删除？",
			Tag:     "plain_text",
		},
		Text: bot.Text{
			Content: "删除后不可恢复",
			Tag:     "plain_text",
		},
	},
}

msg.Actions = []bot.Action{action}
```

### 7. 国际化支持

```go
msg := &bot.FeishuMsg{
	Title: "国际化示例",
	// ... 其他配置
}

// 在 FormatMsg 后可以手动添加国际化配置
card := bot.FormatMsg(msg)
card.Card.I18nElements = &bot.I18nElements{
	ZhCn: &bot.I18nElement{
		Elements: []bot.Element{
			bot.CreateMarkdownElement("中文内容"),
		},
	},
	EnUs: &bot.I18nElement{
		Elements: []bot.Element{
			bot.CreateMarkdownElement("English content"),
		},
	},
}
```

---

## 使用场景推荐

- **MarkdownArray**: 简单键值对场景（推荐）
- **MarkdownItems**: 需要混合键值对和纯内容的复杂场景
- **Markdown**: 兼容现有代码，但顺序不固定
- **卡片2.0交互组件**: 需要用户操作的场景（审批、确认等）
- **宽屏模式**: 展示大量数据或复杂表格
- **自定义图标**: 品牌化或分类标识

---

## 文档与资源

- [飞书机器人官方文档](https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot)
- [飞书卡片2.0文档](https://open.feishu.cn/document/feishu-cards/card-json-v2-structure)
- [卡片搭建工具](https://open.feishu.cn/cardkit)
- [Markdown语法说明](https://open.feishu.cn/document/uAjLw4CM/ukzMukzMukzM/feishu-cards/card-components/content-components/rich-text)
- [交互组件文档](https://open.feishu.cn/document/uAjLw4CM/ukzMukzMukzM/feishu-cards/card-components/interactive-components)

---

## 许可证

MIT License
