本项目是为了能和网关连接起来，然后连接环界云的接口，以便为用户提供更多的服务。


2.ai-proxy 负责解析请求，并转发到 OpenAI 的 API

3.OpenAI 返回数据后，ai-proxy 再把数据返回给前端

4.网关可以把多个 AIproxy 服务连接起来，让流量更稳定

API 代理（handler/）：接收用户请求，调用 service/ 处理。


AI 调用（service/）：转发请求到 OpenAI、Claude、DeepSeek。

数据库存储（repository/）：存储请求日志和用户信息。

微服务部署（Dockerfile + deployment.yaml）。
