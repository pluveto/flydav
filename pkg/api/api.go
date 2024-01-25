// pkg/api/middleware.go
package api

import (
    "net/http"
)

// Middleware 是一个用于定义中间件的类型
type Middleware func(http.Handler) http.Handler

// 使用示例中间件
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 实现日志记录的逻辑
        next.ServeHTTP(w, r)
    })
}
