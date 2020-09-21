# 数据类型说明

```protobuf
syntax = "proto3";

package protobuf;

message Request {
  int32  type = 1;
  string path = 2;
  string location = 3;
  bytes data = 4;
  bytes sign = 5;
}

message Response {
  int32 code = 1;
  string msg = 2;
  int32 type = 3;
  bytes data = 4;
  bool render = 5;
  Render html = 6;
}

message Render {
  string code = 1;
  int32 type = 2;
  string id = 3;
  bool iframe = 4;
}
```

## Requests

### `type` 消息类型

-   **0** 为 系统消息
    -   此类型下无强制字段要求
    -   此类型下可能传输心跳包

-   **1** 为 请求 (发起一个类似于 HTTP 的请求)
    -   此类型下必须包含 path

-   **2** 为 binary
    -   此类型下可能传送文件等
    -   此类型必须包含 `path` 用于鉴别文件
    -   此类型下必须包含 `sign` 用于检验消息的完整性
    
-   **3** 为 Table.js 的请求渲染

### `location` 请求页面的路径

发起请求的页面所有的路径

### `path` 请求的路由信息

此字段仅在 `type` 为 1 时发送

### `data` 请求的数据

此字段传输请求的数据, 类 HTTP Request Body

### `sign` 检验签名

计算方法为 **salt** + `data`

## Response

### `code` 状态码

此为错误状态码, 0 代表成功

### `msg` 消息

对此次响应的说明信息

### `type` 消息类型

-   **0** 为指令
    -   通常命令客户端, 如跳出验证码, 或强制断开客户端连接, 课堂签到
    -   在此类型下, 一般包含 `render` 用来渲染前端
-   **1** 为 `salt`
    -   服务端向客户端发送一个 `salt` 用于计算 `salt`
    -   通常会在客户端连接服务端时主动发送
    -   此消息类型时 `render` 不会存在消息
-   **2** 为 响应 (一个类似于 HTTP 的响应)
    -   此消息类型时 **不一定** 会包含 `render`
    -   此时 `data` 常为 `json` 类型, 通常用于 `table` 的数据渲染
-   **3** 为 响应 (一个类型于 HTTP 的响应, 且带前端渲染)
    -   此消息类型与 2 相似, 但一定包含了 `render` 来使前端渲染
-   **4** 为 binary
    -   此类型时, 可能是由服务器端主动向客户端发送文件

### `render` 是否渲染

为 `true` 时渲染页面

### `html` 前端渲染代码

`id`  为渲染部分的元素 ID
`iframe` 为是否为元素 ID 是否存在在 iframe 里
`code` 为 html 代码
`type` 为渲染位置

-   **0** 为 `body`
    -   渲染前端的 body 位置, 通常为页面渲染
-   **1** 为 `div`
    -   这个是渲染 `<div id="1"></div>` 的
    -   一般为消息框或者为弹出层
