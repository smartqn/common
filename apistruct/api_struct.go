package apistruct

type Ctx interface {
	Persist() error
	Load() error
}

type Msg struct {
	Question           string `json:"question"`
	Answer             string `json:"answer"`
	TargetUserType     string `json:"target_user_type,omitempty"`
	TargetUserTypeKind string `json:"target_user_type_kind,omitempty"`

	ResponseNodeID string `json:"response_node_id"`
	CondID         string `json:"cond_id"`
	ActionCode     string `json:"action_code"`
	ReqID          int64  `json:"req_id"`
	ReqTime        int64  `json:"req_time"`
	Index          int    `json:"index"`
	BreakTts       string `json:"_breakTts_"`
	PlayLast       string `json:"_playLast_"`
	ProblemStatus  int    `json:"_problemStatus_"`
}
