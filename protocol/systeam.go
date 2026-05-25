package protocol

type SysFeature struct {
	OverallLatency     int64 `json:"overall_latency"`     // 总体延迟
	CollectionLatency  int64 `json:"collection_latency"`  // 采集延迟
	ComputationLatency int64 `json:"computation_latency"` // 计算延迟
}
