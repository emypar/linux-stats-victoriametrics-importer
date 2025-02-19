// Configuration for the stats collector

package lsvmi

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

// The configuration is stored in one object to make it easy to load it from a
// file. Most of the configuration parameters are based on the file settings and
// a few can be overridden by command line arguments.
//
// The decreasing order of precedence for parameter values:
//   - command line arg (if applicable)
//   - config file
//   - built-in default
//
// Notes:
//  1. Each component will have its specific configuration, which may be
//     defined elsewhere, for instance in the files(s) providing the implementation.
//  2. The object will be loaded from a YAML file, therefore all configuration
//     parameters should be public and they should have tag annotations.

const (
	CONFIG_FILE_DEFAULT = "lsvmi-config.yaml"

	GLOBAL_CONFIG_INSTANCE_DEFAULT           = "lsvmi"
	GLOBAL_CONFIG_USE_SHORT_HOSTNAME_DEFAULT = true
	GLOBAL_CONFIG_PROCFS_ROOT_DEFAULT        = "/proc"
)

type LsvmiConfig struct {
	GlobalConfig                *GlobalConfig                `yaml:"global_config"`
	ProcStatMetricsConfig       *ProcStatMetricsConfig       `yaml:"proc_stat_metrics_config"`
	ProcNetDevMetricsConfig     *ProcNetDevMetricsConfig     `yaml:"proc_net_dev_metrics_config"`
	ProcInterruptsMetricsConfig *ProcInterruptsMetricsConfig `yaml:"proc_interrupts_metrics_config"`
	ProcSoftirqsMetricsConfig   *ProcSoftirqsMetricsConfig   `yaml:"proc_softirqs_metrics_config"`
	ProcNetSnmpMetricsConfig    *ProcNetSnmpMetricsConfig    `yaml:"proc_net_snmp_metrics_config"`
	ProcNetSnmp6MetricsConfig   *ProcNetSnmp6MetricsConfig   `yaml:"proc_net_snmp6_metrics_config"`
	ProcDiskstatsMetricsConfig  *ProcDiskstatsMetricsConfig  `yaml:"proc_diskstats_metrics_config"`
	ProcPidMetricsConfig        *ProcPidMetricsConfig        `yaml:"proc_pid_metrics_config"`
	StatfsMetricsConfig         *StatfsMetricsConfig         `yaml:"statfs_metrics_config"`
	QdiscMetricsConfig          *QdiscMetricsConfig          `yaml:"qdisc_metrics_config"`
	InternalMetricsConfig       *InternalMetricsConfig       `yaml:"internal_metrics_config"`
	SchedulerConfig             *SchedulerConfig             `yaml:"scheduler_config"`
	CompressorPoolConfig        *CompressorPoolConfig        `yaml:"compressor_pool_config"`
	HttpEndpointPoolConfig      *HttpEndpointPoolConfig      `yaml:"http_endpoint_pool_config"`
	LoggerConfig                *LoggerConfig                `yaml:"log_config"`
}

type GlobalConfig struct {
	// All metrics have the instance and hostname labels.

	// The instance name, default "lsvmi". It may be overridden by --instance
	// command line arg.
	Instance string `yaml:"instance"`

	// Whether to use short hostname or not as the value for hostname label.
	// Typically the hostname is determined from the hostname system call and if
	// the flag below is in effect, it is stripped of domain part. However if
	// the hostname is overridden by --hostname command line arg, that value is
	// used as-is.
	UseShortHostname bool `yaml:"use_short_hostname"`

	// procfs root. It may be overridden by --procfs-root command line arg.
	ProcfsRoot string `yaml:"procfs_root"`
}

var ErrConfigFileArgNotProvided = errors.New("config file arg not provided")

var lsvmiConfigFile = flag.String(
	"config",
	CONFIG_FILE_DEFAULT,
	`Config file to load`,
)

var hostnameArg = flag.String(
	"hostname",
	"",
	FormatFlagUsage(
		`Override the the value returned by hostname syscall`,
	),
)

var instanceArg = flag.String(
	"instance",
	"",
	FormatFlagUsage(
		`Override the "global_config.instance" config setting`,
	),
)

var procfsRootArg = flag.String(
	"procfs-root",
	"",
	FormatFlagUsage(
		`Override the "global_config.procfs_root" config setting`,
	),
)

var loggerLevelArg = flag.String(
	"log-level",
	"",
	FormatFlagUsage(fmt.Sprintf(
		`Override the "log_config.log_level" config setting, it should be one of the %q values`,
		GetLogLevelNames(),
	)),
)

var loggerFilelArg = flag.String(
	"log-file",
	"",
	FormatFlagUsage(
		`Override the config "log_config.log_file" config setting`,
	),
)

var httpPoolEndpointsArg = flag.String(
	"http-pool-endpoints",
	"",
	FormatFlagUsage(
		`Override the "http_endpoint_pool_config.endpoints" config setting`,
	),
)

func DefaultGlobalConfig() *GlobalConfig {
	return &GlobalConfig{
		Instance:         GLOBAL_CONFIG_INSTANCE_DEFAULT,
		UseShortHostname: GLOBAL_CONFIG_USE_SHORT_HOSTNAME_DEFAULT,
		ProcfsRoot:       GLOBAL_CONFIG_PROCFS_ROOT_DEFAULT,
	}
}

func DefaultLsvmiConfig() *LsvmiConfig {
	return &LsvmiConfig{
		GlobalConfig:                DefaultGlobalConfig(),
		ProcStatMetricsConfig:       DefaultProcStatMetricsConfig(),
		ProcNetDevMetricsConfig:     DefaultProcNetDevMetricsConfig(),
		ProcInterruptsMetricsConfig: DefaultProcInterruptsMetricsConfig(),
		ProcSoftirqsMetricsConfig:   DefaultProcSoftirqsMetricsConfig(),
		ProcNetSnmpMetricsConfig:    DefaultProcNetSnmpMetricsConfig(),
		ProcNetSnmp6MetricsConfig:   DefaultProcNetSnmp6MetricsConfig(),
		ProcDiskstatsMetricsConfig:  DefaultProcDiskstatsMetricsConfig(),
		ProcPidMetricsConfig:        DefaultProcPidMetricsConfig(),
		StatfsMetricsConfig:         DefaultStatfsMetricsConfig(),
		QdiscMetricsConfig:          DefaultQdiscMetricsConfig(),
		InternalMetricsConfig:       DefaultInternalMetricsConfig(),
		SchedulerConfig:             DefaultSchedulerConfig(),
		CompressorPoolConfig:        DefaultCompressorPoolConfig(),
		HttpEndpointPoolConfig:      DefaultHttpEndpointPoolConfig(),
	}
}

func LoadLsvmiConfig(cfgFile string) (*LsvmiConfig, error) {
	f, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	cfg := DefaultLsvmiConfig()
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("file: %q: %v", cfgFile, err)
	}
	return cfg, nil
}

func LoadLsvmiConfigFromArgs() (*LsvmiConfig, error) {
	if *lsvmiConfigFile == "" {
		return nil, ErrConfigFileArgNotProvided
	}
	cfg, err := LoadLsvmiConfig(*lsvmiConfigFile)
	if err != nil {
		return nil, err
	}

	// Apply command line overrides:
	if *instanceArg != "" {
		cfg.GlobalConfig.Instance = *instanceArg
	}
	if *procfsRootArg != "" {
		cfg.GlobalConfig.ProcfsRoot = *procfsRootArg
	}
	if *loggerLevelArg != "" {
		cfg.LoggerConfig.Level = *loggerLevelArg
	}
	if *loggerFilelArg != "" {
		cfg.LoggerConfig.LogFile = *loggerFilelArg
	}
	if *httpPoolEndpointsArg != "" {
		cfg.HttpEndpointPoolConfig.OverrideEndpoints(*httpPoolEndpointsArg)
	}

	return cfg, nil
}
