# WebHooks
一个基于beego的github WebHooks服务。

# 安装
git clone git@github.com:yuanzhangcai/WebHooks.git

# 配置
配置app.conf文件主要配置：
1、配置github webhook secret
    格式：repository name + Secret = github webhook secret
    例： repository name为 WebHooks 配置为 WebHooksSecret = webhook_secret

2、配置执行脚本
    格式：repository name + Shell = 脚本文件名
    例： repository name为 WebHooks 配置为 WebHooksShell = gitPull.sh

其他log配置，脚本路径配置就不说了。

# 测试
1、进行Webhooks目录
2、bee run
3、请求url: http://localhost:9990/hook/

