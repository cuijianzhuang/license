package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed build
var embeddedFiles embed.FS

func main() {
	// 设置API服务端口
	apiPort := "15000"
	
	// 使用环境变量覆盖默认端口（如果有）
	if os.Getenv("API_PORT") != "" {
		apiPort = os.Getenv("API_PORT")
	}

	// 设置路由
	mux := http.NewServeMux()

	// 处理API请求（这里集成你的API处理器）
	// apiHandler := yourAPIHandler()
	// mux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	// 处理前端静态文件
	frontendFS, err := fs.Sub(embeddedFiles, "build")
	if err != nil {
		log.Fatal("无法加载前端文件:", err)
	}

	// 静态文件处理器
	staticHandler := http.FileServer(http.FS(frontendFS))

	// 处理请求路径中的静态资源
	mux.Handle("/static/", staticHandler)
	mux.Handle("/assets/", staticHandler)

	// 处理其他文件（favicon、manifest等）
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 检查是否是请求特定文件
		if filepath.Ext(r.URL.Path) != "" {
			staticHandler.ServeHTTP(w, r)
			return
		}

		// 对于所有其他请求（包括SPA路由），返回index.html
		indexFile, err := frontendFS.Open("index.html")
		if err != nil {
			http.Error(w, "无法加载页面", http.StatusInternalServerError)
			return
		}
		defer indexFile.Close()

		stat, err := indexFile.Stat()
		if err != nil {
			http.Error(w, "无法加载页面", http.StatusInternalServerError)
			return
		}

		buffer := make([]byte, stat.Size())
		_, err = indexFile.Read(buffer)
		if err != nil {
			http.Error(w, "无法加载页面", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(buffer)
	})

	// 启动HTTP服务
	log.Printf("服务器已启动，端口: %s\n", apiPort)
	log.Fatal(http.ListenAndServe(":"+apiPort, mux))
} 