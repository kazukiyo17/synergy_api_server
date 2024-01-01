
## redis

### 登录验证相关
predix | description
--- | ---
`username` | 用户名，标记是否已注册
`token` | 用户登录标记

### 业务相关
predix | description
--- | ---
`scene` | 用户信息

## Docker
### Build
```bash
docker build -t synergy_api_server .
```

### Run
```bash
docker run -d -p 8080:8080 synergy_api_server
```