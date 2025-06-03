package ilog

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func printLog(ctx context.Context, event string) {
	level := CtxLevel(ctx)
	EventInfo(ctx, event, "this_is_key", "this_is_val", "level", level)
	EventDebug(ctx, event, "this_is_key", "this_is_val", "level", level)
	EventWarn(ctx, event, "this_is_key", "this_is_val", "this_is_key", 12345, "level", level)
	EventError(ctx, errors.New("this is error"), event, "this_is_key", "this_is_val", "level", level)
	EventError(ctx, nil, event, "this_is_key", "this_is_val", "level", level)
	fmt.Println()
}

func TestMain(t *testing.M) {
	ctx := context.Background()
	printLog(ctx, "test_common")
	ctx = context.WithValue(ctx, LogIDKey, "this_is_log_id")
	printLog(ctx, "test_with_logid")
	SetGlobalLogLevel(LevelWarn)
	printLog(ctx, "test_with_global_level")
	ctx = SetCtxLogLevel(ctx, LevelInfo)
	printLog(ctx, "test_with_ctx_level")
}
