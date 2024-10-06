package imports

import (
	"os"
	"strings"

	"path/filepath"

	"github.com/YReshetko/go-annotation/internal/utils/arrays"
)

type Path string

const EmptyPath Path = "."

func (p Path) Intersection(inPath Path) Path {
	result, _, _ := arrays.LCS(p.split(), inPath.split())
	return join(result)
}

func (p Path) Left(inPath Path) Path {
	path := p.split()
	_, i, _ := arrays.LCS(path, inPath.split())
	if i < 0 {
		return EmptyPath
	}
	return join(path[:i])
}

func (p Path) LeftJoin(inPath Path) Path {
	path := p.split()
	result, i, _ := arrays.LCS(path, inPath.split())
	if i < 0 {
		return EmptyPath
	}
	return join(path[:i+len(result)])
}

func (p Path) Right(inPath Path) Path {
	path := p.split()
	result, i, _ := arrays.LCS(path, inPath.split())
	if i < 0 {
		return EmptyPath
	}
	return join(path[i+len(result):])
}

func (p Path) RightJoin(inPath Path) Path {
	path := p.split()
	_, i, _ := arrays.LCS(path, inPath.split())
	if i < 0 {
		return EmptyPath
	}
	return join(path[i:])
}

func (p Path) FullJoin(inPath Path) Path {
	pathSplitted := p.split()
	inPathSplitted := inPath.split()
	result, i, j := arrays.LCS(pathSplitted, inPathSplitted)
	if len(result) == 0 {
		return EmptyPath
	}
	return join(append(pathSplitted[:i+len(result)], inPathSplitted[j+len(result):]...))
}

func (p Path) IsEmpty() bool {
	return p == EmptyPath || len(p.split()) == 0
}

func (p Path) String() string {
	return string(p)
}

func Of(s string) Path {
	return Path(s)
}

func (p Path) split() []string {
	items := strings.FieldsFunc(p.String(), func(r rune) bool {
		return r == '\\' || r == '/'
	})

	var retItems []string
	for _, item := range items {
		if len(item) == 0 || item == "." {
			continue
		}
		retItems = append(retItems, item)
	}
	return retItems
}

func join(strs []string) Path {
	if len(strs) == 0 {
		return EmptyPath
	}
	if strings.Contains(strs[0], ":") {
		return Path(strs[0] + string(os.PathSeparator) + filepath.Join(strs[1:]...))
	}
	return Path(filepath.Join(strs...))
}
