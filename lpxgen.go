package lpxgen

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
)

const (
	TimeFormat = "2006-01-02T15:04:05.000000+00:00"
)

var (
	UniqTokens int
	TokenPool  map[string][]string
)

func init() {
	TokenPool = make(map[string][]string)
}

type Log interface {
	PrivalVersion() string
	Time() string
	Hostname() string
	Name() string
	Procid() string
	Msgid() string
	Msg() string
	String() string
}

type LPXGenerator struct {
	Mincount int
	Maxcount int
	Log      Log
}

func NewGenerator(minBatch, maxBatch int, loggen Log) *LPXGenerator {
	return &LPXGenerator{
		Mincount: minBatch,
		Maxcount: maxBatch,
		Log:      loggen,
	}
}

func (g *LPXGenerator) Generate(url string) *http.Request {
	var body bytes.Buffer
	batchSize := g.Mincount + rand.Intn(g.Maxcount-g.Mincount)

	for i := 0; i < batchSize; i++ {
		body.WriteString(g.Log.String())
	}

	request, _ := http.NewRequest("POST", url, bytes.NewReader(body.Bytes()))
	request.Header.Add("Content-Length", string(body.Len()))
	request.Header.Add("Content-Type", "application/logplex-1")
	request.Header.Add("Logplex-Msg-Count", string(batchSize))
	request.Header.Add("Logplex-Frame-Id", randomFrameId())
	request.Header.Add("Logplex-Drain-Token", randomToken("d."))
	return request
}

func FormatSyslog(l Log) string {
	return fmt.Sprintf("%s %s %s %s %s %s %s\n",
		l.PrivalVersion(),
		l.Time(),
		l.Hostname(),
		l.Name(),
		l.Procid(),
		l.Msgid(),
		l.Msg())
}

var hexchars = []byte("0123456789abcdef")

func randomHexString(i int) string {
	x := make([]byte, i)
	for i > 0 {
		i--
		x[i] = hexchars[rand.Intn(16)]
	}
	return string(x)
}

func UUID4() string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		randomHexString(8),
		randomHexString(4),
		randomHexString(4),
		randomHexString(4),
		randomHexString(12))
}

func randomFrameId() string {
	return randomHexString(32)
}

func randomToken(prefix string) string {
	tokens, ok := TokenPool[prefix]
	if !ok {
		tokens = make([]string, UniqTokens)
		for i := 0; i < UniqTokens; i++ {
			tokens[i] = prefix + UUID4()
		}
		TokenPool[prefix] = tokens
	}

	return tokens[rand.Intn(len(tokens))]
}

func randomIPv4() string {
	return fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(255),
		rand.Intn(255),
		rand.Intn(255),
		rand.Intn(255))
}
