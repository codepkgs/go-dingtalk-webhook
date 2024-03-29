# 说明

> `Go` 版本的钉钉 `Webhook`机器人`SDK`

# 功能列表

## 支持的消息类型

[自定义机器人接入](https://open.dingtalk.com/document/orgapp/custom-robot-access#title-jfe-yo9-jl2)

- [x] 普通文本消息 `client.Text`
- [x] Markdown 消息 `client.Markdown`
- [x] 链接（Link）消息 `client.Link`
- [x] ActionCard 消息 `client.ActionCard`
- [x] FeedCard 消息 `client.FeedCard`

# 示例

- 初始化 Client

  > WebhookAddress 为创建机器人时产生的 Webhook 地址。
  > 如果创建的机器人的安全设置采用的是 自定义关键词 或 IP 地址(段)，在创建 client 的时候，`Secret` 传入空字符串即可。
  > 如果创建的机器人的安全设置采用的是 加签，在创建 client 的时候，`Secret` 传入产生的密钥即可。

  ```go
  // 初始化一个未采用加签的机器人
  client, err := dingtalk.NewClient(
      "https://oapi.dingtalk.com/robot/send?access_token=xxxxx",
      "")
  if err != nil {
      fmt.Println(err)
  }
  ```

  ```go
  // 初始化一个采用加签的机器人
  client, err := dingtalk.NewClient(
      "https://oapi.dingtalk.com/robot/send?access_token=xxxxx",
      "SECxxxxx")
  if err != nil {
      fmt.Println(err)
  }
  ```

- 发送文本消息

  > `atMobiles` 和 `isAtAll`: 如果 `isAtAll` 为 true，则会 at 所有人，否则只 at 在 atMobiles 中的用户。

  ```go
  sr, err := client.Text("测试", []string{"18611111111"}, false)
  if err != nil {
      fmt.Println(err)
  } else {
      fmt.Printf("%#v", sr)
  }
  ```

- 发送 Markdown 消息

  ```go
  sr, err := client.Markdown("测试", fmt.Sprintf(`
  # 一级标题
  * 测试消息1
  * 测试消息2
  * <font color="red">测试消息3</font>
  %s
  ## 二级标题
  * [Link](https://www.baidu.com)
  %[1]s
  `, strings.Repeat("-", 30)), []string{"18611111111"}, false)
  if err != nil {
      fmt.Println(err)
  } else {
      fmt.Printf("%#v", sr)
  }
  ```

- 发送 Link 消息

  ```go
  sr, err := client.Link("测试消息", "这是一个测试的Link类型的消息", "https://www.baidu.com", "https://blog.itpub.net/ueditor/php/upload/image/20200211/1581400086713823.png")
  if err != nil {
      fmt.Println(err)
  } else {
      fmt.Printf("%#v", sr)
  }
  ```

- 发送独立跳转的 ActionCard 类型的消息

  注：`content` 支持 `markdown` 格式的消息。

  > `dingtalk.Vertical`：按钮垂直排列
  > `dingtalk.Horizontal`：按钮水平排列
  > `ActionType: dingtalk.WEB`：表示在单独的浏览器中打开（电脑端 APP 应用）
  > `ActionType: dingtalk.APP`：表示在 APP 侧边栏中打开（电脑端 APP 应用。如果没有指定 ActionType，则默认为在 app 侧边栏中打开）

  ```go
  sr, err := client.ActionCard("ActionCard消息", "这是一个测试的ActionCard类型的消息，按钮垂直排列", dingtalk.Vertical, []dingtalk.ActionCardButton{
      {Title: "跳转到百度", ActionURL: "https://www.baidu.com", ActionType: dingtalk.WEB},
      {Title: "跳转到京东", ActionURL: "https://www.jd.com", ActionType: dingtalk.APP},
  })
  if err != nil {
      fmt.Println(err)
  } else {
      fmt.Printf("%#v", sr)
  }
  ```

- 发送 FeedCard 类型的消息

  > `Title`：单条信息文本
  > `MessageURL`：点击单条信息到跳转链接
  > `PicURL`：单条信息后面图片的 URL
  > `ActionType: dingtalk.WEB`：表示在单独的浏览器中打开（电脑端 APP 应用）
  > `ActionType: dingtalk.APP`：表示在 APP 侧边栏中打开（电脑端 APP 应用。如果没有指定 ActionType，则默认为在 app 侧边栏中打开）

  ```go
  sr, err := client.FeedCard([]dingtalk.FeedCardLink{
      {Title: "跳转到百度", MessageURL: "https://www.baidu.com", PicURL: "https://blog.itpub.net/ueditor/php/upload/image/20200211/1581400086713823.png", ActionType: dingtalk.WEB},
      {Title: "跳转到京东", MessageURL: "https://www.jd.com", PicURL: "https://blog.itpub.net/ueditor/php/upload/image/20200211/1581400086713823.png"},
  })
  if err != nil {
      fmt.Println(err)
  } else {
      fmt.Printf("%#v", sr)
  }
  ```
