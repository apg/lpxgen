package lpxgen

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

type ProbLog struct {
	logs       map[Log]float32
	sum        float32
	mutex      sync.Mutex
	currentLog Log
}

func logFromString(l string) Log {
	switch l {
	case "router":
		return Router
	case "dynomem":
		return DynoMem
	case "dynoload":
		return DynoLoad
	case "default":
		return DefaultLog{}
	default:
		fmt.Fprintf(os.Stderr, "WARNING: Invalid logtype %q: returning DefaultLog\n", l)
		return DefaultLog{}
	}
}

func ProbLogFromString(ds string) Log {
	logs := make(map[Log]float32)

	bits := strings.Split(ds, ",")
	for _, bit := range bits {
		logspec := strings.Split(bit, ":")
		switch len(logspec) {
		case 2:
			if val, err := strconv.ParseFloat(logspec[1], 32); err != nil {
				fmt.Printf("Invalid log spec: %q is not a valid number\n", logspec[1])
				os.Exit(1)
			} else {
				logs[logFromString(logspec[0])] = float32(val)
			}
		case 1:
			logs[logFromString(logspec[0])] = float32(1.0 / len(bits))
		default:
			fmt.Printf("Invalid log spec: %q\n", bit)
			os.Exit(1)
		}
	}

	return NewProbLog(logs)
}

func NewProbLog(logs map[Log]float32) *ProbLog {
	var sum float32
	for _, v := range logs {
		sum += v
	}

	return &ProbLog{
		logs:       logs,
		currentLog: DefaultLog{},
	}
}

func (p *ProbLog) Add(l Log, prob float32) {
	defer p.mutex.Unlock()
	p.mutex.Lock()

	p.logs[l] = prob
	p.sum += prob
}

func (p *ProbLog) PrivalVersion() string {
	return p.currentLog.PrivalVersion()
}

func (p *ProbLog) Time() string {
	return p.currentLog.Time()
}

func (p *ProbLog) Hostname() string {
	return p.currentLog.Hostname()
}

func (p *ProbLog) Name() string {
	return p.currentLog.Name()
}

func (p *ProbLog) Procid() string {
	return p.currentLog.Procid()
}

func (p *ProbLog) Msgid() string {
	return p.currentLog.Msgid()
}

func (p *ProbLog) Msg() string {
	return p.currentLog.Msg()
}

func (p *ProbLog) String() string {
	defer p.mutex.Unlock()

	p.mutex.Lock()
	p.currentLog = p.nextLog()
	s := FormatSyslog(p.currentLog)
	return fmt.Sprintf("%d %s", len(s), s)
}

func (p *ProbLog) nextLog() Log {
	var l Log
	var v float32

	if len(p.logs) == 0 {
		return p.currentLog
	}

	value := rand.Float32() * p.sum

	for l, v = range p.logs {
		value -= v
		if v < 0 {
			return l
		}
	}

	return l
}
