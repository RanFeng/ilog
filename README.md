# ilog

## 使用方法

```
go get -u "github.com/RanFeng/ilog@v1.0.0"
```
## demo
``` go
import (
	"context"

	ilog "github.com/RanFeng/ilog"
)

func main() {
	ctx := context.Background()
	ilog.EventInfo(ctx, "test_event_info", "this_is_key", "this_is_val")
	ctx = context.WithValue(ctx, ilog.LogIDKey, "this_is_log_id")
	ilog.EventInfo(ctx, "test_event_info", "this_is_key", "this_is_val")
	ilog.EventDebug(ctx, "test_event_debug", "this_is_key", "this_is_val")
	ilog.EventWarn(ctx, "test_event_warn", "this_is_key", "this_is_val", "this_is_key", 12345)
	ilog.EventError(ctx, errors.New("this is error"), "test_event_error", "this_is_key", "this_is_val")
	ilog.EventError(ctx, nil, "test_event_nil_error", "this_is_key", "this_is_val")
}

```