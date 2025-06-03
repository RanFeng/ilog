# ilog

## 使用方法

```
go get -u "github.com/RanFeng/ilog@v1.0.1"
```
## demo
``` go
import (
	"context"

	ilog "github.com/RanFeng/ilog"
)

func printLog(ctx context.Context, event string) {
	level := ilog.CtxLevel(ctx)
	ilog.EventInfo(ctx, event, "this_is_key", "this_is_val", "level", level)
	ilog.EventDebug(ctx, event, "this_is_key", "this_is_val", "level", level)
	ilog.EventWarn(ctx, event, "this_is_key", "this_is_val", "this_is_key", 12345, "level", level)
	ilog.EventError(ctx, errors.New("this is error"), event, "this_is_key", "this_is_val", "level", level)
	ilog.EventError(ctx, nil, event, "this_is_key", "this_is_val", "level", level)
	fmt.Println()
}


func main() {
	ctx := context.Background()
	printLog(ctx, "test_common")
	
	ctx = context.WithValue(ctx, LogIDKey, "this_is_log_id")
	printLog(ctx, "test_with_logid")
	
	ilog.SetGlobalLogLevel(LevelWarn)
	printLog(ctx, "test_with_global_level")
	
	ctx = ilog.SetCtxLogLevel(ctx, LevelInfo)
	printLog(ctx, "test_with_ctx_level")
}

```