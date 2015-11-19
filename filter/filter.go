package filter

import (
	"net"
	"strings"
	"sync"

	manet "gx/QmNT7d1e4Xcp3KcsvxyzUHVtqrR43uypoxLLzdKj6YZga2/go-multiaddr-net"
	ma "gx/QmVUi2ncqnU48zsPgR1rQosDGwY3SSZ1Ndp33j33YjXdsj/go-multiaddr"
)

type Filters struct {
	mu      sync.RWMutex
	filters map[string]*net.IPNet
}

func NewFilters() *Filters {
	return &Filters{
		filters: make(map[string]*net.IPNet),
	}
}

func (fs *Filters) AddDialFilter(f *net.IPNet) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.filters[f.String()] = f
}

func (f *Filters) AddrBlocked(a ma.Multiaddr) bool {
	_, addr, err := manet.DialArgs(a)
	if err != nil {
		// if we cant parse it, its probably not blocked
		return false
	}

	ipstr := strings.Split(addr, ":")[0]
	ip := net.ParseIP(ipstr)
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, ft := range f.filters {
		if ft.Contains(ip) {
			return true
		}
	}
	return false
}

func (f *Filters) Filters() []*net.IPNet {
	var out []*net.IPNet
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, ff := range f.filters {
		out = append(out, ff)
	}
	return out
}

func (f *Filters) Remove(ff *net.IPNet) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.filters, ff.String())
}
