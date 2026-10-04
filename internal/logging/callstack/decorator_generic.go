package callstack

import (
	"context"
)

// DecorateFunc wraps a generic func() T.
func DecorateFunc[T any](name string, fn func() T) func() T {
	return func() T {
		frame := captureCallerFrame(name, 2)
		DefaultStack().Push(frame)
		defer DefaultStack().Pop()
		return fn()
	}
}

// DecorateFuncErr wraps a generic func() (T, error).
func DecorateFuncErr[T any](name string, fn func() (T, error)) func() (T, error) {
	return func() (T, error) {
		frame := captureCallerFrame(name, 2)
		DefaultStack().Push(frame)
		defer DefaultStack().Pop()
		return fn()
	}
}

// DecorateContextFunc wraps a generic func(context.Context) (T, error).
func DecorateContextFunc[T any](name string, fn func(context.Context) (T, error)) func(context.Context) (T, error) {
	return func(ctx context.Context) (T, error) {
		ctxWithStack, cs := EnsureContext(ctx)
		frame := captureCallerFrame(name, 2)
		cs.Push(frame)
		defer cs.Pop()
		return fn(ctxWithStack)
	}
}

// DecorateCtx wraps a generic func(context.Context) T.
func DecorateCtx[T any](name string, fn func(context.Context) T) func(context.Context) T {
	return func(ctx context.Context) T {
		ctxWithStack, cs := EnsureContext(ctx)
		frame := captureCallerFrame(name, 2)
		cs.Push(frame)
		defer cs.Pop()
		return fn(ctxWithStack)
	}
}
