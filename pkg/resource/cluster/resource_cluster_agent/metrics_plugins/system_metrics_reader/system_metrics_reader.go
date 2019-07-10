package system_metrics_reader

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/syunkitada/goapp/pkg/config"
	"github.com/syunkitada/goapp/pkg/lib/logger"
	"github.com/syunkitada/goapp/pkg/resource/resource_model"
)

type SystemMetricsReader struct {
	name         string
	enableLogin  bool
	enableCpu    bool
	enableMemory bool
	cacheLength  int
	uptimeStats  []UptimeStat
	loginStats   []LoginStat
	cpuStats     []CpuStat
	numaNodes    []resource_model.NumaNodeSpec
}

func New(conf *config.ResourceMetricsSystemConfig) *SystemMetricsReader {
	// f, _ := os.Open("/sys/devices/system/node/online")
	// /sys/devices/system/node/node0/meminfo
	nodeOnline, err := os.Open("/sys/devices/system/node/online")
	if err != nil {
		logger.StdoutFatalf("Failed Initialize SystemMetricsReader: %v", err)
	}
	defer nodeOnline.Close()

	nodeOnlineBytes, err := ioutil.ReadAll(nodeOnline)
	if err != nil {
		logger.StdoutFatalf("Failed Initialize SystemMetricsReader: %v", err)
	}
	splitedNodes := strings.Split(strings.TrimRight(string(nodeOnlineBytes), "\n"), ",")

	numaNodes := []resource_model.NumaNodeSpec{}
	for _, node := range splitedNodes {
		id, err := strconv.Atoi(node)
		if err != nil {
			logger.StdoutFatalf("Failed Initialize SystemMetricsReader: %v", err)
		}
		numaNodes = append(numaNodes, resource_model.NumaNodeSpec{
			Id: id,
		})
	}

	return &SystemMetricsReader{
		name:         "system",
		enableLogin:  conf.EnableLogin,
		enableCpu:    conf.EnableCpu,
		enableMemory: conf.EnableMemory,
		cacheLength:  conf.CacheLength,
		uptimeStats:  make([]UptimeStat, 0, conf.CacheLength),
		loginStats:   make([]LoginStat, 0, conf.CacheLength),
		cpuStats:     make([]CpuStat, 0, conf.CacheLength),
		numaNodes:    numaNodes,
	}
}

type UptimeStat struct {
	reportStatus int // 0, 1(GetReport), 2(Reported)
	timestamp    time.Time
	uptime       int64
}

type LoginStat struct {
	users     []UserStat
	timestamp time.Time
}

type UserStat struct {
	user  string
	tty   string
	from  string
	login string
	idle  string
	jcpu  string
	pcpu  string
	what  string
}

type CpuStat struct {
	reportStatus  int // 0, 1(GetReport), 2(Reported)
	timestamp     time.Time
	cpuMap        map[string][]int64
	intr          int64
	ctx           int64
	btime         int64
	processes     int64
	procs_running int64
	procs_blocked int64
	softirq       int64
}

func (reader *SystemMetricsReader) GetNumaNodes(tctx *logger.TraceContext) []resource_model.NumaNodeSpec {
	return reader.numaNodes
}

func (reader *SystemMetricsReader) Read(tctx *logger.TraceContext) error {
	timestamp := time.Now()

	// Read /proc/uptime
	// uptime(s)  idle(s)
	// 2906.26 5507.43
	fmt.Println("READ /proc/uptime")
	procUptime, _ := os.Open("/proc/uptime")
	defer procUptime.Close()
	scanner := bufio.NewScanner(procUptime)
	scanner.Scan()
	uptimeText := scanner.Text()
	uptimeWords := strings.Split(uptimeText, " ")
	uptime, _ := strconv.ParseInt(uptimeWords[0], 10, 64)
	uptimeStat := UptimeStat{
		reportStatus: 0,
		timestamp:    timestamp,
		uptime:       uptime,
	}
	if len(reader.uptimeStats) > reader.cacheLength {
		reader.uptimeStats = reader.uptimeStats[1:]
	}
	reader.uptimeStats = append(reader.uptimeStats, uptimeStat)

	if reader.enableLogin {
		// Don't read /var/run/utmp, because of this is binary
		// Read w -h
		// USER    TTY      FROM          LOGIN@   IDLE   JCPU   PCPU  WHAT
		// hoge    pts/8    192.168.1.1   09:34    2.00s  0.10s  0.00s tmux a
		out, err := exec.Command("w", "-h").Output()
		users := []UserStat{}
		if err != nil {
			fmt.Println(err)
		}
		for _, line := range strings.Split(string(out), "\n") {
			l := strings.Split(line, " ")
			if len(l) != 8 {
				continue
			}
			users = append(users, UserStat{
				user:  l[0],
				tty:   l[1],
				from:  l[2],
				login: l[3],
				idle:  l[4],
				jcpu:  l[5],
				pcpu:  l[6],
				what:  l[7],
			})
		}
		if len(reader.loginStats) > reader.cacheLength {
			reader.loginStats = reader.loginStats[1:]
		}
		reader.loginStats = append(reader.loginStats, LoginStat{
			timestamp: timestamp,
			users:     users,
		})
	}

	if reader.enableCpu {
		// Read /proc/stat
		//      user   nice system idle    iowait irq softirq steal guest guest_nice
		// cpu  264230 262  60792  8237284 20685  0   2652    0     0     0
		// cpu0 126387 2    30266  4124610 11105  0   1011    0     0     0
		// cpu1 137842 260  30525  4112674 9580   0   1641    0     0     0
		// intr 18316761 ...
		// ctxt 57087643
		// btime 1546819593
		// processes 227393
		// procs_running 1
		// procs_blocked 0
		// softirq 11650881 ...

		fmt.Println("READ /proc/stat")
		f, _ := os.Open("/proc/stat")
		defer f.Close()
		scanner = bufio.NewScanner(f)
		lines := make([]string, 0, 20)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		var cpu []string
		cpuMap := map[string][]int64{}
		lastIndex := 1
		for i := 0; i < lastIndex; i++ {
			cpu = strings.Split(lines[i], " ")
			v := make([]int64, len(cpu)-1)
			for j, c := range cpu[1:] {
				v[j], _ = strconv.ParseInt(c, 10, 64)
			}
			cpuMap[cpu[0]] = v
		}

		intr, _ := strconv.ParseInt(strings.Split(lines[lastIndex], " ")[1], 10, 64)
		ctx, _ := strconv.ParseInt(strings.Split(lines[lastIndex+1], " ")[1], 10, 64)
		btime, _ := strconv.ParseInt(strings.Split(lines[lastIndex+2], " ")[1], 10, 64)
		processes, _ := strconv.ParseInt(strings.Split(lines[lastIndex+3], " ")[1], 10, 64)
		procs_running, _ := strconv.ParseInt(strings.Split(lines[lastIndex+4], " ")[1], 10, 64)
		procs_blocked, _ := strconv.ParseInt(strings.Split(lines[lastIndex+5], " ")[1], 10, 64)
		softirq, _ := strconv.ParseInt(strings.Split(lines[lastIndex+6], " ")[1], 10, 64)
		stat := CpuStat{
			reportStatus:  0,
			timestamp:     timestamp,
			cpuMap:        cpuMap,
			intr:          intr,
			ctx:           ctx,
			btime:         btime,
			processes:     processes,
			procs_running: procs_running,
			procs_blocked: procs_blocked,
			softirq:       softirq,
		}

		if len(reader.cpuStats) > reader.cacheLength {
			reader.cpuStats = reader.cpuStats[1:]
		}
		reader.cpuStats = append(reader.cpuStats, stat)
	}

	if reader.enableMemory {
		// Read /proc/vmstat
		fmt.Println("READ Memory")
	}

	fmt.Println(reader.cpuStats)
	return nil
}

func (reader *SystemMetricsReader) GetName() string {
	return reader.name
}

func (reader *SystemMetricsReader) Report() ([]resource_model.Metric, []resource_model.Alert) {
	metrics := make([]resource_model.Metric, 0, 100)
	alerts := make([]resource_model.Alert, 0, 100)

	for _, stat := range reader.uptimeStats {
		timestamp := strconv.FormatInt(stat.timestamp.UnixNano(), 10)
		metrics = append(metrics, resource_model.Metric{
			Name: "system_uptime",
			Time: timestamp,
			Tag:  map[string]string{},
			Metric: map[string]int64{
				"uptime": stat.uptime,
			},
		})
	}

	for _, stat := range reader.loginStats {
		timestamp := strconv.FormatInt(stat.timestamp.UnixNano(), 10)
		metrics = append(metrics, resource_model.Metric{
			Name: "system_login",
			Time: timestamp,
			Tag:  map[string]string{},
			Metric: map[string]int64{
				"users": int64(len(stat.users)),
			},
		})
	}

	for _, stat := range reader.cpuStats {
		timestamp := strconv.FormatInt(stat.timestamp.UnixNano(), 10)
		for cpuName, cpu := range stat.cpuMap {
			metrics = append(metrics, resource_model.Metric{
				Name: "system_cpu",
				Time: timestamp,
				Tag: map[string]string{
					"cpu": cpuName,
				},
				Metric: map[string]int64{
					"user":       cpu[0],
					"nice":       cpu[1],
					"system":     cpu[2],
					"idle":       cpu[3],
					"iowait":     cpu[4],
					"irq":        cpu[5],
					"softirq":    cpu[6],
					"steal":      cpu[7],
					"guest":      cpu[8],
					"guest_nice": cpu[9],
				},
			})
		}

		metrics = append(metrics, resource_model.Metric{
			Name: "system_cpu",
			Time: timestamp,
			Tag: map[string]string{
				"cpu": "cpu",
			},
			Metric: map[string]int64{
				"intr":          stat.intr,
				"ctx":           stat.ctx,
				"btime":         stat.btime,
				"processes":     stat.processes,
				"procs_running": stat.procs_running,
				"procs_blocked": stat.procs_blocked,
				"softirq":       stat.softirq,
			},
		})

		stat.reportStatus = 1
	}

	// TODO check metrics and issue alerts

	return metrics, alerts
}

func (reader *SystemMetricsReader) Reported() {
	for _, stat := range reader.cpuStats {
		stat.reportStatus = 2
	}
}