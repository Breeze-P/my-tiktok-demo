# Kitex version of TikTok Demo

### Tech Stack

Hertz、Kitex、Redis、MySQL、RabbitMQ、MongoDB、MinIO、etcd、Jaeger

`docker-compose up` for necessary dependencies

### Desc
1. 项目采用微服务架构，分离了用户管理与视频流服务，设置独立网关，可以分布式部署
2. 基于Bcrypt加密算法进行密码存储，利用全局拦截器进行JWT权限校验
3. 使用etcd作为注册中心，实现服务注册与发现功能，使用Jaeger监控服务状况
4. 使用Redis缓存响应数据，加速数据访问；存储鉴权信息，避免重复登录
5. 使用RabbitMQ缓存视频发布请求，优化发布流程；缓存聊天内容，顺序发送消息
6. 使用MongoDB实现聊天功能，分治消息的私聊与群聊、用户的在线与离线
7. 配套实现了RBAC后台管理系统，方便服务数据、用户权限的管理

附：<a href="https://github.com/Breeze-P/my-tiktok-demo/blob/main/DEVLOG.md">开发日志</a>