# Go飞书机器人Markdown消息工具

## 项目简介

一个基于Go语言封装的飞书群机器人消息发送工具，支持Markdown格式消息的快速构建与发送，提供简洁的API和丰富的消息格式支持。

---

## 安装

```bash
go get github.com/crazykun/feishu-bot-markdown
```

---

## 快速开始

### 1. 配置机器人Hook地址

```go
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

#### 方式二：使用 MarkdownItems（推荐，保持顺序）
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
- [x] **新增**：自动回退机制，当 Markdown 为空时使用 MarkdownItems

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

使用 `MarkdownItems` 替代 `Markdown` map 来保证内容顺序：

```go
// 问题：map 在 JSON 序列化时顺序不固定
msg := &bot.FeishuMsg{
	Markdown: map[string]any{
		"第三步": "完成",
		"第一步": "开始", 
		"第二步": "进行中",
	},
}

// 解决方案：使用 MarkdownItems 保持固定顺序
msg := &bot.FeishuMsg{
	MarkdownItems: []bot.Text{
		{Tag: "第一步", Content: "开始"},
		{Tag: "第二步", Content: "进行中"},
		{Tag: "第三步", Content: "完成"},
	},
}
```

### 自动回退机制

当 `Markdown` 字段为空时，系统会自动使用 `MarkdownItems`：

```go
msg := &bot.FeishuMsg{
	Title: "自动回退示例",
	Markdown: map[string]any{}, // 空的 map
	MarkdownItems: []bot.Text{
		{Tag: "自动使用", Content: "MarkdownItems 内容"},
	},
}
// 系统会自动使用 MarkdownItems 的内容
```

---

## 文档与资源

- [飞书机器人官方文档](https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot)
- [卡片搭建指南](https://open.feishu.cn/cardkit)
- [Markdown语法说明](https://open.feishu.cn/document/uAjLw4CM/ukzMukzMukzM/feishu-cards/card-components/content-components/rich-text)
