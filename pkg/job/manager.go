package job

import (
	"context"
	"sync"

	"github.com/fzerorubigd/balloon/pkg/assert"
	"github.com/fzerorubigd/balloon/pkg/config"
	"github.com/fzerorubigd/balloon/pkg/initializer"
	"github.com/fzerorubigd/balloon/pkg/redis"
	"github.com/fzerorubigd/chapar/middlewares/storage"
	"github.com/fzerorubigd/chapar/workers"
)

var (
	prefix  = config.RegisterString("jobs.manager.prefix", "w_", "job prefix in redis driver")
	manager *workers.Manager
	wrkrs   []workerSetup
	lock    sync.Mutex
	once    sync.Once
)

type workerSetup struct {
	name string
	w    workers.Worker
	opts []workers.WorkerOptions
}

type initSetup struct {
}

func (is *initSetup) Initialize(ctx context.Context) {
	driver := redis.NewDriver(ctx, prefix.String())
	manager = workers.NewManager(driver, driver)
	manager.RegisterMiddleware(
		storage.NewStorageMiddleware(redis.NewJobStore("job_" + prefix.String())),
	)
}

// Process all worker queues
func Process(ctx context.Context, opts ...workers.ProcessOptions) {
	assert.NotNil(manager, "make sure using this after finishing initialization")

	// This lock guarantee all call to this function will be blocked
	lock.Lock()
	defer lock.Unlock()

	// and only one of them (the first one) actually do thing
	once.Do(func() {
		for i := range wrkrs {
			assert.Nil(manager.RegisterWorker(wrkrs[i].name, wrkrs[i].w, wrkrs[i].opts...))
		}
		manager.Process(ctx, opts...)
	})
}

// RegisterWorker try to register new workers in system
func RegisterWorker(name string, worker workers.Worker, opt ...workers.WorkerOptions) {
	wrkrs = append(wrkrs, workerSetup{
		name: name,
		w:    worker,
		opts: opt,
	})
}

// EnqueueJob try to enqueue job in the queue
func EnqueueJob(ctx context.Context, queue string, data []byte, opts ...workers.EnqueueHandler) error {
	assert.NotNil(manager, "manager is empty")
	return manager.Enqueue(ctx, queue, data, opts ...)
}

func init() {
	initializer.Register(&initSetup{}, 10000)
}
