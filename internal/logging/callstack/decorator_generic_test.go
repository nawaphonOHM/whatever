package callstack

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testGeneric42  = 42
	testGeneric100 = 100
	testWorkers30  = 30
)

func testDecorateFunc(t *testing.T) {
	f1 := DecorateFunc("calc", func() int {
		assert.Equal(t, 1, DefaultStack().Len())
		return testGeneric42
	})
	assert.Equal(t, testGeneric42, f1())
	assert.Equal(t, 0, DefaultStack().Len())
}

func testDecorateFuncErr(t *testing.T) {
	f2 := DecorateFuncErr("fetch", func() (string, error) {
		assert.Equal(t, 1, DefaultStack().Len())
		return "val", nil
	})
	val, err := f2()
	assert.NoError(t, err)
	assert.Equal(t, "val", val)
	assert.Equal(t, 0, DefaultStack().Len())
}

func testDecorateContextFunc(t *testing.T) {
	f3 := DecorateContextFunc("ctxFetch", func(ctx context.Context) (int, error) {
		cs := FromContext(ctx)
		require.NotNil(t, cs)
		assert.Equal(t, 1, cs.Len())
		return testGeneric100, nil
	})
	num, err := f3(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, testGeneric100, num)
}

func testDecorateCtx(t *testing.T) {
	f4 := DecorateCtx("ctxVal", func(ctx context.Context) string {
		cs := FromContext(ctx)
		require.NotNil(t, cs)
		assert.Equal(t, 1, cs.Len())
		return "result"
	})
	res := f4(context.Background())
	assert.Equal(t, "result", res)
}

func TestGenericDecorators(t *testing.T) {
	SetDefaultStack(New())
	t.Run("DecorateFunc", testDecorateFunc)
	t.Run("DecorateFuncErr", testDecorateFuncErr)
	t.Run("DecorateContextFunc", testDecorateContextFunc)
	t.Run("DecorateCtx", testDecorateCtx)
}

func TestDecorateContext_Concurrency(t *testing.T) {
	var wg sync.WaitGroup

	for range testWorkers30 {
		wg.Go(func() {
			fn := DecorateContext("concurTask", func(ctx context.Context) error {
				cs := FromContext(ctx)
				require.NotNil(t, cs)
				assert.Equal(t, 1, cs.Len())
				return nil
			})
			require.NoError(t, fn(context.Background()))
		})
	}
	wg.Wait()
}
