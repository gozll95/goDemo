go func(ctx context.Context, e *Entry) {
	defer func() {
		r := recover()
		if r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", err)
			}
			err = fmt.Errorf("panic: %v, stacktrace: %s", err, string(debug.Stack()))
			go c.errorsHandler(ctx, e.Job, err)
		}
	}()

	if c.funcCtx != nil {
		ctx = c.funcCtx(ctx, e.Job)
	}

	m, err := c.etcdclient.NewMutex(fmt.Sprintf("etcd_cron/%s/%d", e.Job.canonicalName(), effective.Unix()))
	if err != nil {
		go c.etcdErrorsHandler(ctx, e.Job, errors.Wrapf(err, "fail to create etcd mutex for job '%v'", e.Job.Name))
		return
	}
	lockCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err = m.Lock(lockCtx)
	if err == context.DeadlineExceeded {
		return
	} else if err != nil {
		go c.etcdErrorsHandler(ctx, e.Job, errors.Wrapf(err, "fail to lock mutex '%v'", m.Key()))
		return
	}

	err = e.Job.Run(ctx)
	if err != nil {
		go c.errorsHandler(ctx, e.Job, err)
		return
	}




很重要,好好利用context
	// err = m.Lock(lockCtx)
	// if err == context.DeadlineExceeded {
	// 	return
	// } else if err != nil {
	// 	go c.etcdErrorsHandler(ctx, e.Job, errors.Wrapf(err, "fail to lock mutex '%v'", m.Key()))
	// 	return
	// }