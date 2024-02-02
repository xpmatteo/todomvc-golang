package todo

type Item struct {
	Title  string
	IsDone bool
}

type List struct {
	Items []Item
}

func NewList() List {
	return List{}
}

func (l *List) Add(title string) {
	if len(title) == 0 {
		return
	}
	l.Items = append(l.Items, Item{title, false})
}