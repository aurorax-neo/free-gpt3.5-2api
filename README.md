# [free-gpt3.5-2api](https://github.com/aurorax-neo/free-gpt3.5-2api)

## 接口

#### /v1/tokens

```
curl --location --request GET 'http://127.0.0.1:9846/v1/tokens' \
--header 'Authorization: Bearer abc'
```

返回示例说明：`count`为授权池中可用授权数，如果` count` 为 `0`请检查`ip`是否支持 `openai`

```
{
    "count": 0
}
```

#### /v1/chat/completions

###### 支持返回stream和json

```
http://<ip>:<port>/v1/chat/completions
```

##### 示例

```
curl http://127.0.0.1:9846
```

```
curl --location --request POST 'http://127.0.0.1:9846/v1/chat/completions' \
--header 'Authorization: Bearer abc' \
--header 'Content-Type: application/json' \
--data-raw '{
    "model": "gpt-3.5-turbo",
    "messages": [
        {
            "role": "user",
            "content": "西红柿炒钢丝球怎么做?"
        }
    ],
    "stream": false
}'
```

## 配置

### 环境变量

```
LOG_LEVEL=info    # debug, info, warn, error
BIND=0.0.0.0      # 127.0.0.1
PORT=3040
PROXY=            # http://127.0.0.1:7890,http://127.0.0.1:7890 已支持多个代理（使用英文 "," 分隔）
AUTHORIZATIONS=   # abc,bac (英文 , 分隔)
POOL_MAX_COUNT=64 # max number of connections to keep in the pool 默认：64
AUTH_ED=600       # expiration time for the authorization in seconds 默认：600
```

###### 也可使用与程序同目录下 `.env` 文件配置上述字段


### docker部署

##### 1 .创建文件夹

```
mkdir -p $PWD/free-gpt3.5-2api
```

##### 2.拉取镜像启动

```
docker run -itd  --name=free-gpt3.5-2api -p 9846:3040 -v $PWD/free-gpt3.5-2api/logs:/app/logs ghcr.io/aurorax-neo/free-gpt3.5-2api
```

### Koyeb部署

###### 注意：`Regions`请选择支持`openai`免登的区域！！！

[![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)](https://app.koyeb.com/deploy?type=docker&name=free-gpt3-5-2api&region=par&ports=3040;http;/&image=ghcr.io/aurorax-neo/free-gpt3.5-2api)

