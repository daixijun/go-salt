package salt

type WorkerThread struct {
	Requests        int     `json:"Requests"`
	BytesRead       int     `json:"Bytes Read"`
	BytesWritten    int     `json:"Bytes Written"`
	WorkTime        float64 `json:"Work Time"`
	ReadThroughput  float64 `json:"Read Throughput"`
	WriteThroughput float64 `json:"Write Throughput"`
}

type HTTPServer struct {
	Enabled         bool                    `json:"Enabled"`
	BindAddress     string                  `json:"Bind Address"`
	RunTime         int                     `json:"Run time"`
	Accepts         int                     `json:"Accepts"`
	AcceptsSec      float64                 `json:"Accepts/sec"`
	Queue           int                     `json:"Queue"`
	Threads         int                     `json:"Threads"`
	ThreadsIdle     int                     `json:"Threads Idle"`
	SocketErrors    int                     `json:"Socket Errors"`
	Requests        int                     `json:"Requests"`
	BytesRead       int                     `json:"Bytes Read"`
	BytesWritten    int                     `json:"Bytes Written"`
	WorkTime        int                     `json:"Work Time"`
	ReadThroughput  float64                 `json:"Read Throughput"`
	WriteThroughput float64                 `json:"Write Throughput"`
	WorkerThreads   map[string]WorkerThread `json:"Worker Threads"`
}
