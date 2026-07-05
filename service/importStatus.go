package service

import (
	"sync"
	"time"
)

const (
	ImportStatusIdle      = "idle"
	ImportStatusImporting = "importing"
	ImportStatusSuccess   = "success"
	ImportStatusError     = "error"
)

type ImportStatus struct {
	Status     string `json:"status"`
	Phase      string `json:"phase"`
	Message    string `json:"message"`
	Progress   int    `json:"progress"`
	FileName   string `json:"fileName"`
	FileSize   int64  `json:"fileSize"`
	StartedAt  string `json:"startedAt"`
	FinishedAt string `json:"finishedAt"`
	Error      string `json:"error"`
}

var importState = struct {
	sync.Mutex
	status ImportStatus
}{
	status: ImportStatus{
		Status:   ImportStatusIdle,
		Phase:    "idle",
		Message:  "等待导入",
		Progress: 0,
	},
}

func StartImportStatus(fileName string, fileSize int64) (ImportStatus, bool) {
	importState.Lock()
	defer importState.Unlock()

	if importState.status.Status == ImportStatusImporting {
		return importState.status, false
	}

	importState.status = ImportStatus{
		Status:    ImportStatusImporting,
		Phase:     "upload",
		Message:   "已接收文件，准备解析",
		Progress:  5,
		FileName:  fileName,
		FileSize:  fileSize,
		StartedAt: time.Now().Format(time.DateTime),
	}
	return importState.status, true
}

func UpdateImportStatus(phase, message string, progress int) {
	importState.Lock()
	defer importState.Unlock()

	if importState.status.Status != ImportStatusImporting {
		return
	}
	importState.status.Phase = phase
	importState.status.Message = message
	importState.status.Progress = progress
}

func FinishImportStatus(err error) {
	importState.Lock()
	defer importState.Unlock()

	importState.status.FinishedAt = time.Now().Format(time.DateTime)
	if err != nil {
		importState.status.Status = ImportStatusError
		importState.status.Phase = "error"
		importState.status.Message = "导入失败"
		importState.status.Progress = 100
		importState.status.Error = err.Error()
		return
	}

	importState.status.Status = ImportStatusSuccess
	importState.status.Phase = "done"
	importState.status.Message = "导入完成，统计数据已刷新"
	importState.status.Progress = 100
	importState.status.Error = ""
}

func GetImportStatus() ImportStatus {
	importState.Lock()
	defer importState.Unlock()

	return importState.status
}
