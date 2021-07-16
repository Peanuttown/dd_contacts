package contacts_callback

type DingDingRequest struct {
	Encrypt string `json:"encrypt"`
}

type DingDingRequestDecryptedData struct {
}

type EventDataBase struct {
}

type EventUser struct {
	UserIds []string `json:"UserId"`
}

type EventDept struct {
	DeptIds []uint64 `json:"DeptId"`
}

// 用户添加
type EventUserAddData struct {
	EventUser
}

// 用户删除
type EventUserLeaveData struct {
	EventUser
}

type EventUserUpdate struct {
	EventUser
}

// 部门删除
type EventDeptDelete struct {
	EventDept
}

// 部门添加
type EventDeptAdd struct {
	EventDept
}

// 部门更新
type EventDeptUpdate struct {
	EventDept
}
