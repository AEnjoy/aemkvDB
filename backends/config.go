package backends

import "context"

// UsingGlobalDB 当指定了这个参数为true, 则在程序空间内，无论创建几个DB对象, 都只会使用第一个DB对象.
var UsingGlobalDB = false

// GlobalContext 全局上下文
var GlobalContext context.Context

// Split 全局str分隔符
var Split = "###|*^*|###"

func init() {
	GlobalContext = context.Background()
}
