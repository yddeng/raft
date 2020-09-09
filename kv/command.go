package kv

type SetCmd struct {
	Key   string
	Value string
}

func (this *SetCmd) CommandName() string {
	return "set"
}

// handler struct

type SetReq struct {
	Key   string
	Value string
}

type SetResp struct {
	OK       bool
	IsLeader bool
}

type GetReq struct {
	Key string
}

type GetResp struct {
	OK    bool
	Value string
}
