package log

import (
	"context"
	"errors"
	"testing"
)

func TestMain(t *testing.M) {
	ctx := context.Background()
	EventInfo(ctx, "test_event_info", "this_is_key", "this_is_val")
	ctx = context.WithValue(ctx, LogIDKey, "this_is_log_id")
	EventInfo(ctx, "test_event_info", "this_is_key", "this_is_val")
	EventDebug(ctx, "test_event_debug", "this_is_key", "this_is_val")
	EventWarn(ctx, "test_event_warn", "this_is_key", "this_is_val", "this_is_key", 12345)
	EventError(ctx, errors.New("this is error"), "test_event_error", "this_is_key", "this_is_val")
	EventError(ctx, nil, "test_event_nil_error", "this_is_key", "this_is_val")
}
