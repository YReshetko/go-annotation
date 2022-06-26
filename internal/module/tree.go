package module

import "fmt"

type pkg string

type selector struct {
	pkg  pkg
	next *selector
}

func (s *selector) put(v string) selector {
	newSelector := selector{}
	if s.next == nil {
		newSelector.pkg = pkg(v)
		return newSelector
	}
	newSelector.pkg = s.pkg
	next := s.next.put(v)
	newSelector.next = &next
	return newSelector
}

type tree struct {
	m *Module
	s map[pkg]*tree
}

func newTree() *tree {
	return &tree{
		s: make(map[pkg]*tree),
	}
}

func (t *tree) put(s selector, m *Module) error {
	if s.next == nil {
		if t.m != nil {
			return fmt.Errorf("unable to replace module in tree for pkg %s", s.pkg)
		}
		t.m = m
		return nil
	}

	if t.m == nil {
		return fmt.Errorf("unable to create modules tree node as no existing module for %s", s.pkg)
	}

	s = *s.next

	nt, ok := t.s[s.pkg]
	if !ok {
		nt = newTree()
		t.s[s.pkg] = nt
	}
	return nt.put(s, m)
}

func (t *tree) get(s selector) (*Module, bool) {
	if s.next == nil {
		return t.m, t.m != nil
	}
	s = *s.next
	nt, ok := t.s[s.pkg]
	if !ok {
		return nil, false
	}
	return nt.get(s)
}

func (s *selector) String() string {
	if s.next == nil {
		return string(s.pkg)
	}
	return string(s.pkg) + "\n"
}
