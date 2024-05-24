# [free-gpt3.5-2api](https://github.com/aurorax-neo/free-gpt3.5-2api)



## 一、支持

#### 1.支持免登录chat2api

#### 2.支持账号chat2api（Authorization Bearer eyJhbGciOiJSUz***）

#### 3.支持账号ACCESS_TOKEN池（Authorization Bearer ac-***）

## 二、配置

#### 环境变量

```
LOG_LEVEL=info    # debug, info, warn, error
LOG_PATH=         # 日志文件路径，默认为空（不生成日志文件）
BIND=0.0.0.0      # 127.0.0.1
PORT=3040
ACCESS_TOKENS=    # 账号token池（eyJhbGci****，eyJhbGci****）
PROXY=            # http://127.0.0.1:7890,http://127.0.0.1:7890 已支持多个代理（英文 "," 分隔）
AUTHORIZATIONS=   # abc,bac (英文 "," 分隔)
BASE_URL=         # 默认：https://chat.openai.com
POOL_MAX_COUNT=64 # max number of connections to keep in the pool 默认：64
AUTH_ED=60       # expiration time for the authorization in seconds 默认：60
```

###### 也可使用与程序同目录下 `.env` 文件配置上述字段

##### 若要使用ACCESS_TOKENS内的账号，AUTHORIZATIONS字段内必须配置`ac-`开头的AUTHORIZATION并使用ac-***调用本程序，若ACCESS_TOKENS无可用账号则默认使用免登

## 三、部署

### 1.docker部署

##### 1 .创建文件夹

```
mkdir -p $PWD/free-gpt3.5-2api
```

##### 2.拉取镜像启动

```
docker run -itd  --name=free-gpt3.5-2api -p 9846:3040 ghcr.io/aurorax-neo/free-gpt3.5-2api
```

##### 3.更新容器

```
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock containrrr/watchtower -cR free-gpt3.5-2api --debug
```

### 2.Koyeb部署

###### 注意：`Regions`请选择支持`openai`免登的区域！！！现原生ip已不支持免登，请配置代理使用！！！

[![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)](https://app.koyeb.com/deploy?type=docker&name=free-gpt3-5-2api&region=par&ports=3040;http;/&image=ghcr.io/aurorax-neo/free-gpt3.5-2api)

## 四、接口

#### 1./v1/tokens

```
curl --location --request GET 'http://127.0.0.1:9846/v1/tokens' \
--header 'Authorization: Bearer abc'
```

返回示例说明：`count`为授权池中可用授权数，如果`count` 为 `0`请检查`ip`是否支持 `openai`

```
{
    "count": 0
}
```

#### 2./v1/accTokens

```
curl --location --request GET 'http://127.0.0.1:9846/v1/accTokens' \
--header 'Authorization: Bearer abc'
```

返回示例说明：`count`为授权池中可用授权数，如果`count` 为 `0`请检查`ip`是否支持 `openai`

```
{
    "count": 1,
    "canUseCount": 1
}
```

#### 3./v1/chat/completions

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

## 参考项目

- https://github.com/aurora-develop/aurora

- https://github.com/xqdoo00o/ChatGPT-to-API
