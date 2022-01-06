package conf

var (
	LenStackBuf = 4096

	// log
	LogLevel string
	LogPath  string
	LogFlag  int

	// console
	ConsolePort   int
	ConsolePrompt string = "Leaf# "
	ProfilePath   string

	// cluster
	ListenAddr      string
	ConnAddrs       []string  //A,B两个进程 分别监听不同的端口,相互连通  是两个单独的client,单独的连接. 感觉作者像这个意思.
	PendingWriteNum int
)
