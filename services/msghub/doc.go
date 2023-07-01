/*
	消息中心服务模块
	Dep：
		Db
		Nats
	Feature：
		- 消息中心
		- 消息分发
			- nats
				-> 持久化模块
				-> 长连接推送模块
			- rpc (rpc调用为nats消息无应答的降级方案)
				-> 持久化模块
				-> 长连接推送模块
*/

package msghub
