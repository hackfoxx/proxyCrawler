# API文档
## 代理相关接口
### 总览
| 接口      | 描述                       |
|---------| -------------------------- |
| proxies | 获取符合条件的所有代理     |
| proxy   | 随机获取一个符合条件的代理 |
| add     | 添加代理                   |
### 获取符合条件的所有代理
**描述**
<p>根据get参数的值返回所有符合条件的代理，默认返回所有代理</p>

**请求**
```http request
GET /proxies?param1=value1&param2=value2 HTTP/1.1
Host: example.com:port
Authorization: <AccessToken>
```
| 参数名称      | 必填 | 描述     | 可选值                                        |
|-----------|----|--------|--------------------------------------------|
| protocol  | 否  | 指定代理协议 | http<br/>https<br/>vmess<br/>socks<br/>... |
| country   | 否  | 指定所在国家 | CN<br/>US<br/>UK<br/>...                   |
| validated | 否  | 是否通过验证 | 0 (未通过验证)<br/>1 (已通过验证)<br/>               |
| rtype     | 否  | 指定返回格式 | full (返回所有字段)<br/>link (仅返回链接格式)<br/>      |

**响应**
```
1. rtype=link
["http://114.231.8.101:8888","https://117.69.236.166:8089"]

2. rtype=full
{
  "FetcherType": "ip3366",
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
### 随机获取一个符合条件的代理
**描述**

<p>根据get参数的值随机返回一个符合条件的代理，默认返回一个随机代理</p>

**请求**

```http request
GET /proxy?param1=value1&param2=value2 HTTP/1.1
Host: example.com:port
Authorization: <AccessToken>
```
| 参数名称      | 必填 | 描述     | 可选值                                        |
|-----------|----|--------|--------------------------------------------|
| protocol  | 否  | 指定代理协议 | http<br/>https<br/>vmess<br/>socks<br/>... |
| country   | 否  | 指定所在国家 | CN<br/>US<br/>UK<br/>...                   |
| validated | 否  | 是否通过验证 | 0 (未通过验证)<br/>1 (已通过验证)<br/>               |
| rtype     | 否  | 指定返回格式 | full (返回所有字段)<br/>link (仅返回链接格式)<br/>      |

**响应**

```
1. rtype=link
"http://114.*.*.101:8888"
2. rtype=full
[{
  "FetcherType": "ip3366",
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
  "Link": "http://114.231.8.101:8888",
  "Country": "CN"
}, {
  "FetcherType": "ip3366",
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
  "Link": "https://117.69.236.166:8089",
  "Country": "CN"
}]
```
### 添加代理
**描述**

<p>从互联网上传代理</p>
目前支持的格式：

- **xui** : `https?:.*:54321/?`
- **local** : `{protocol}://{user}:{pass}@{host}:{port}/?`

**请求**

```http request
GET /add?fetcher_type=xui&link=http://example.com:port/ips.txt HTTP/1.1
Host: example.com:port
Authorization: <AccessToken>
```
| 参数名称         | 必填 | 描述        | 可选值                             |
|--------------|----|-----------|---------------------------------|
| fetcher_type | 是  | 指定fetcher | xui<br/>local<br/>...<br/>      |
| link         | 是  | 需要上传的文件   | http://example.com:port/ips.txt |

**响应**

```
成功添加n个代理
```