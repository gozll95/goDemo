encoderCache.RLock()
	f := encoderCache.m[t.String()+"-"+tag]
	encoderCache.RUnlock()
	if f != nil {
		return f
	}

	// To deal with recursive types, populate the map with an
	// indirect func before we build it. This type waits on the
	// real func (f) to be ready and then calls it. This indirect
	// func is only used for recursive types.
	encoderCache.Lock()
	if encoderCache.m == nil {
		encoderCache.m = make(map[string]encoderFunc)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	encoderCache.m[t.String()+"-"+tag] = func(e *encodeState, v reflect.Value, opts encOpts) {
		wg.Wait()
		f(e, v, opts)
	}
	encoderCache.Unlock()

	// Compute fields without lock.
	// Might duplicate effort but won't hold other computations back.
	f = newTypeEncoder(t, tag, true)
	wg.Done()
	encoderCache.Lock()
	encoderCache.m[t.String()+"-"+tag] = f
	encoderCache.Unlock()
	fmt.Println("f is ", f)
	return f