package natssub

type Options struct {
	jetstream   bool
	queue       string
	subject     string
	concurrency int
	buffer      int
	handler     MessageHandler
}

func (o *Options) validate() error {
	return nil
}

type Option = func(o *Options)

func NewOptions() *Options {
	return &Options{
		concurrency: 1,
		buffer:      0,
	}
}

// TODO: Stream Configuration needed
func WithJetStream() Option {
	return func(o *Options) {
		o.jetstream = true
	}
}

func WithQueue(queue string) Option {
	return func(o *Options) {
		o.queue = queue
	}
}

func WithSubject(subject string) Option {
	return func(o *Options) {
		o.subject = subject
	}
}

func WithConcurrency(concurrency int) Option {
	return func(o *Options) {
		o.concurrency = concurrency
	}
}

func WithBuffer(buffer int) Option {
	return func(o *Options) {
		o.buffer = buffer
	}
}

func WithHandler(handler MessageHandler) Option {
	return func(o *Options) {
		o.handler = handler
	}
}
