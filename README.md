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

### 公共响应体
- 所有http接口的公共响应体，在此统一解释，一下文档中不在一一说明
```
{
    "code": 0, // 1：code不为0时error字段返回错误信息，data字段不存在；2：code为0时error字段不存在，data返回响应数据
    "error": "xxx错误",
    "data": null  // code不为0时data内容不存在
}
```

### 网关服务（所有服务接口集成于网关中）

#### 核心接口（消息收、发）

- 消息推送服务
    ```json
    method：ws、poll(http)、sse
    url: {gateway server}/notify
    params: 
        token: {用户端token}
        protoc: {指定协议：ws|poll|sse}
    MsgData:

   {
        "MsgId": "EDCZL-SY32-243Y-VWQ",
        "ConversationID": "87ba7679-b682-47e7-8499-0385dda22b66",
        "FromSession": "test-client-session-id",
        "SendTime": "2023-10-22T22:59:04.069413421+08:00",
        "Status": 0,
        "Type": 1,
        "Content": "消息内容"
    }
    ```

- 发送消息
    ```json
    method：POST
    url: {gateway server}/send_msg
    header:
        Authorization: Bearer {用户端token}
    Request:
        {
            "conversation_id": "87ba7679-b682-47e7-8499-0385dda22b66",
            "status": 0,
            "type": 1,
            "content": "消息内容"
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": 
            {
                "MsgID": "EDCZM-UUHJ-243Y-VWQ"
            }
        }
    ```


#### 会话相关

##### 公开接口：可供客户端直接

- 检查当前用户是否在会话中
    ```json
    method：POST
    url: {gateway server}/check_joined_conversation
    header:
        Authorization: Bearer {用户端token}
    Request:
        {
            "conversation_id": "会话ID"
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data":  // code不为0时data内容不存在
            {
                "Joined": false // false or true
            }
        }
    ```

- 获取当前用户加入的会话列表
    ```json
    method：POST
    url: {gateway server}/get_joined_conversation
    header:
        Authorization: Bearer {用户端token}
    Request:
        {
            "SessionId": "用户ID"
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": 
            {
                "Conversations": [] // 会话ID列表
            }
        }
        
    ```

- 获取会话信息
    ```json
    method：POST
    url: {gateway server}/get_conversation_detail
    header:
        Authorization: Bearer {用户端token}
    Request:
        {
            "conversation_id": "会话ID"
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": 
            {
                "ConversationID":   "会话ID",                   // string
                "LastMsgID":        "EDCZR-GL2A-243Y-VWQ",      // *string // 最后一条消息的ID
                "LastSendTime":     "2023-10-22T22:59:04.069413421+08:00",               // *time.Time 最后一条消息的发送时间
                "IsDelete":         false,                      // bool 是否删除
                "ConversationType": 0,                          // int64 会话类型，创建是指定
                "LastSendSession":  "87ba7679-b682-47e7-8499-0385dda22b66", // *string 最后发送消息的用户sessionid
                "IsArchive":        false                       // bool 是否归档
            }
        }
        
    ```

- 加入会话
    ```json
    method：POST
    url: {gateway server}/join_conversation
    header:
        Authorization: Bearer {用户端token}
    Request:
        {
            "conversation_id": "会话ID"
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": 
            {
                "ConversationId": "会话ID"
            }
        }
        
    ```

- 创建会话
    ```json
    method：POST
    url: {gateway server}/create_conversation
    header:
        Authorization: Bearer {用户端token}
    Request:
        {
            "conversation_type": 0, // int64 会话类型
        }`
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": 
            {
                "ConversationId": "会话ID"
            }
        }
    ```

- 离开会话
    ```json
    method：POST
    url: {gateway server}/leave_conversation
    header:
        Authorization: Bearer {用户端token}
    Request:
        {
            "conversation_id": "会话ID"
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": 
            {
                "ConversationId": "会话ID"
            }
        }
    ```

- 获取会话的最后一条消息ID
    ```json
    method：POST
    url: {gateway server}/get_lastone_msgid
    header:
        Authorization: Bearer {用户端token}
    Request:
        {
            "conversation_id": "会话ID"
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": 
            {
                "MsgId": "消息ID"
            }
        }
        
    ```

##### 内部接口：不建议直接公开给客户端

- 获取Token
    ```json
    method：POST
    url: {gateway server}/inner/get_token
    Request:
        {
            "session": "50864896-8136-4a43-8a48-1d3325a7f78f", // 用户身份唯一ID，建议使用uuid
            "platform": 1   // 平台标识，多平台登录使用
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": "xxxxx" // token
        }
    ```

- 创建会话并初始化成员
    ```json
    method：POST
    url: {gateway server}/inner/create_conversation
    Request:
        {
            "conversation_type": 0, // 自定义类型会话类型
	        "sessions": [] // 初始化的用户sessionid列表
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": 
            {
                "ConversationId": "会话ID"
            }
        }
    ```

- 踢出会话成员
    ```json
    method：POST
    url: {gateway server}/inner/kickout_conversation
    Request:
        {
            "conversation_id": "87ba7679-b682-47e7-8499-0385dda22b66", // 操作的会话id
	        "sessions": [] // 要踢出会话的sessionid列表
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": 
            {
                "ConversationId": "会话ID"
            }
        }
    ```

- 指定用户加入会话
    ```json
    method：POST
    url: {gateway server}/inner/join_conversation
    Request:
        {
            "conversation_id": "87ba7679-b682-47e7-8499-0385dda22b66", // 操作的会话id
	        "sessions": [] // 要加入会话的sessionid列表
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": 
            {
                "ConversationId": "会话ID"
            }
        }
    ```

- 会话状态设为归档/非归档
    ```json
    method：POST
    url: {gateway server}/inner/set_archive_conversations
    Request:
        {
            "conversation_ids" [], // 要设置的会话id列表
	        "is_archive": true // true or false
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
        }
    ```

#### 消息服务

##### 公开接口：可供客户端直接

- 获取会话指定消息ID范围的信息
    ```json
    method：POST
    url: {gateway server}/get_conversation_msg_by_range
    Request:
        {
            "conversation_id": "", // 操作的会话ID
            "start_msg_id": "",  // 起始消息ID
            "end_msg_id": "", // 结束消息ID
            "desc": true, // true or false 升降序
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": {
                "Msg": [
                    {
                        "MsgId": "", // string 消息ID
                        "ConversationID": "", // string 会话ID
                        "FromSession": "", // string 发送sessionId
                        "SendTime": "2023-10-22T22:59:04.069413421+08:00", // time.Time 发送时间
                        "Status": 0, // int32 消息状态，与发送时一致
                        "Type":0"", // int32 消息类型标识，与发送时一致
                        "Content": "", // string 消息内容
                    },
                ]
            }
        }
    ```

- 获取指定消息ID之前的30条消息
    ```json
    method：POST
    url: {gateway server}/get_conversation_30msg_before_the_id
    Request:
        {
            "conversation_id": "", // 要操作的会话id
            "msg_id": "", // 指定消息id
            "desc": false, // 升降序
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": {
                "Msg": [
                    {
                        "MsgId": "", // string 消息ID
                        "ConversationID": "", // string 会话ID
                        "FromSession": "", // string 发送sessionId
                        "SendTime": "2023-10-22T22:59:04.069413421+08:00", // time.Time 发送时间
                        "Status": 0, // int32 消息状态，与发送时一致
                        "Type":0"", // int32 消息类型标识，与发送时一致
                        "Content": "", // string 消息内容
                    },
                ]
            }
        }
    ```

- 获取指定消息ID之后的30条消息
    ```json
    method：POST
    url: {gateway server}/get_conversation_30msg_after_the_id
    Request:
        {
            "conversation_id": "", // 要操作的会话id
            "msg_id": "", // 指定消息id
            "desc": false, // 升降序
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": {
                "Msg": [
                    {
                        "MsgId": "", // string 消息ID
                        "ConversationID": "", // string 会话ID
                        "FromSession": "", // string 发送sessionId
                        "SendTime": "2023-10-22T22:59:04.069413421+08:00", // time.Time 发送时间
                        "Status": 0, // int32 消息状态，与发送时一致
                        "Type":0"", // int32 消息类型标识，与发送时一致
                        "Content": "", // string 消息内容
                    },
                ]
            }
        }
    ```

- 获取会话的最后30条消息
    ```json
    method：POST
    url: {gateway server}/get_conversation_last_30_msg
    Request:
        {
            "conversation_id": "", // 要操作的会话id
            "desc": false, // 升降序
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": {
                "Msg": [
                    {
                        "MsgId": "", // string 消息ID
                        "ConversationID": "", // string 会话ID
                        "FromSession": "", // string 发送sessionId
                        "SendTime": "2023-10-22T22:59:04.069413421+08:00", // time.Time 发送时间
                        "Status": 0, // int32 消息状态，与发送时一致
                        "Type":0"", // int32 消息类型标识，与发送时一致
                        "Content": "", // string 消息内容
                    },
                ]
            }
        }
    ```

- 获取会话之后一条消息
    ```json
    method：POST
    url: {gateway server}/get_conversation_lastone_msg
    Request:
        {
            "conversation_id": "", // 要操作的会话id
            "desc": false, // 升降序
        }
    Response:
        {
            "code": 0,
            "error": "xxx错误", // code不为0时返回错误信息
            "data": {
                "Msg": [
                    {
                        "MsgId": "", // string 消息ID
                        "ConversationID": "", // string 会话ID
                        "FromSession": "", // string 发送sessionId
                        "SendTime": "2023-10-22T22:59:04.069413421+08:00", // time.Time 发送时间
                        "Status": 0, // int32 消息状态，与发送时一致
                        "Type":0"", // int32 消息类型标识，与发送时一致
                        "Content": "", // string 消息内容
                    },
                ]
            }
        }
    ```

##### 内部接口：不建议直接公开给客户端