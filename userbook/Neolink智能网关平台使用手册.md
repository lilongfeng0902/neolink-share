# Neolink智能网关平台使用手册

## 1. 概览

__Neolink 智能网关平台__是一个面向大模型应用场景的 __统一 AI 服务转发与管理平台__，主要用于连接与分发多个主流大模型生态，包括 __ChatGPT、Claude、Gemini__ 等。通过 Neolink，开发者可以以统一接口调用不同模型，简化接入成本，提升调用效率，并实现更灵活的模型选择。

- 模型调用Base URL：http://120.133.40.59/api
- 模型调用API Key：xxxxxxxx

## 2. 模型API说明

### 2.1. 模型列表

现提供如下模型列表，可根据后续章节的使用示例来调用不同模型。

| 模型Id | 支持协议 | 说明 |
| --- | --- | --- |
| gpt-4.1 | OpenAI | 具有强推理、文本理解与生成等能力 |
| gpt-5 | OpenAI | 具有强推理（长上下文）、多模态等能力 |
| gpt-5.2 | OpenAI | 具有最强推理、多步规划等能力 |
| gpt-5.2-search | OpenAI | 默认开启搜索能力，以便更好的响应用户请求 |
| gpt-5.2-pro | OpenAI | 具有支持长上下文多轮对话的能力，对复杂指令、业务规则理解能力强 |
| gpt-oss-120b | OpenAI | 具有支持可配置的推理深度、完整的思维链访问，以及原生工具调用能力 |
| sora-2-guan | OpenAI | 具有参照文本提示生成视频或者文本提示+图片生成视频的能力 |
| sora-2-pro-guan | OpenAI | 具有参照文本提示生成视频或者文本提示+图片生成视频的能力 |
| claude-sonnet-4-5-20250929 | OpenAI、Claude | 建议使用Claude协议 |
| claude-opus-4-5-20251101 | OpenAI、Claude | 建议使用Claude协议 |
| claude-haiku-4-5-20251001 | OpenAI、Claude | 建议使用Claude协议 |
| claude-sonnet-4-5-20250929-thinking | OpenAI、Claude | 具有复杂推理、长逻辑链，文档理解、系统设计能力，建议使用Claude协议 |
| grok-code-fast-1 | OpenAI | 具有代码生成 / 补全 / 重构，响应速度快（fast），适合高并发的能力 |
| gemini-2.5-flash-image | OpenAI、Gemini | 具有文生图等能力，不适合复杂文本推理，建议使用Gemini协议 |
| gemini-3-pro-preview | OpenAI、Gemini | 具有强推力、多模态理解等能力，建议使用Gemini协议 |
| gemini-3-pro-preview-search | OpenAI、Gemini | 默认开启搜索能力，以便更好的响应用户请求，建议使用Gemini协议 |
| gemini-3-pro-image-preview | OpenAI、Gemini | 具有文生图、图像修改等能力，建议使用Gemini协议 |
| gemini-3-flash-preview | OpenAI、Gemini | 支持文本 + 图片理解，速度快成本低，建议使用Gemini协议 |
| Doubao-1.5-thinking-pro | OpenAI | 具有中文复杂推理能力，对政策、业务逻辑、长问题回答友好 |
| qwen3-max | OpenAI | 具有推理、生成、总结等能力 |
| qwen3-coder-plus | OpenAI | 具有代码理解 / 生成 / 修复的能力，对 Java / Go / Python / 前端支持较好 |
| deepseek-v3.2 | OpenAI | 具有数学、逻辑、代码推理能力，成本相对低，性价比极强 |
| GLM-4.6 | OpenAI | 具有中文理解、表达自然等能力，对政企、知识问答友好 |
| kimi-k2-0905-preview | OpenAI | 具有超长上下文，阅读、总结长文档等能力，推理能力较稳，不追求极限速度 |
| moonshot-v1-128k-vision-preview | OpenAI | 具有超长上下文，支持理解文本 + 图片等能力 |

说明：上面表格“支持协议”是指各家大模型厂商在接入方式（API）、消息格式（roles/parts）、流式输出、工具调用（function/tool calling）、鉴权与计费、以及安全/内容政策等方面的“接口规范与使用约定”。

### 2.2. 对话模型

#### 2.2.1. 创建对话请求（OpenAI）

1. curl请求示例

```bash
curl --location --request POST 'http://120.133.40.59/api/v1/chat/completions' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer <API Key>' \\
--data-raw '{
"model": "<Model ID>",
"messages": [
{
"role": "user",
"content": "hello"
}
]
}'
```

1. python编码示例

```python
import requests
url = "http://120.133.40.59/api/v1/chat/completions"

headers = {
"Content-Type": "application/json",
"Authorization": "Bearer <API Key>",
}

payload = {
"model": "<Model ID>",
"messages": [
{
"role": "user",
"content": "hello"
}
]
}
response = requests.post(url, headers=headers, json=payload, timeout=30)
response.raise_for_status()
result = response.json()
print(result)
if "choices" in result:
print("Assistant:", result["choices"][0]["message"]["content"])
```

1. 支持模型Id

- gpt-4.1
- gpt-5
- gpt-5.2
- gpt-oss-120b
- gemini-3-pro-preview
- Doubao-1.5-thinking-pro
- qwen3-max
- qwen3-coder-plus
- deepseek-v3.2
- GLM-4.6
- kimi-k2-0905-preview
- moonshot-v1-128k-vision-preview

#### 2.2.2. 创建对话请求（Anthropic）

1. curl请求示例

```bash
curl --location --request POST 'http://120.133.40.59/api/v1/messages' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer <API Key>' \\
--data-raw '{
"model": "<Model ID>",
"messages": [
{
"role": "user",
"content": "hello"
}
]
}'
```

1. python编码示例

```python
import requests

url = "http://120.133.40.59/api/v1/messages"

headers = {
"Authorization": "Bearer <API Key>",
"Content-Type": "application/json",
}

payload = {
"model": "<Model ID>",
"messages": [
{
"role": "user",
"content": "hello"
}
]
}

resp = requests.post(url, headers=headers, json=payload, timeout=30)
resp.raise_for_status()
data = resp.json()

reply = "".join(
item["text"]
for item in data.get("content", [])
if item.get("type") == "text"
)

print("Assistant reply:")
print(reply)
```

1. 支持模型Id

- claude-sonnet-4-5-20250929
- claude-opus-4-5-20251101
- claude-haiku-4-5-20251001
- claude-sonnet-4-5-20250929-thinking

#### 2.2.3. 创建对话请求（Google）

1. curl请求示例

```bash
curl --location --request POST 'http://120.133.40.59/api/v1beta/models/<Model ID>:generateContent' \\
--header 'Authorization: Bearer <API Key>' \\
--header 'Content-Type: application/json' \\
--data-raw '{
"contents": [
{
"parts": [
{
"text": "How does AI work?"
}
]
}
]
}'

注：请求header中默认使用Authorization传递鉴权参数，如：--header 'Authorization: Bearer <API Key>'，如遇到不能识别Authorization参数的情形，可使用x-goog-api-key传递鉴权信息进行处理，如：--header 'x-goog-api-key: <API Key>'。
```

1. python编码示例

```python
import requests

url = "http://120.133.40.59/api/v1beta/models/<Model ID>:generateContent"

headers = {
"x-goog-api-key": "<API Key>",
"Content-Type": "application/json",
}

payload = {
"contents": [
{
"parts": [
{
"text": "How does AI work?"
}
]
}
]
}

resp = requests.post(url, headers=headers, json=payload, timeout=30)
resp.raise_for_status()
data = resp.json()
print(data )

注：请求header中默认使用Authorization传递鉴权参数，如 "Authorization": "Bearer <API Key>"，如遇到不能识别Authorization参数的情形，可使用x-goog-api-key传递鉴权信息进行处理，将 "Authorization": "Bearer <API Key>"替换为"x-goog-api-key": "<API Key>"。
```

1. 支持模型Id

- gemini-2.5-flash-image
- gemini-3-pro-image-preview
- gemini-3-pro-preview
- gemini-3-flash-preview

### 2.3. 多模态模型

#### 2.3.1. 创建多模态请求（OpenAI）

1. curl请求示例

```bash
curl --location --request POST 'http://120.133.40.59/api/v1/responses' \\
--header 'Authorization: Bearer <API Key>' \\
--header 'Content-Type: application/json' \\
--data-raw '{
"model": "<Model ID>",
"input": [
{
"role": "user",
"content": [
{ "type": "input_text", "text": "描述这张图，并生成一张类似的" },
{
"type": "input_image",
"image_url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250925/thtclx/input1.png"
}
]
}
]
}'
```

1. python编码示例

```python
import requests
url = "http://120.133.40.59/api/v1/responses"

headers = {
"Content-Type": "application/json",
"Authorization": "Bearer <API Key>",
}

payload = {
"model": "<Model ID>",
"input": [
{
"role": "user",
"content": [
{ "type": "input_text", "text": "描述这张图，并生成一张类似的" },
{
"type": "input_image",
"image_url": "https://help-static-aliyun-doc.aliyuncs.com/file-manage-files/zh-CN/20250925/thtclx/input1.png"
}
]
}
]
}
response = requests.post(url, headers=headers, json=payload, timeout=30)
response.raise_for_status()
result = response.json()
print(result)
```

1. 支持模型Id

- gpt-4.1
- gpt-5
- gpt-5.2
- gpt-5.2-pro

#### 2.3.2. 搜索请求示例（OpenAI）

##### 2.3.2.1. 创建原生模型搜索请求

1. curl请求示例

```bash
curl --location --request POST 'http://120.133.40.59/api/v1/responses' \\
--header 'Authorization: Bearer <API Key>' \\
--header 'Content-Type: application/json' \\
--data-raw '{
"model": "<Model ID>",
"reasoning": {
"effort": "low"
},
"tools": [
{
"type": "web_search"
}
],
"tool_choice": "auto",
"input": "今天北京天气怎么样？"
}'
```

1. python编码示例

```python
SQL
import requests
url = "http://120.133.40.59/api/v1/responses"

headers = {
"Content-Type": "application/json",
"Authorization": "Bearer <API Key>",
}

payload = {
"model": "<Model ID>",
"reasoning": {
"effort": "low"
},
"tools": [
{
"type": "web_search"
}
],
"tool_choice": "auto",
"input": "今天北京天气怎么样？"
}
response = requests.post(url, headers=headers, json=payload, timeout=30)
response.raise_for_status()
result = response.json()
print(result)
```

1. 支持模型Id

- gpt-5
- gpt-5.2

##### 2.3.2.2. 创建默认搜索请求

1. curl请求示例

```bash
curl --location --request POST 'http://120.133.40.59/api/v1/responses' \\
--header 'Authorization: Bearer <API Key>' \\
--header 'Content-Type: application/json' \\
--data-raw '{
"model": "<Model ID>",
"input": "今天北京天气怎么样？"
}'
```

1. python编码示例

```python
SQL
import requests
url = "http://120.133.40.59/api/v1/responses"

headers = {
"Content-Type": "application/json",
"Authorization": "Bearer <API Key>",
}

payload = {
"model": "<Model ID>",
"input": "今天北京天气怎么样？"
}
response = requests.post(url, headers=headers, json=payload, timeout=30)
response.raise_for_status()
result = response.json()
print(result)
```

1. 支持模型Id

- gpt-5.2-search

#### 2.3.3. 搜索请求示例（Google）

##### 2.3.3.1. 创建原生模型搜索请求

1. curl请求示例

```bash
curl --location --request POST 'http://120.133.40.59/api/v1beta/models/<Model ID>:generateContent' \\
--header 'Authorization: Bearer <API Key>' \\
--header 'Content-Type: application/json' \\
--data-raw '{
"contents": [
{
"parts": [
{
"text": "北京今天天气怎么样？"
}
]
}
],
"tools": [
{
"googleSearch": {}
},
{
"urlContext": {}
}
]
}'
```

1. python编码示例

```python
SQL
import requests
url = "http://120.133.40.59/api/v1beta/models/<Model ID>:generateContent"

headers = {
"Content-Type": "application/json",
"Authorization": "Bearer <API Key>",
}

payload = {
"contents": [
{
"parts": [
{
"text": "北京今天天气怎么样？"
}
]
}
],
"tools": [
{
"googleSearch": {}
},
{
"urlContext": {}
}
]
}
response = requests.post(url, headers=headers, json=payload, timeout=30)
response.raise_for_status()
result = response.json()
print(result)
```

1. 支持模型Id

- gemini-3-pro-preview

##### 2.3.3.2. 创建默认搜索请求

1. curl请求示例

```bash
curl --location --request POST 'http://120.133.40.59/api/v1beta/models/<Model ID>:generateContent' \\
--header 'Authorization: Bearer <API Key>' \\
--header 'Content-Type: application/json' \\
--data-raw '{
"contents": [
{
"parts": [
{
"text": "北京今天天气怎么样？"
}
]
}
]
}'
```

1. python编码示例

```python
SQL
import requests
url = "http://120.133.40.59/api/v1beta/models/<Model ID>:generateContent"

headers = {
"Content-Type": "application/json",
"Authorization": "Bearer <API Key>",
}

payload = {
"contents": [
{
"parts": [
{
"text": "北京今天天气怎么样？"
}
]
}
]
}
response = requests.post(url, headers=headers, json=payload, timeout=30)
response.raise_for_status()
result = response.json()
print(result)
```

1. 支持模型Id

- gemini-3-pro-preview-search

#### 2.3.4. 创建视频理解请求（Google）

1. curl请求示例

```bash
curl --location --request POST 'http://120.133.40.59/api/v1beta/models/<Model ID>:generateContent' \\
--header 'Authorization: Bearer <API Key>' \\
--header 'Content-Type: application/json' \\
--data-raw '{
"contents": [
{
"parts": [
{
"inline_data": {
"mime_type": "video/mp4",
"data": "<Base64 Code>"
}
},
{
"text": "请总结这个视频的主要内容，并列出三个关键时刻。"
}
]
}
]
}'

注：“Base64 Code”参数是将大模型识别的视频进行base64编码，参考脚本命令如下：

Bash
\#\!/usr/bin/env bash
set -euo pipefail
VIDEO_PATH="./demo.mp4"

if [[ "$(base64 --version 2>&1)" = \*"FreeBSD"\* ]]; then
B64FLAGS="--input"
else
B64FLAGS="-w0"
fi
base64 $B64FLAGS "$VIDEO_PATH" > base64.txt
```

1. python编码示例

```python
import requests
url = "http://120.133.40.59/api/v1beta/models/<Model ID>:generateContent"

headers = {
"Content-Type": "application/json",
"Authorization": "Bearer <API Key>",
}

payload = {
"contents": [
{
"parts": [
{
"inline_data": {
"mime_type": "video/mp4",
"data": "<Base64 Code>"
}
},
{
"text": "请总结这个视频的主要内容，并列出三个关键时刻。"
}
]
}
]
}

resp = requests.post(url, headers=headers, json=payload, timeout=30)
resp.raise_for_status()
data = resp.json()
print(data )

注：“Base64 Code”参数是将大模型识别的视频进行base64编码，参考脚本命令如下：

Bash
\#\!/usr/bin/env bash
set -euo pipefail
VIDEO_PATH="./demo.mp4"

if [[ "$(base64 --version 2>&1)" = \*"FreeBSD"\* ]]; then
B64FLAGS="--input"
else
B64FLAGS="-w0"
fi
base64 $B64FLAGS "$VIDEO_PATH" > base64.txt
```

1. 支持模型Id

- gemini-3-pro-preview
- gemini-3-flash-preview

#### 2.3.5. 创建视频

##### 2.3.5.1. 创建视频请求（OpenAI）

1. curl请求示例

```bash
curl --location --request POST 'http://120.133.40.59/api/v1/videos' \\
--header 'Authorization: Bearer <API Key>' \\
--form 'model="<Model ID>"' \\
--form 'prompt="A calico cat playing a piano on stage"' \\
--form 'seconds="4"' \\
--form 'size="720x1280"'
```

1. python编码示例

```python
import requests

url = "http://120.133.40.59/api/v1/videos"

headers = {
"Authorization": "Bearer <API Key>"
}

files = {
"model": (None, "<Model ID>"),
"prompt": (None, "A calico cat playing a piano on stage"),
"seconds": (None, "4"),
"size": (None, "720x1280"),
}

response = requests.post(url, headers=headers, files=files)

print(response.status_code)
print(response.text)
```

1. 支持模型Id

- sora-2-guan
- sora-2-pro-guan

##### 2.3.5.2. 获取视频内容请求（OpenAI）

1. curl请求示例

```bash
curl --location --request GET 'http://120.133.40.59/api/v1/videos/<video Id>/content?video_id=<video Id>' \\
--header 'Authorization: Bearer <API Key>' \\
--output output.mp4
```

1. python编码示例

```python
import requests

API_KEY = "<API Key>"
VIDEO_ID = "<video Id>"

url = f"http://120.133.40.59/api/v1/videos/{VIDEO_ID}/content"

headers = {
"Authorization": f"Bearer {API_KEY}"
}

params = {
"video_id": VIDEO_ID
}

response = requests.get(url, headers=headers, params=params, stream=True)
response.raise_for_status()

with open("output.mp4", "wb") as f:
for chunk in response.iter_content(chunk_size=8192):
if chunk:
f.write(chunk)

print("视频已保存为 output.mp4")

注：变量video Id是创建视频请求成功后返回的参数
```
