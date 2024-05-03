package ontapCollectors

type Aggregates struct {
	Records []Aggregate `json:"records"`
}

type Aggregate struct {
	Name string `json:"name"`
	Node struct {
		Name string `json:"name"`
	} `json:"home_node"`
	Space struct {
		Block struct {
			Size                 float64 `json:"size"`
			Available            float64 `json:"available"`
			Used                 float64 `json:"used"`
			InactiveData         float64 `json:"inactive_user_data"`
			PhysicalUsed         float64 `json:"physical_used"`
			Metadata             float64 `json:"aggregate_metadata"`
			UsedWithSnapReserve  float64 `json:"used_including_snapshot_reserve"`
			DataCompactCount     float64 `json:"data_compacted_count"`
			DataCompactSaved     float64 `json:"data_compaction_space_saved"`
			VolDedupeSharedCount float64 `json:"volume_deduplication_shared_count"`
			VolDedupeSaved       float64 `json:"volume_deduplication_space_saved"`
			//VolDedupeSavedPct      float64 `json:"volume_deduplication_space_saved_percent"`
			//UsedPct                float64 `json:"used_percent"`
			//DataCompactSavedPct    float64 `json:"data_compaction_space_saved_percent"`
			//UsedWithSnapReservePct float64 `json:"used_including_snapshot_reserve_percent"`
			//MetadataPct            float64 `json:"aggregate_metadata_percent"`
			//PhysicalUsedPct        float64 `json:"physical_used_percent"`
			//InactiveDataPct        float64 `json:"inactive_user_data_percent"`
			//FullThresholdPct     float64 `json:"full_threshold_percent"`
		} `json:"block_storage"`
	} `json:"space"`
	Metric struct {
		Timestamp  string `json:"timestamp"`
		Duration   string `json:"duration"`
		Status     string `json:"status"`
		Throughput struct {
			Read  float64 `json:"read"`
			Write float64 `json:"write"`
			Other float64 `json:"other"`
			//Total float64 `json:"total"`
		} `json:"throughput"`
		Latency struct {
			Read  float64 `json:"read"`
			Write float64 `json:"write"`
			Other float64 `json:"other"`
			//Total float64 `json:"total"`
		} `json:"latency"`
		IOps struct {
			Read  float64 `json:"read"`
			Write float64 `json:"write"`
			Other float64 `json:"other"`
			//Total float64 `json:"total"`
		} `json:"iops"`
	} `json:"metric"`
}
