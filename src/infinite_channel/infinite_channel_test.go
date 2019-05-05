package infinite_channel

import (
	"testing"
)

func TestInfiniteChannel(t *testing.T) {
	ch := NewInfiniteChannel()
	testChannel(t, "infinite channel", ch)

	ch = NewInfiniteChannel()
	testChannelPair(t, "infinite channel", ch, ch)

	ch = NewInfiniteChannel()
	testChannelConcurrentAccessors(t, "infinite channel", ch)

}

func testChannel(t *testing.T, name string, ch *InfiniteChannel) {
	go func() {
		for i := 0; i < 1000; i++ {
			ch.In() <- i
		}
		ch.Close()
	}()
	for i := 0; i < 1000; i++ {
		val := <-ch.Out()
		if i != val.(int) {
			t.Fatal(name, "expected", i, "but got", val.(int))
		}
	}
}

func testChannelPair(t *testing.T, name string, in *InfiniteChannel, out *InfiniteChannel) {
	go func() {
		for i := 0; i < 1000; i++ {
			in.In() <- i
		}
		in.Close()
	}()
	for i := 0; i < 1000; i++ {
		val := <-out.Out()
		if i != val.(int) {
			t.Fatal("pair", name, "expected", i, "but got", val.(int))
		}
	}
}

func testChannelConcurrentAccessors(t *testing.T, name string, ch *InfiniteChannel) {
	// no asserts here, this is just for the race detector's benefit
	go ch.Len()
	go ch.Cap()

	go func() {
		ch.In() <- nil
	}()

	go func() {
		<-ch.Out()
	}()
}

func BenchmarkInfiniteChannelSerial(b *testing.B) {
	ch := NewInfiniteChannel()
	for i := 0; i < b.N; i++ {
		ch.In() <- nil
	}
	for i := 0; i < b.N; i++ {
		<-ch.Out()
	}
}

func BenchmarkInfiniteChannelParallel(b *testing.B) {
	ch := NewInfiniteChannel()
	go func() {
		for i := 0; i < b.N; i++ {
			<-ch.Out()
		}
		ch.Close()
	}()
	for i := 0; i < b.N; i++ {
		ch.In() <- nil
	}
	<-ch.Out()
}

func BenchmarkInfiniteChannelTickTock(b *testing.B) {
	ch := NewInfiniteChannel()
	for i := 0; i < b.N; i++ {
		ch.In() <- nil
		<-ch.Out()
	}
}
