# 人人公众号

英文名：`powerchat`

`linux`/`windows` 加密聊天，分享文件，分享内容（markdown或html），TCP加密隧道

### Android 版主页

- 主页：[https://gitee.com/sonichy/PowerChat_Android](https://gitee.com/sonichy/PowerChat_Android)
- 下载地址：[https://gitee.com/sonichy/PowerChat_Android/releases](https://gitee.com/sonichy/PowerChat_Android/releases)

## 编译说明

发行版测试平台是 ubuntu18.04 和 debian9.8 ，其它平台如有异常可自行编译。

编译器：go编译器1.11以上、vala编译器0.40以上

编译图形界面(ui)需要 vala-0.40 或者 vala-0.42。
测试过 vala-0.38，编译后运行异常，其他版本未经测试。

编译依赖：glib2、gtk3、gee、json-glib

## 通用功能

1. 文字聊天。
2. 发送图片。
3. 发送文件。

## 特色功能

**个人网站**
用户只需运行了客户端，就可以对所有联系人提供网页服务。用户只需点击左侧`我的朋友`列表中的联系人名字右侧的`WEB`按钮，就可以打开对方提供的网页。你既可以使用默认功能提供个人博客，也可以花点功夫配置自己的特色网站。

### 一、默认个人网站
默认状态下，内置个人`Blog`服务，当访问客户打开你的主页时，会浏览`HOME/ChatShare`目录下的`.html`、`.htm`、`.md`三种类型的文件。可以用目录分类，并且会隐藏以`_`开头的目录，以便存放你的图片、其他文件、CSS、JS脚本等资源。

### 二、代理本地网站
真正懂得网站设计的朋友可以在本地配置个人网站，然后用本客户端做代理，供联系人访问，只需设置代理相应的端口。

### 三、TCP加密隧道
在两台客户端之间形成一个加密TCP隧道，可用于玩游戏、使用局域网工具。
