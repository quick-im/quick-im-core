# quick-im-core

![](./docs/quick-im-core-layer.png)

## 项目介绍

- qucik-im-core项目初衷不是打造一个完整的即时通信项目，而是提供一个易于集成、扩展且轻量级并高性能的实时消息模块。

- 本项目提供完整的HTTP API接入文档，可接入任何适用场景，包括且不限于*客服系统*、*OA系统*等依赖实时消息功能的产品中。

- 本项目适用Go语言开发，目标是*轻*、*快*、*稳*，并提供完整的链路追踪功能，以便于快速梳理模块间依赖关系以便对可能产生的问题进行快速定位以及排查。

- 本项目支持单用户*多协议*、*多平台*同时登录，目前集成消息网关协议有*websocket*、*sse*、*poll*，终端可根据场景以及习惯来自定选择不同协议进行接入。

- 本项目不对任何*类型*数据进行硬限制，包括且不限于*会话类型*、*消息类型*等，为接入此项目的产品提供最自由的方案进行接入。

## 接口文档

项目结构简单易懂、面向产品端的接口主要分为*消息网关*、*会话服务相关*、*消息服务相关*等模块。

### 网关服务（所有服务接口集成于网关中）

#### 核心接口（消息收、发）

- 消息推送服务

    ```
    method：ws、poll(http)、sse
    url: {gateway server}/notify
    params: 
        token: {用户端token}
        protoc: {指定协议：ws|poll|sse}
    MsgData:

    ```

- 发送消息
    ```
    method：POST
    url: {gateway server}/send_msg
    header:
        Authorization: Bearer {用户端token}
    Request:

    Response:

    ```


#### 会话相关

##### 公开接口：可供客户端直接

- 检查当前用户是否在会话中
    ```
    method：POST
    url: {gateway server}/check_joined_conversation
    header:
        Authorization: Bearer {用户端token}
    Request:

    Response:

    ```

- 获取当前用户加入的会话列表
    ```
    method：POST
    url: {gateway server}/get_joined_conversation
    header:
        Authorization: Bearer {用户端token}
    Request:

    Response:

    ```

- 获取会话信息
    ```
    method：POST
    url: {gateway server}/get_conversation_detail
    header:
        Authorization: Bearer {用户端token}
    Request:

    Response:

    ```

- 加入会话
    ```
    method：POST
    url: {gateway server}/join_conversation
    header:
        Authorization: Bearer {用户端token}
    Request:

    Response:

    ```

- 创建会话
    ```
    method：POST
    url: {gateway server}/create_conversation
    header:
        Authorization: Bearer {用户端token}
    Request:

    Response:

    ```

- 离开会话
    ```
    method：POST
    url: {gateway server}/leave_conversation
    header:
        Authorization: Bearer {用户端token}
    Request:

    Response:

    ```

- 获取会话的最后一条消息ID
    ```
    method：POST
    url: {gateway server}/get_lastone_msgid
    header:
        Authorization: Bearer {用户端token}
    Request:

    Response:

    ```

##### 内部接口：不建议直接公开给客户端

- 获取Token
    ```
    method：POST
    url: {gateway server}/inner/get_token
    Request:

    Response:

    ```

- 创建会话并初始化成员
    ```
    method：POST
    url: {gateway server}/inner/create_conversation
    Request:

    Response:

    ```

- 踢出会话成员
    ```
    method：POST
    url: {gateway server}/inner/kickout_conversation
    Request:

    Response:

    ```

- 指定用户加入会话
    ```
    method：POST
    url: {gateway server}/inner/join_conversation
    Request:

    Response:

    ```

- 会话状态设为归档
    ```
    method：POST
    url: {gateway server}/inner/set_archive_conversations
    Request:

    Response:

    ```

#### 消息服务

##### 公开接口：可供客户端直接

- 获取会话指定消息ID范围的信息
    ```
    method：POST
    url: {gateway server}/get_conversation_msg_by_range
    Request:

    Response:

    ```

- 获取指定消息ID之前的30条消息
    ```
    method：POST
    url: {gateway server}/get_conversation_30msg_before_the_id
    Request:

    Response:

    ```

- 获取指定消息ID之后的30条消息
    ```
    method：POST
    url: {gateway server}/get_conversation_30msg_after_the_id
    Request:

    Response:

    ```

- 获取会话的最后30条消息
    ```
    method：POST
    url: {gateway server}/get_conversation_last_30_msg
    Request:

    Response:

    ```

- 获取会话之后一条消息
    ```
    method：POST
    url: {gateway server}/get_conversation_lastone_msg
    Request:

    Response:

    ```

##### 内部接口：不建议直接公开给客户端