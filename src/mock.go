//types

// An encodeState encodes JSON into a bytes.Buffer.
type encodeState struct {
	bytes.Buffer // accumulated output
	scratch      [64]byte
}

type encOpts struct {
	// quoted causes primitive fields to be encoded inside JSON strings.
	quoted bool
	// escapeHTML causes '<', '>', and '&' to be escaped in JSON strings.
	escapeHTML bool
}

type encoderFunc func(e *encodeState, v reflect.Value, opts encOpts)


//这里直接是var
var encoderCache struct {
	sync.RWMutex
	m map[reflect.Type]encoderFunc
}


var (
	marshalerType     = reflect.TypeOf(new(Marshaler)).Elem()
	textMarshalerType = reflect.TypeOf(new(encoding.TextMarshaler)).Elem()
)

type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

var fieldCache struct {
	value atomic.Value // map[reflect.Type][]field
	mu    sync.Mutex   // used only by writers
}


// A field represents a single field found in a struct.
type field struct {
	name      string
	nameBytes []byte                 // []byte(name)
	equalFold func(s, t []byte) bool // bytes.EqualFold or equivalent

	tag       bool
	index     []int
	typ       reflect.Type
	omitEmpty bool
	quoted    bool
}

type structEncoder struct {
	fields    []field
	fieldEncs []encoderFunc
}



b, err := json.Marshal(s)
		- func Marshal(v interface{}) ([]byte, error) 
			- e := &encodeState{}
			- err := e.marshal(v, encOpts{escapeHTML: true})
					- func (e *encodeState) marshal(v interface{}, opts encOpts) (err error) 
						- //defer recover from panic
						- defer func(){...}()
						- e.reflectValue(reflect.ValueOf(v), opts)
								- func (e *encodeState) reflectValue(v reflect.Value, opts encOpts)
									- valueEncoder(v)(e, v, opts) //前面是一个encoderFunc
											-  func valueEncoder(v reflect.Value) encoderFunc
												- return typeEncoder(v.Type())
														- func typeEncoder(t reflect.Type) encoderFunc
															- encoderCache.RLock()
															- f := encoderCache.m[t] //哪里set的?
															- encoderCache.RUnlock()
															- f 不为空 return f 
															- f 为空
																- encoderCache.Lock()
																- encoderCache.m[t] = func(e *encodeState, v reflect.Value, opts encOpts) {
																							wg.Wait()
																							f(e, v, opts)
																						}
																- encoderCache.Unlock()
																- f = newTypeEncoder(t, true)
																		- func newTypeEncoder(t reflect.Type, allowAddr bool) encoderFunc
																			- //如果t满足marshalerType接口
																			- if t.Implements(marshalerType) -return marshalerEncoder
																														- func marshalerEncoder(e *encodeState, v reflect.Value, opts encOpts)
																															- if v.Kind()==reflect.Ptr && v.IsNil() -> e.WriteString("null") -> return 
																															- //还原m为本体
																															-  m, ok := v.Interface().(Marshaler) 
																															- b, err := m.MarshalJSON()
																															- // copy JSON into buffer, checking validity.
																															- //待看
																															- err = compact(&e.Buffer, b, opts.escapeHTML)
																																	- func compact(dst *bytes.Buffer, src []byte, escape bool) error
																																		- origLen := dst.Len()
																			- if t.Kind()!=reflect.Ptr && allowAddr
																				- //如果t的指针实现了marshalerType
																				- if reflect.PtrTo(t).Implements(marshalerType)
																					- return newCondAddrEncoder(addrMarshalerEncoder, newTypeEncoder(t, false))
																							- func addrMarshalerEncoder(e *encodeState, v reflect.Value, _ encOpts)
																								- va := v.Addr()
																								- if va.IsNil() -> e.WriteString("null") -> return 
																								- m := va.Interface().(Marshaler)
																								- b, err := m.MarshalJSON()																				- err = compact(&e.Buffer, b, true)
																			- ... 
																			- //如果t满足textMarshalerType接口
																			- if t.Implements(textMarshalerType) - return textMarshalerEncoder
																			- ...  
																			- switch k.Kind() case reflect.Struct: 
																				- return newStructEncoder(t)
																						- func newStructEncoder(t reflect.Type) encoderFunc
																							- fields := cachedTypeFields(t)
																									- func cachedTypeFields(t reflect.Type) []field 
																										- //使用atomic.Value的Load
																										- m, _ := fieldCache.value.Load().(map[reflect.Type][]field)
																										- f:=m[t]
																										- if f !=nil -> return f 
																										- f = typeFields(t)
																													- func typeFields(t reflect.Type) []field 
																														- fieldCache.mu.Lock()
																														- m, _ = fieldCache.value.Load().(map[reflect.Type][]field)
																														- newM := make(map[reflect.Type][]field, len(m)+1)
																														- for k, v := range m -> newM[k]=v 
																														- newM[t] = f //设置t的方法f
																														- fieldCache.value.Store(newM) //存储到内存中
																														- fieldCache.mu.Unlock()
																										- if f ==nil
																											- f=[]field{}


																							- 	se := &structEncoder{
																										fields:    fields,
																										fieldEncs: make([]encoderFunc, len(fields)),
																									}
																							- for i,f:=range fields 
																								- se.fieldEncs[i]=typeEncoder(typeByIndex(t, f.index))
																							- return se.encode
																									- func (se *structEncoder) encode(e *encodeState, v reflect.Value, opts encOpts) 
																										- e.WriteByte('{')
																										- first := true
																										- for i, f := range se.fields
																											- //递归
																											- fv := fieldByIndex(v, f.index)
																													- func fieldByIndex(v reflect.Value, index []int) reflect.Value 
																														- 
																											- ...  
																											- e.string(f.name, opts.escapeHTML)
																											- e.WriteByte(':')
																											- opts.quoted = f.quoted
																											- //递归
																											- se.fieldEncs[i](e, fv, opts)
																										- e.WriteByte('}')
																										

							
							
																- wg.Done()
																- encoderCache.Lock()
																- encoderCache.m[t] = f
																- encoderCache.Unlock()

			- return e.Bytes(), nil





//需要重复看的
func (e *encodeState) marshal(v interface{}, opts encOpts) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if s, ok := r.(string); ok {
				panic(s)
			}
			err = r.(error) //err被defer里定义
		}
	}()
	e.reflectValue(reflect.ValueOf(v), opts)
	return nil
}


func typeEncoder(t reflect.Type) encoderFunc {
	encoderCache.RLock()
	f := encoderCache.m[t]
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
		encoderCache.m = make(map[reflect.Type]encoderFunc)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	encoderCache.m[t] = func(e *encodeState, v reflect.Value, opts encOpts) {
		wg.Wait()
		f(e, v, opts)
	}
	encoderCache.Unlock()

	// Compute fields without lock.
	// Might duplicate effort but won't hold other computations back.
	f = newTypeEncoder(t, true)
	wg.Done()
	encoderCache.Lock()
	encoderCache.m[t] = f
	encoderCache.Unlock()
	return f
}


 atomic.Value ???

//主要函数
 func newTypeEncoder(t reflect.Type, allowAddr bool) encoderFunc {
	if t.Implements(marshalerType) {
		return marshalerEncoder
	}
	if t.Kind() != reflect.Ptr && allowAddr {
		if reflect.PtrTo(t).Implements(marshalerType) {
			return newCondAddrEncoder(addrMarshalerEncoder, newTypeEncoder(t, false))
		}
	}

	if t.Implements(textMarshalerType) {
		return textMarshalerEncoder
	}
	if t.Kind() != reflect.Ptr && allowAddr {
		if reflect.PtrTo(t).Implements(textMarshalerType) {
			return newCondAddrEncoder(addrTextMarshalerEncoder, newTypeEncoder(t, false))
		}
	}

	switch t.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncoder
	case reflect.Float32:
		return float32Encoder
	case reflect.Float64:
		return float64Encoder
	case reflect.String:
		return stringEncoder
	case reflect.Interface:
		return interfaceEncoder
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	case reflect.Ptr:
		return newPtrEncoder(t)
	default:
		return unsupportedTypeEncoder
	}
}


//实现error方法
type MarshalerError struct {
	Type reflect.Type
	Err  error
}

func (e *MarshalerError) Error() string {
	return "json: error calling MarshalJSON for type " + e.Type.String() + ": " + e.Err.Error()
}

...  
	if err != nil {
		e.error(&MarshalerError{v.Type(), err})
	}


//函数值


//递归
func fieldByIndex(v reflect.Value, index []int) reflect.Value {
	for _, i := range index {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return reflect.Value{}
			}
			v = v.Elem()
		}
		v = v.Field(i)
	}
	return v
}

copy(index, f.index)