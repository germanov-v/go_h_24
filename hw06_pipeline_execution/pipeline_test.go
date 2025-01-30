package hw06pipelineexecution

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

var isFullTesting = true

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})
}

func TestAllStageStop(t *testing.T) {
	if !isFullTesting {
		return
	}
	wg := sync.WaitGroup{}
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		wg.Wait()

		require.Len(t, result, 0)

	})
}

func TestAdditionalPipeline(t *testing.T) {
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Empty additional", func(v interface{}) interface{} { return v }),
		g("Additional икс 2", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Additional +100", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Itoa", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("empty staging", func(t *testing.T) {
		in := make(Bi)
		data := []int{10, 20, 30}
		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()
		out := ExecutePipeline(in, nil)
		var result []interface{}
		for val := range out {
			result = append(result, val)
		}
		require.Equal(t, []interface{}{10, 20, 30}, result)
	})

	t.Run("empty IN", func(t *testing.T) {
		in := make(Bi)
		close(in)
		out := ExecutePipeline(in, nil, stages...)
		var result []interface{}
		for val := range out {
			result = append(result, val)
		}
		require.Empty(t, result)
	})

	t.Run("1 stage", func(t *testing.T) {
		in := make(Bi)
		go func() {
			in <- 5
			in <- 10
			in <- 15
			close(in)
		}()
		singleStage := g("x2", func(v interface{}) interface{} { return v.(int) * 2 })
		out := ExecutePipeline(in, nil, singleStage)
		var result []int
		for val := range out {
			result = append(result, val.(int))
		}
		require.Equal(t, []int{10, 20, 30}, result)
	})

	t.Run("max In", func(t *testing.T) {
		in := make(Bi)
		dataSize := 20 //1_000
		go func() {
			for i := 0; i < dataSize; i++ {
				in <- i
			}
			close(in)
		}()
		out := ExecutePipeline(in, nil, stages...)
		var result []string
		for val := range out {
			result = append(result, val.(string))
		}
		require.Len(t, result, dataSize)
		require.Equal(t, "100", result[0])
		require.Equal(t, "102", result[1])
	})

}
