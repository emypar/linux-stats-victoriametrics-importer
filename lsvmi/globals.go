// Globals for LSVMI

package lsvmi

var (
	GlobalLsvmiConfig                    *LsvmiConfig
	GlobalHttpEndpointPool               *HttpEndpointPool
	GlobalCompressorPool                 *CompressorPool
	GlobalMetricsQueue                   MetricsQueue
	GlobalScheduler                      *Scheduler
	GlobalInstance                       string
	GlobalHostname                       string
	GlobalProcfsRoot                     string
	GlobalMetricsGeneratorStatsContainer *MetricsGeneratorStatsContainer
)
