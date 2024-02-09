到华为云函数创建python3.9的项目然后复制main.py的内容到华为云函数index.py里, 之后在触发器里创建APIG网关, 安全认证选None, 仅支持使用params的方式给参数, 不支持body, 示例: https://push.example.com/?key=100&msg=hello
