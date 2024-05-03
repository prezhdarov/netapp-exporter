package ontapCollectors

/*
const (
	clusterSubsystem = "cluster"
)

var clusterCollectorFlag = flag.Bool(fmt.Sprintf("collector.%s", clusterSubsystem), collector.DefaultEnabled, fmt.Sprintf("Enable the %s collector (default: %v)", clusterSubsystem, collector.DefaultEnabled))

type clusterCollector struct {
	logger log.Logger
}

func init() {
	collector.RegisterCollector("cluster", clusterCollectorFlag, NewClusterCollector)
}

// NewMeminfoCollector returns a new Collector exposing memory stats.
func NewClusterCollector(logger log.Logger) (collector.Collector, error) {
	return &clusterCollector{logger}, nil
}

func (c *clusterCollector) Update(ch chan<- prometheus.Metric, namespace string, clientAPI collector.ClientAPI, loginData map[string]interface{}, params map[string]string) error {

	extraConfig := make(map[string]interface{}, 0)

	extraConfig["api"] = "/api/cluster/?fields=metric"
	body, err := clientAPI.Get(loginData, extraConfig, c.logger)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*body.(*[]byte), &data)
	if err != nil {
		return err
	}

	//clusterLabels := map[string]string{"nacluster": loginData["target"].(string)}

	return nil
}
*/
