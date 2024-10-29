package services

type Base struct {
	name string
	cfg  Config
}

func (s *Base) Name() string {
	return s.name
}

func (s *Base) Config() Config {
	return s.cfg
}

func NewService(name string, opts ...InitializeOption) *Base {
	s := &Base{
		name: name,
		cfg:  Config{},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type Config struct {
	HTTPServerPort   string
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePass     string
	DatabaseName     string
	DatabaseEngine   DatabaseEngineType
	Environment      string
	OtelCollectorURL string
}

type InitializeOption func(*Base)

func Database(opts ...DatabaseOption) InitializeOption {
	return func(s *Base) {
		for _, opt := range opts {
			opt(s)
		}
	}
}

type DatabaseOption func(*Base)

func DatabaseEngine(engine DatabaseEngineType) DatabaseOption {
	return func(s *Base) {
		s.cfg.DatabaseEngine = engine
	}
}

func DatabaseHost(host string) DatabaseOption {
	return func(s *Base) {
		s.cfg.DatabaseHost = host
	}
}

func DatabasePort(port string) DatabaseOption {
	return func(s *Base) {
		s.cfg.DatabasePort = port
	}
}

func DatabaseUser(user string) DatabaseOption {
	return func(s *Base) {
		s.cfg.DatabaseUser = user
	}
}

func DatabasePass(pass string) DatabaseOption {
	return func(s *Base) {
		s.cfg.DatabasePass = pass
	}
}

func DatabaseName(name string) DatabaseOption {
	return func(s *Base) {
		s.cfg.DatabaseName = name
	}
}

func Environment(env string) InitializeOption {
	return func(s *Base) {
		s.cfg.Environment = env
	}
}

func Traces(url string) InitializeOption {
	return func(s *Base) {
		s.cfg.OtelCollectorURL = url
	}
}

type HTTPApiOption func(*Base)

func HTTPServer(opts ...HTTPApiOption) InitializeOption {
	return func(s *Base) {
		for _, opt := range opts {
			opt(s)
		}
	}
}

func Port(port string) HTTPApiOption {
	return func(s *Base) {
		s.cfg.HTTPServerPort = port
	}
}
