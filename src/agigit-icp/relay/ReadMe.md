## 交互架构图

https://metalanguage.notion.site/Nostr-AGI-Git-NAG-Protocol-5646934339884d908877a508afab2011

## 安装需求

- Go 环境
- PostgreSQL 数据库

## 启动方法

```bash
go mod tidy
cd src/
go build -o relayer-basic
输入自己的数据库账号密码
POSTGRESQL_DATABASE=postgres://name:pass@localhost:5432/dbname ./relayer-basic  
```