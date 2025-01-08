package refactor_pattern1

import "time"

func New() *BatchProcessingWorkflow {
	ds := NewDoSomethingImpl(time.Second)
	bp := NewBatchProcessor(ds, 0, 20)
	bpw := NewBatchProcessingWorkflow(bp)
	return bpw
}
