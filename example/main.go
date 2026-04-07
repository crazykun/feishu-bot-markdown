package main

import (
	"fmt"
	"log"

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

	// 示例5: 卡片2.0 - 宽屏模式 + 允许转发
	msgWideScreen := &bot.FeishuMsg{
		Title: "卡片2.0 - 宽屏模式",
		MarkdownArray: [][2]string{
			{"功能", "宽屏模式展示"},
			{"说明", "适合展示更多内容"},
		},
		WideScreen:    true, // 启用宽屏模式
		EnableForward: true, // 允许转发
		Note:          "卡片2.0新特性：宽屏模式",
		HeaderColor:   bot.ColorTurquoise,
	}

	// 示例6: 卡片2.0 - 带按钮交互
	msgWithActions := &bot.FeishuMsg{
		Title: "卡片2.0 - 交互按钮",
		MarkdownArray: [][2]string{
			{"通知类型", "重要提醒"},
			{"内容", "请点击下方按钮查看详情"},
		},
		Actions: []bot.Action{
			bot.CreatePrimaryButtonElement("查看详情", "https://www.feishu.cn"),
			bot.CreateButtonElement("取消", ""),
		},
		Note:        "卡片2.0新特性：交互组件",
		HeaderColor: bot.ColorOrange,
	}

	// 示例7: 卡片2.0 - 自定义图标
	msgWithIcon := &bot.FeishuMsg{
		Title: "卡片2.0 - 自定义图标",
		MarkdownArray: [][2]string{
			{"特性", "支持自定义头部图标"},
			{"说明", "使用图片token配置"},
		},
		CustomIcon: &bot.Icon{
			Tag:   "standard_icon",
			Token: "bell_outlined", // 飞书内置图标token
		},
		Note:        "卡片2.0新特性：自定义图标",
		HeaderColor: bot.ColorRed,
	}

	// 示例8: 卡片2.0 - 多列布局
	msgWithColumns := &bot.FeishuMsg{
		Title: "卡片2.0 - 多列布局",
		MarkdownArray: [][2]string{
			{"部署环境", "生产环境"},
			{"部署状态", "<font color='green'>成功</font>"},
			{"部署时间", "2024-01-01 12:00:00"},
			{"耗时", "3分25秒"},
		},
		Note:        "卡片2.0支持更丰富的布局",
		HeaderColor: bot.ColorGreen,
	}

	// 生成消息卡片来查看结构
	fmt.Println("=== 使用 Map 的消息 ===")
	card1 := bot.FormatMsg(msgWithMap)
	fmt.Printf("标题: %s\n", card1.Card.Header.Title.Content)

	fmt.Println("\n=== 使用 MarkdownItems 的消息 ===")
	card2 := bot.FormatMsg(msgWithItems)
	fmt.Printf("标题: %s\n", card2.Card.Header.Title.Content)

	fmt.Println("\n=== 使用 MarkdownArray 的消息 ===")
	card3 := bot.FormatMsg(msgWithArray)
	fmt.Printf("标题: %s\n", card3.Card.Header.Title.Content)

	fmt.Println("\n=== 回退机制的消息 ===")
	card4 := bot.FormatMsg(msgFallback)
	fmt.Printf("标题: %s\n", card4.Card.Header.Title.Content)

	fmt.Println("\n=== 卡片2.0 - 宽屏模式 ===")
	card5 := bot.FormatMsg(msgWideScreen)
	fmt.Printf("标题: %s\n", card5.Card.Header.Title.Content)
	if card5.Card.Config != nil {
		fmt.Printf("宽屏模式: %v, 允许转发: %v\n", card5.Card.Config.WideScreenMode, card5.Card.Config.EnableForward)
	}

	fmt.Println("\n=== 卡片2.0 - 交互按钮 ===")
	card6 := bot.FormatMsg(msgWithActions)
	fmt.Printf("标题: %s\n", card6.Card.Header.Title.Content)
	fmt.Printf("按钮数量: %d\n", len(card6.Card.Elements))

	fmt.Println("\n=== 卡片2.0 - 自定义图标 ===")
	card7 := bot.FormatMsg(msgWithIcon)
	fmt.Printf("标题: %s\n", card7.Card.Header.Title.Content)
	if card7.Card.Header.UdIcon != nil {
		fmt.Printf("图标Token: %s\n", card7.Card.Header.UdIcon.Token)
	}

	fmt.Println("\n=== 卡片2.0 - 多列布局 ===")
	card8 := bot.FormatMsg(msgWithColumns)
	fmt.Printf("标题: %s\n", card8.Card.Header.Title.Content)

	// 如果需要发送消息，取消注释以下代码并设置正确的 webhook URL
	webhook := "https://open.feishu.cn/open-apis/bot/v2/hook/xxxx-xxxx-xxxx-xxxx"

	// 发送带交互按钮的消息
	if err := bot.SendFeishuMsg(webhook, msgWithIcon); err != nil {
		log.Printf("发送消息失败: %v", err)
	} else {
		log.Println("带按钮的消息发送成功")
	}
}
