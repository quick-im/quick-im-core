## 接口文档

项目结构简单易懂、面向产品端的接口主要分为*消息网关*、*会话服务相关*、*消息服务相关*等模块。

### 公共响应体
- 所有http接口的公共响应体，在此统一解释，一下文档中不在一一说明
```json
{
    "code": 0,
    "error": "xxx错误",
    "data": null
}
```
- `code`不为0时`error`字段返回错误信息，`data`字段不存在；`code`为0时`error`字段不存在，`data`返回响应数据


### 网关服务（所有服务接口集成于网关中）

#### 核心接口（消息收、发）

- 消息推送服务
    - **method**：ws、poll(http)、sse
    - **url**: `{gateway server}`/notify
    - **params**: 
        - token: `{用户端token}`
        - protoc: `{指定协议：ws|poll|sse}`
    - **MsgData**:
    ```json
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
    - **method**：POST
    - **url**: `{gateway server}`/send_msg
    - **header**:
        - Authorization: Bearer `{用户端token}`
    - **Request**:
    ```json
    {
        "conversation_id": "87ba7679-b682-47e7-8499-0385dda22b66",
        "status": 0,
        "type": 1,
        "content": "消息内容"
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data": 
        {
            "MsgID": "EDCZM-UUHJ-243Y-VWQ"
        }
    }
    ```


#### 会话相关

##### 公开接口：可供客户端直接

- 检查当前用户是否在会话中
    - **method**：POST
    - **url**: `{gateway server}`/check_joined_conversation
    - **header**:
        - Authorization: Bearer `{用户端token}``
    - **Request**:
    ```json
    {
        "conversation_id": "会话ID"
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data":
        {
            "Joined": false
        }
    }
    ```

- 获取当前用户加入的会话列表
    - **method**：POST
    - **url**: `{gateway server}`/get_joined_conversation
    - **header**:
        - Authorization: Bearer `{用户端token}`
    - **Request**:
    ```json
    {
        "SessionId": "用户ID"
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data": 
        {
            "Conversations": []
        }
    } 
    ```

- 获取会话信息
    - **method**：POST
    - **url**: `{gateway server}`/get_conversation_detail
    - **header**:
        - **Authorization**: Bearer `{用户端token}`
    - **Request**:
    ```json
    {
        "conversation_id": "会话ID"
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data": 
        {
            "ConversationID":   "会话ID",
            "LastMsgID":        "EDCZR-GL2A-243Y-VWQ",
            "LastSendTime":     "2023-10-22T22:59:04.069413421+08:00",
            "IsDelete":         false,
            "ConversationType": 0,
            "LastSendSession":  "87ba7679-b682-47e7-8499-0385dda22b66",
            "IsArchive":        false
        }
    }
    ```

- 加入会话
    - **method**：POST
    - **url**: `{gateway server}`/join_conversation
    - **header**:
        - **Authorization**: Bearer `{用户端token}`
    - **Request**:
    ```json
    {
        "conversation_id": "会话ID"
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data": 
        {
            "ConversationId": "会话ID"
        }
    }   
    ```

- 创建会话
    - **method**：POST
    - **url**: `{gateway server}`/create_conversation
    - **header**:
        - Authorization: Bearer `{用户端token}`
    - **Request**:
    ```json
    {
        "conversation_type": 0, // int64 会话类型
    }`
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data": 
        {
            "ConversationId": "会话ID"
        }
    }
    ```

- 离开会话
    - **method**：POST
    - **url**: `{gateway server}`/leave_conversation
    - **header**:
        - Authorization: Bearer `{用户端token}`
    - **Request**:
    ```json
    {
        "conversation_id": "会话ID"
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data": 
        {
            "ConversationId": "会话ID"
        }
    }
    ```

- 获取会话的最后一条消息ID
    - **method**：POST
    - **url**: `{gateway server}`/get_lastone_msgid
    - **header**:
        - Authorization: Bearer `{用户端token}`
    - **Request**:
    ```json
    {
        "conversation_id": "会话ID"
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data": 
        {
            "MsgId": "消息ID"
        }
    }  
    ```

##### 内部接口：不建议直接公开给客户端

- 获取Token
    - **method**：POST
    - **url**:` {gateway server}`/inner/get_token
    - **Request**:
    ```json
    {
        "session": "50864896-8136-4a43-8a48-1d3325a7f78f",
        "platform": 1
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data": "xxxxx"
    }
    ```

- 创建会话并初始化成员
    - **method**：POST
    - **url**: `{gateway server}`/inner/create_conversation
    - **Request**:
    ```json
    {
        "conversation_type": 0,
        "sessions": []
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data": 
        {
            "ConversationId": "会话ID"
        }
    }
    ```

- 踢出会话成员
    - **method**：POST
    - **url**: `{gateway server}`/inner/kickout_conversation
    - **Request**:
    ```json
    {
        "conversation_id": "87ba7679-b682-47e7-8499-0385dda22b66",
        "sessions": []
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data": 
        {
            "ConversationId": "会话ID"
        }
    }
    ```

- 指定用户加入会话
    - **method**：POST
    - **url**: `{gateway server}`/inner/join_conversation
    - **Request**:
    ```json
    {
        "conversation_id": "87ba7679-b682-47e7-8499-0385dda22b66",
        "sessions": []
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data": 
        {
            "ConversationId": "会话ID"
        }
    }
    ```

- 会话状态设为归档/非归档
    - **method**：POST
    - **url**: `{gateway server}`/inner/set_archive_conversations
    - **Request**:
    ```json
    {
        "conversation_ids" [],
        "is_archive": true
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
    }
    ```

#### 消息服务

##### 公开接口：可供客户端直接

- 获取会话指定消息ID范围的信息
    - **method**：POST
    - **url**: `{gateway server}`/get_conversation_msg_by_range
    - **Request**:
    ```json
    {
        "conversation_id": "",
        "start_msg_id": "",
        "end_msg_id": "",
        "desc": true,
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data": {
            "Msg": [
                {
                    "MsgId": "",
                    "ConversationID": "",
                    "FromSession": "",
                    "SendTime": "2023-10-22T22:59:04.069413421+08:00",
                    "Status": 0,
                    "Type":0"",
                    "Content": "",
                },
            ]
        }
    }
    ```

- 获取指定消息ID之前的30条消息
    - method：POST
    - url: `{gateway server}`/get_conversation_30msg_before_the_id
    - Request:
    ```json
    {
        "conversation_id": "",
        "msg_id": "",
        "desc": false,
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误", // code不为0时返回错误信息
        "data": {
            "Msg": [
                {
                    "MsgId": "",
                    "ConversationID": "",
                    "FromSession": "",
                    "SendTime": "2023-10-22T22:59:04.069413421+08:00",
                    "Status": 0,
                    "Type":0"",
                    "Content": "",
                },
            ]
        }
    }
    ```

- 获取指定消息ID之后的30条消息
    - **method**：POST
    - **url**: `{gateway server}`/- get_conversation_30msg_after_the_id
    - **Request**:
    ```json
    {
        "conversation_id": "",
        "msg_id": "",
        "desc": false,
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data": {
            "Msg": [
                {
                    "MsgId": "",
                    "ConversationID": "",
                    "FromSession": "",
                    "SendTime": "2023-10-22T22:59:04.069413421+08:00",
                    "Status": 0,
                    "Type":0"",
                    "Content": "",
                },
            ]
        }
    }
    ```

- 获取会话的最后30条消息
    - **method**：POST
    - **url**: `{gateway server}`/- get_conversation_last_30_msg
    - **Request**:
    ```json
    {
        "conversation_id": "",
        "desc": false,
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误", // code不为0时返回错误信息
        "data": {
            "Msg": [
                {
                    "MsgId": "",
                    "ConversationID": "",
                    "FromSession": "",
                    "SendTime": "2023-10-22T22:59:04.069413421+08:00",
                    "Status": 0,
                    "Type":0"",
                    "Content": "",
                },
            ]
        }
    }
    ```

- 获取会话之后一条消息
    - **method**：POST
    - **url**: `{gateway server}`/get_conversation_lastone_msg
    - **Request**:
    ```json
    {
        "conversation_id": "",
        "desc": false,
    }
    ```
    - **Response**:
    ```json
    {
        "code": 0,
        "error": "xxx错误",
        "data": {
            "Msg": [
                {
                    "MsgId": "",
                    "ConversationID": "",
                    "FromSession": "",
                    "SendTime": "2023-10-22T22:59:04.069413421+08:00",
                    "Status": 0,
                    "Type":0"",
                    "Content": "",
                },
            ]
        }
    }
    ```

##### 内部接口：不建议直接公开给客户端