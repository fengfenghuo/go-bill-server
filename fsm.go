package approve

import (
	"bpm/models/db"

	"github.com/astaxie/beego/logs"
	"github.com/looplab/fsm"
)

// ApprovalFSM 状态机模型
type ApprovalFSM struct {
	Fsm *fsm.FSM
}

// 审批动作定义
const (
	ApprovalBegin     string = "toApproval"
	ApprovalCompleted string = "complete"
	ApprovalRejected  string = "reject"
	ApprovalClose     string = "close"
)

func NewApproveFsm() *ApprovalFSM {
	af := &ApprovalFSM{}

	af.Fsm = fsm.NewFSM(db.ApprovalStateNone,
		fsm.Events{
			{
				Name: ApprovalBegin,
				Src:  []string{db.ApprovalStateNone},
				Dst:  db.ApprovalStateApproval,
			},
			{
				Name: ApprovalCompleted,
				Src:  []string{db.ApprovalStateApproval},
				Dst:  db.ApprovalStateApproved,
			},
			{
				Name: ApprovalRejected,
				Src:  []string{db.ApprovalStateApproval},
				Dst:  db.ApprovalStateRejected,
			},
			{
				Name: ApprovalClose,
				Src:  []string{db.ApprovalStateApproval, db.ApprovalStateApproved, db.ApprovalStateRejected},
				Dst:  db.ApprovalStateNone,
			},
		},
		fsm.Callbacks{
			"enter_" + db.ApprovalStateNone:     func(e *fsm.Event) { af.enterClosed(e) },
			"enter_" + db.ApprovalStateApproval: func(e *fsm.Event) { af.enterApprovaling(e) },
			"enter_" + db.ApprovalStateApproved: func(e *fsm.Event) { af.enterApproved(e) },
			"enter_" + db.ApprovalStateRejected: func(e *fsm.Event) { af.enterRejected(e) },
		},
	)
	return af
}

func (af *ApprovalFSM) enterClosed(e *fsm.Event) {
	if len(e.Args) == 0 {
		logs.Error("approve:ApprovalFSM.enterClosed args not exist")
		return
	}
	bpmDetail, ok := e.Args[0].(*db.BPMObjectModel)
	if !ok {
		logs.Error("approve:ApprovalFSM.enterClosed args error")
		return
	}
	logs.Debug("approve:ApprovalFSM.enterClosed %s", bpmDetail.BPMID)
}

func (af *ApprovalFSM) enterApprovaling(e *fsm.Event) {
	if len(e.Args) == 0 {
		logs.Error("approve:ApprovalFSM.enterApprovaling args not exist")
		return
	}
	bpmDetail, ok := e.Args[0].(*db.BPMObjectModel)
	if !ok {
		logs.Error("approve:ApprovalFSM.enterApprovaling args error")
		return
	}
	logs.Debug("approve:ApprovalFSM.enterApprovaling %s", bpmDetail.BPMID)

	processor, ok := e.Args[1].(ApprovalProcessor)
	if !ok {
		logs.Error("approve:ApprovalFSM.enterApprovaling args error")
		return
	}

	if err := processor.Initiate(bpmDetail); err != nil {
		logs.Error("approve:ApprovalFSM.enterApprovaling.Initiate %s err: %v", bpmDetail.BPMID, err)
	}
}

func (af *ApprovalFSM) enterApproved(e *fsm.Event) {
	if len(e.Args) == 0 {
		logs.Error("approve:ApprovalFSM.enterApproved args not exist")
		return
	}
	bpmDetail, ok := e.Args[0].(*db.BPMObjectModel)
	if !ok {
		logs.Error("approve:ApprovalFSM.enterApproved args error")
		return
	}
	logs.Debug("approve:ApprovalFSM.enterApproved %s", bpmDetail.BPMID)
}

func (af *ApprovalFSM) enterRejected(e *fsm.Event) {
	if len(e.Args) == 0 {
		logs.Error("approve:ApprovalFSM.enterRejected args not exist")
		return
	}
	bpmDetail, ok := e.Args[0].(*db.BPMObjectModel)
	if !ok {
		logs.Error("approve:ApprovalFSM.enterRejected args error")
		return
	}
	logs.Debug("approve:ApprovalFSM.enterRejected %s", bpmDetail.BPMID)
}
