## 项目概述

本项目包含三个主要组成部分：用户（User）、智能合约（Contract）和中继（Relay）。用户通过智能合约发起请求，中继处理这些请求并与用户进行交互。此包中的是中继模块

## 交互架构图

https://metalanguage.notion.site/Nostr-AGI-Git-NAG-Protocol-5646934339884d908877a508afab2011

## 安装需求

- Go 环境
- PostgreSQL 数据库

## 启动方法

```bash
go mod tidy
cd relay-service/
go build -o relayer-basic
输入自己的数据库账号密码
POSTGRESQL_DATABASE=postgres://name:pass@localhost:5432/dbname ./relayer-basic  
```
