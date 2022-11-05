// Invalid example but that would be nice to explore the case deeply

package builder

func Execute() {
	f := NewFactory[CommonInterface]()
	f.Register("s1", NewStructureForExample1Builder)

	f.Register("s2", NewStructureForExample2Builder)

	f.New("s1").Do()
	f.New("s2").Do()
	f.New("s2").Do()
	f.New("s1").Do()
	f.New("s2").Do()
}
