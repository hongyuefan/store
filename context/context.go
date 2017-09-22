package context

type Context struct {
	IntervalTime    uint32
	OnceUpdateCount int
	Port            string
	DBUrl           string
	Threads         int
}

func NewContext(interval uint32, threadsNum int, onceUpdateCount int, port string, dburl string) *Context {

	return &Context{
		IntervalTime:    interval,
		OnceUpdateCount: onceUpdateCount,
		Port:            port,
		DBUrl:           dburl,
		Threads:         threadsNum,
	}
}
