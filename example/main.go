package main

import (
	"fmt"

	bot "github.com/crazykun/feishu-bot-markdown"
)

func main() {
	// 示例1: 使用传统的 map 方式（可能无序）
	msgWithMap := &bot.FeishuMsg{
		Title: "使用 Map 的消息",
		Markdown: map[string]interface{}{
			"状态": "运行中",
			"时间": "2024-01-01 12:00:00",
			"结果": "成功",
		},
		Note:        "使用 map 可能导致顺序不固定",
		HeaderColor: bot.ColorBlue,
	}

	// 示例2: 使用新的 MarkdownItems 方式（保持顺序）
	msgWithItems := &bot.FeishuMsg{
		Title: "使用 MarkdownItems 的消息",
		MarkdownItems: []bot.Text{
			{Tag: "第一步", Content: "初始化系统"},
			{Tag: "第二步", Content: "加载配置文件"},
			{Tag: "第三步", Content: "启动服务"},
			{Tag: "状态", Content: "<font color='green'>成功</font>"},
		},
		Note:        "使用 MarkdownItems 保持固定顺序",
		HeaderColor: bot.ColorGreen,
	}

	// 示例3: 使用 MarkdownArray 方式（最简洁）
	msgWithArray := &bot.FeishuMsg{
		Title: "使用 MarkdownArray 的消息",
		MarkdownArray: [][2]string{
			{"状态", "<font color='green'>成功</font>"},
			{"进度", "100%"},
			{"耗时", "2.5秒"},
			{"结果", "任务完成"},
		},
		Note:        "使用 MarkdownArray 最简洁的方式",
		HeaderColor: bot.ColorViolet,
	}

	// 示例4: 当 Markdown 为空时，自动使用 MarkdownItems
	msgFallback := &bot.FeishuMsg{
		Title:    "回退机制示例",
		Markdown: map[string]interface{}{}, // 空的 map
		MarkdownItems: []bot.Text{
			{Tag: "自动回退", Content: "当 Markdown 为空时使用此内容"},
			{Tag: "优势", Content: "解决了 JSON 序列化时 map 无序的问题"},
		},
		Note:        "演示回退机制",
		HeaderColor: bot.ColorYellow,
	}

	// 生成消息卡片来查看结构
	fmt.Println("=== 使用 Map 的消息 ===")
	card1 := bot.FormatMsg(msgWithMap)
	fmt.Printf("标题: %s\n", card1.Card.Header.Title.Content)

	fmt.Println("=== 使用 MarkdownItems 的消息 ===")
	card2 := bot.FormatMsg(msgWithItems)
	fmt.Printf("标题: %s\n", card2.Card.Header.Title.Content)

	fmt.Println("=== 使用 MarkdownArray 的消息 ===")
	card3 := bot.FormatMsg(msgWithArray)
	fmt.Printf("标题: %s\n", card3.Card.Header.Title.Content)

	fmt.Println("=== 回退机制的消息 ===")
	card4 := bot.FormatMsg(msgFallback)
	fmt.Printf("标题: %s\n", card4.Card.Header.Title.Content)

	// 如果需要发送消息，取消注释以下代码并设置正确的 webhook URL
	/*
		webhook := "https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url"

		if err := bot.SendFeishuMsg(webhook, msgWithArray); err != nil {
			log.Printf("发送消息失败: %v", err)
		} else {
			log.Println("消息发送成功")
		}
	*/
}
