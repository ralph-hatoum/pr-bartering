package worker

import (
	"go.uber.org/zap"
)

type Workload struct {
	WorkloadType byte
	Cid          string
	PeerIP       string
}

func Worker(workloadChannel chan Workload, logger *zap.Logger) {
	logger.Info("Worker started")

	for workload := range workloadChannel {
		switch workload.WorkloadType {
		case 0x01:
			handle_storage_request(workload)
		case 0x02:
			handle_proof_request(workload)
		}
	}
}

func handle_storage_request(workload Workload) {
	return
}

func handle_proof_request(workload Workload) {
	return
}
