# feishu-bot-markdown
使用飞书群机器人发送markdown消息, 使用飞书群机器人卡片结构, 发送markdown格式, 本项目是对飞书群机器人发送消息的封装, 使得发送消息更加方便



# 安装
`go get github.com/crazykun/feishu-bot-markdown`


# 使用
bot_test.go

```
import (
	bot "github.com/crazykun/feishu-bot-markdown"
)

// 替换自己的机器人地址
var Hook = "https://open.feishu.cn/open-apis/bot/v2/hook/xxxx-xxxx-xxxx-xxxx"

var TmplTestFeishu = &bot.FeishuMsg{
	Title: "飞书markdown消息测试！",
	Markdown: map[string]any{
		"标题": "标题测试",
		"内容": "这是内容**粗体**, *斜体*, ~~删除线~~",
		"状态": "<font color='green'>成功</font> <font color='red'>失败</font> <font color='grey'>灰色</font>",
		"列表": "\n- 列表1\n- 列表2\n- 列表3",
		"链接": "[飞书机器人助手](https://www.feishu.cn/hc/zh-CN/articles/236028437163-%E6%9C%BA%E5%99%A8%E4%BA%BA%E6%B6%88%E6%81%AF%E5%86%85%E5%AE%B9%E6%94%AF%E6%8C%81%E7%9A%84%E6%96%87%E6%9C%AC%E6%A0%B7%E5%BC%8F)",
		"时间": `2021-08-12 12:00:00`,
	},
	Note:        "这是备注",
	Link:        "http://www.baidu.com",
	HeaderColor: bot.ColorWathet,
}

bot.SendFeishuMsg(Hook, TmplTestFeishu)

```

# 截图
![截图](https://raw.githubusercontent.com/crazykun/feishu-bot-markdown/main/src/screenshot.jpg)



# markdown支持格式
目前富文本组件仅支持 Markdown 语法的子集，详情参见下表。
 [文档](https://open.feishu.cn/document/uAjLw4CM/ukzMukzMukzM/feishu-cards/card-components/content-components/rich-text)

| 名称  | 语法 | 效果 | 注意事项 |
| ------------- | ------------- | ------------- | ------------- |
| 换行  |  `\n` | 文本换行  | 无  |
| 斜体  | `*斜体*`  |*斜体*  | 无  |
| 加粗  | `**粗体**` 或 `__粗体__`  |__粗体__ | 无  |
| 删除线 | `~~删除线~~` |~~删除线~~ | 无  |
| @指定人 | `<at id=user_id>张三</at>` | @张三 | 自定义机器人仅支持使用 open_id、user_id @指定人。  |
| 链接 | `<a href='https://open.feishu.cn'></a>` | [@飞书](https://open.feishu.cn) | 无  |
| 彩色文本 | `<font color='green'>绿色</font><font color='red'>红色</font><font color='grey'>灰色</font>` | <font color='green'>绿色</font><br><font color='red'>红色</font><br><font color='grey'>灰色</font> | 无  |
| 列表 | `1. 有序列表 - 无序列表` | - 无序列表1<br>  - 无序列表 1.1<br>- 无序列表2 | 无  |
| 代码块 | ` ``` fmt.Println("Hello World") ``` ` | ``` fmt.Println("Hello World") ``` | 无  |




# 飞书文档
[文档](https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot)

[卡片搭建](https://open.feishu.cn/cardkit)
