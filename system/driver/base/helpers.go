package base

type FuncRun struct {
	F    func()
	Done chan struct{}
}
