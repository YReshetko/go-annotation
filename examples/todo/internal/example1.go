package internal

type ExploringTodoComment struct {
	Field struct {
		AnotherField struct {
			// @TODO(msg="Refactor the structure")
			Internal string
		}
	}
}
