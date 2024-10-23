# proxyCrawler

## 功能介绍

实现定时爬取公开代理并进行验证，提供了API接口可以快速获取到符合条件的代理


## API文档
### 代理相关接口
#### 总览
| 接口      | 描述                       |
|---------| -------------------------- |
| proxies | 获取符合条件的所有代理     |
| proxy   | 随机获取一个符合条件的代理 |
| add     | 添加代理                   |
#### 获取符合条件的所有代理
**描述**
<p>根据get参数的值返回所有符合条件的代理，默认返回所有代理</p>

**请求**
```http request
GET /proxies?param1=value1&param2=value2 HTTP/1.1
Host: example.com:port
Authorization: <AccessToken>
```
| 参数名称 | 必填 | 描述         | 可选值                                            |
| -------- | ---- | ------------ | ------------------------------------------------- |
| p        | 否   | 指定代理协议 | http https vmess socks ...                        |
| c        | 否   | 指定所在国家 | CN US UK ...                                      |
| v        | 否   | 是否通过验证 | 0 (未通过验证) 1 (已通过验证 \|默认)              |
| t        | 否   | 指定返回格式 | full (返回所有字段 \| 默认) link (仅返回链接格式) |

**响应**

```
1. rtype=link
["http://114.*.*.101:8888","https://117.*.*.166:8089"]

2. rtype=full
{
  "CType": "ip3366",
  "Protocol": "http",
  "Host": "114.*.*.101",
  "Port": 8888,
  "Validated": true,
  "Latency": 186,
  "ValidateDate": "2024-06-12T17:45:38.4537801+08:00",
  "ToValidateDate": "2024-06-12T18:45:38.4537801+08:00",
  "ValidateFailedCnt": 0,
  "User": "",
  "Pass": "",
  "Link": "http://114.*.*.101:8888",
  "Country": "CN"
}

```
#### 随机获取一个符合条件的代理
**描述**

<p>根据get参数的值随机返回一个符合条件的代理，默认返回一个随机代理</p>

**请求**

```http request
GET /proxy?param1=value1&param2=value2 HTTP/1.1
Host: example.com:port
Authorization: <AccessToken>
```
| 参数名称 | 必填 | 描述         | 可选值                                            |
| -------- | ---- | ------------ | ------------------------------------------------- |
| p        | 否   | 指定代理协议 | http https vmess socks ...                        |
| c        | 否   | 指定所在国家 | CN US UK ...                                      |
| v        | 否   | 是否通过验证 | 0 (未通过验证) 1 (已通过验证 \|默认)              |
| t        | 否   | 指定返回格式 | full (返回所有字段 \| 默认) link (仅返回链接格式) |

**响应**

```
1. rtype=link
"http://114.*.*.101:8888"
2. rtype=full
[{
  "CType": "ip3366",
  "Protocol": "http",
  "Host": "114.*.*.101",
  "Port": 8888,
  "Validated": true,
  "Latency": 186,
  "ValidateDate": "2024-06-12T17:45:38.4537801+08:00",
  "ToValidateDate": "2024-06-12T18:45:38.4537801+08:00",
  "ValidateFailedCnt": 0,
  "User": "",
  "Pass": "",
  "Link": "http://114.*.*.101:8888",
  "Country": "CN"
}, {
  "CType": "ip3366",
  "Protocol": "https",
  "Host": "117.*.*.166",
  "Port": 8089,
  "Validated": false,
  "Latency": 99999,
  "ValidateDate": "2024-06-12T17:45:38.4851554+08:00",
  "ToValidateDate": "2024-06-12T18:45:38.4851554+08:00",
  "ValidateFailedCnt": 1,
  "User": "",
  "Pass": "",
  "Link": "https://117.*.*.166:8089",
  "Country": "CN"
}]
```
#### 通过链接添加代理
**描述**

<p>从互联网上传代理</p>
目前支持的格式：

- **xui** : `https?:.*:54321/?`
- **local** : `{protocol}://{user}:{pass}@{host}:{port}/?`

**请求**

```http request
GET /add?c_type=xui&link=http://example.com:port/ips.txt HTTP/1.1
Host: example.com:port
Authorization: <AccessToken>
```
| 参数名称 | 必填 | 描述           | 可选值                          |
| -------- | ---- | -------------- | ------------------------------- |
| c_type   | 是   | 指定fetcher    | xui<br/>local<br/>...<br/>      |
| link     | 是   | 需要上传的文件 | http://example.com:port/ips.txt |

**响应**

```
成功添加n个代理
```

#### 通过POST添加代理

**描述**

<p>通过POST请求添加代理</p>

目前支持的格式：

- **xui** : `https?:.*:54321/?`
- **local** : `{protocol}://{user}:{pass}@{host}:{port}/?`

**请求**

```http request
POST /add HTTP/1.1
Host: example.com:port
Authorization: <AccessToken>
Content-Type: application/json

{"type":"xui","proxies":["http://1.2.3.4:54321"]}
```

| 参数名称 | 必填 | 描述        | 可选值                     |
| -------- | ---- | ----------- | -------------------------- |
| type     | 是   | 指定crawler | xui<br/>local<br/>...<br/> |

**响应**

```
添加结果
新增: 200 更新: 10 删除: 0 跳过: 0 报错: 0 总数: 210
```

# 开发文档

## 数据库设计

1、AddedData

| 字段名称  | 数据类型 | 说明                  |
| --------- | -------- | --------------------- |
| c_type    | 字符串   | crawler_type          |
| data      | 字符串   | 可供crawler读取的数据 |
| validated | 布尔值   | crawler是否已读取     |

2、代理

| 字段名称            | 数据类型 | 说明                                                       |
| ------------------- | -------- | ---------------------------------------------------------- |
| c_type              | 字符串   | 这个代理来自哪个爬取器                                     |
| protocol            | 字符串   | 代理协议名称，一般为HTTP                                   |
| host                | 字符串   | 代理的IP地址                                               |
| port                | 整数     | 代理的端口号                                               |
| validated           | 布尔值   | 这个代理是否通过了验证，通过了验证表示当前代理可用         |
| latency             | 整数     | 延迟(单位毫秒)，表示上次验证所用的时间，越小则代理质量越好 |
| validate_date       | 时间戳   | 上一次进行验证的时间                                       |
| to_validate_date    | 时间戳   | 下一次进行验证的时间                                       |
| validate_failed_cnt | 整数     | 已经连续验证失败了多少次，会影响下一次验证的时间           |
| user                | 字符串   | 代理用户名                                                 |
| pass                | 字符串   | 代理密码                                                   |
| link                | 字符串   | 链接                                                       |

3、扫描器

| 字段名称         | 数据类型 | 说明                                                         |
| ---------------- | -------- | ------------------------------------------------------------ |
| type             | 字符串   | 爬取器类型                                                   |
| enable           | 布尔值   | 是否启用这个爬取器，被禁用的爬取器不会在之后被运行，但是其之前爬取的代理依然存在 |
| sum_proxies_cnt  | 整数     | 至今为止总共爬取到了多少个代理                               |
| last_proxies_cnt | 整数     | 上次爬取到了多少个代理                                       |
| last_fetch_date  | 时间戳   | 上次爬取的时间                                               |
| to_fetch_date    | 时间戳   | 计划下次爬取的时间                                           |

## 新增扫描器

1. crawler包下新建一个golang文件 tmpCrawler.go

   ```go
   package crawler
   ```

2. 新建一个结构体

   ```go
   package crawler
   type tmpCrawler struct{}
   ```

3. 实现crawl方法

   ````go
   package crawler
   import "ProxyPool/model"
   type tmpCrawler struct{}
   func (f tmpCrawler) crawl() []model.DBProxy{
       var result []model.DBProxy
           // 从互联网获取数据
       resp, err := utils.GetClient().Get(fmt.Sprintf("http://www.ip3366.net/free/?stype=%d&page=%d", 0, 1))
       // 业务逻辑...
       return result
   }
   ```
   ````

4. 注册到扫描器表
   打开crawler.Crawler.go, 在crawlerMap添加类型和结构体

    ```go
   var crawlerMap = map[string]Crawler{
   	"xui":    XUICrawler{},
   	"ip3366": ip3366Crawler{},
   	"local":  localCrawler{}, 
   	"tmp":  tmpCrawler{},
   }
    ```

## config.yml

```yaml
web:
  addr: ":8081" # web服务端口
  base_path: "/v1" # 前置路径，设置为/v1后需访问 /v1/proxies
  authorization: "your-token" # 认证头，需要携带该请求头才可以正常访问，否则会报403 默认 a57f1abe-c6df-4a9b-82ad-a29cf1304399
database:
  #配置MySQL连接参数
  username: "root"  #账号
  password: "123456" #密码
  host: "127.0.0.1" #数据库地址，可以是Ip或者域名
  port: 3306 #数据库端口
  Dbname: "crawler" #数据库名
  timeout: "10s" #连接超时，10秒
```

