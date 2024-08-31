package ipcounter

type IPSet struct {
	set map[uint32]struct{}
}

func (m *IPSet) Add(ip uint32) {
	m.set[ip] = struct{}{}
}

func (m *IPSet) Count() uint64 {
	return uint64(len(m.set))
}
