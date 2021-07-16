package event_type

type EventType string

const (
	EVENT_TYPE_CHECK_URL   = "check_url"
	EVENT_TYPE_USER_ADD    = "user_add_org"
	EVENT_TYPE_USER_LEAVE  = "user_leave_org"
	EVENT_TYPE_USER_UPDATE = "user_modify_org"
	EVENT_TYPE_DEPT_DEL    = "org_dept_remove"
	EVENT_TYPE_DEPT_ADD    = "org_dept_create"
	EVENT_TYPE_DEPT_UPDATE = "org_dept_modify"
)
