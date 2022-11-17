package _const

const (
	/**
	 * trace链路核心信息
	 */
	TRACE_HEAD_ID             = "t-head-traceid"
	TRACE_HEAD_RPC_ID         = "t-head-rpcid"
	TRACE_HEAD_SAMPLED        = "t-head-sampled"
	TRACE_HEAD_USER_ID        = "t-head-userId"
	TRACE_HEAD_USER_NAME      = "t-head-userName"
	TRACE_HEAD_REMOTE_IP      = "t-head-remoteIp"
	TRACE_HEAD_REMOTE_APPNAME = "t-head-remoteAppName"
	TRACE_HEAD_ORIGNAL_URL    = "t-head-orignal-url"

	/**
	 * 附加到ATTR的信息
	 */
	A_REAL_PORT = "a-real-port"
	A_REAL_IP   = "a-real-ip"
	A_PERF      = "a-rerf"
	A_USER_ID   = "a-user-id"
	A_USER_NAME = "a-user-name"
	A_ERROR_MSG = "a-error-msg"
	A_WARN_MSG  = "a-warn-msg"

	/**
	 * 其他非链路关键字
	 */
	ISC_EXCEPT = "isc.except"
)
