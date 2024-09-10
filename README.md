# crawlerctl
Crawler platform under development


```textmate

go-crawler-platform/
│
├── api/                        # API 层
│   ├── task.go                 # 任务管理相关 API
│   ├── log.go                  # 日志查看 API
│   ├── project.go              # 用户上传的爬虫项目的 API
│   ├── env.go                  # 环境管理的 API
│   └── schedule.go             # 定时任务相关 API
│
├── cmd/                        # 主程序入口
│   └── main.go                 # 程序入口，启动 HTTP 服务
│
├── config/                     # 配置文件
│   └── config.yaml             # 全局配置文件
│
├── crawler/                    # 核心任务调度与管理
│   ├── task_manager.go         # 任务管理器，启动/停止任务
│   ├── log_manager.go          # 日志收集与管理
│   ├── worker.go               # 任务执行器，支持 Python 脚本调度
│   ├── env_manager.go          # `pyenv` 环境管理
│   └── scheduler.go            # 定时任务调度器
│
├── db/                         # 数据库操作
│   ├── db.go                   # 数据库连接与模型操作
│   └── migrations/             # 数据库迁移文件
│
├── models/                     # 数据模型
│   ├── task.go                 # 任务模型定义
│   ├── project.go              # 用户上传的项目模型定义
│   └── schedule.go             # 定时任务模型定义
│
├── services/                   # 复杂业务逻辑处理层
│   ├── task_service.go         # 任务管理服务
│   ├── schedule_service.go     # 定时任务服务
│   └── log_service.go          # 日志管理服务
│
├── storage/                    # 数据存储层
│   ├── task_store.go           # 任务存储与状态更新
│   └── schedule_store.go       # 定时任务存储与管理
│
├── test/                       # 单元测试和集成测试
│   ├── test_task.go            # 针对任务调度、进程管理的测试
│   └── test_api.go             # 针对 API 的集成测试
│
├── Dockerfile                  # Docker 配置文件
├── docker-compose.yml          # Docker Compose 文件
├── Makefile                    # Makefile 文件
├── go.mod                      # Go module 依赖管理
├── go.sum                      # Go module 依赖文件
├── README.md                   # 项目说明文档
└── .gitignore                  # Git 忽略文件

```
