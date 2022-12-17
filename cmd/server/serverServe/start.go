package serverServe

import (
	stdContext "context"
	"fmt"
	"github.com/aaronchen2k/deeptest"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1"
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	"github.com/aaronchen2k/deeptest/internal/pkg/config"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/cron"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/middleware"
	"github.com/aaronchen2k/deeptest/internal/pkg/core/module"
	"github.com/aaronchen2k/deeptest/internal/pkg/log"
	commUtils "github.com/aaronchen2k/deeptest/internal/pkg/utils"
	"github.com/aaronchen2k/deeptest/internal/server/core/cache"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	"github.com/aaronchen2k/deeptest/internal/server/modules/service"
	_i118Utils "github.com/aaronchen2k/deeptest/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/facebookgo/inject"
	"github.com/kataras/iris/v12/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/helper/tests"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

var client *tests.Client

// Start 初始化web服务
func Start() {
	inits()

	idleConnClosed := make(chan struct{})
	irisApp := createIrisApp(&idleConnClosed)

	//irisApp.Use(func(ctx iris.Context) {
	//	ctx.Request().Header.Del("Origin")
	//	ctx.Next()
	//})

	initWebSocket(irisApp)

	server := &WebServer{
		app:               irisApp,
		addr:              config.CONFIG.System.ServerAddress,
		timeFormat:        config.CONFIG.System.TimeFormat,
		staticPath:        config.CONFIG.System.StaticPath,
		webPath:           config.CONFIG.System.WebPath,
		idleConnClosed:    idleConnClosed,
		globalMiddlewares: []context.Handler{},
	}

	server.InjectModule()
	server.Start()
}

func inits() {
	consts.RunFrom = consts.FromServer
	consts.WorkDir = commUtils.GetWorkDir()

	config.Init("server")
	zapLog.Init("server")
	_i118Utils.Init(consts.Language, "")

	err := cache.Init()
	if err != nil {
		logUtils.Errorf("init redis cache failed, error %s", err.Error())
		return
	}
}

func createIrisApp(idleConnClosed *chan struct{}) (irisApp *iris.Application) {
	irisApp = iris.New()
	irisApp.Validator = validator.New() //参数验证
	irisApp.Logger().SetLevel(config.CONFIG.System.Level)

	iris.RegisterOnInterrupt(func() { //优雅退出
		timeout := 10 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		irisApp.Shutdown(ctx) // close all hosts

		close(*idleConnClosed)
	})

	return
}

// injectWebsocketModule 注册组件
func injectWebsocketModule(websocketCtrl *handler.WebSocketCtrl) {
	var g inject.Graph
	g.Logger = logrus.StandardLogger()

	if err := g.Provide(
		&inject.Object{Value: dao.GetDB()},
		&inject.Object{Value: websocketCtrl},
	); err != nil {
		logrus.Fatalf("provide usecase objects to the Graph: %v", err)
	}
	err := g.Populate()
	if err != nil {
		logrus.Fatalf("populate the incomplete Objects: %v", err)
	}
}

// Init 启动web服务
func (webServer *WebServer) Start() {
	webServer.app.UseGlobal(webServer.globalMiddlewares...)
	err := webServer.InitRouter()
	if err != nil {
		fmt.Printf("初始化路由错误： %v\n", err)
		panic(err)
	}
	// 添加上传文件路径
	webServer.app.Listen(
		webServer.addr,
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(webServer.timeFormat),
	)
	<-webServer.idleConnClosed
}

// GetAddr 获取web服务地址
func (webServer *WebServer) GetAddr() string {
	return webServer.addr
}

// AddModule 添加模块
func (webServer *WebServer) AddModule(module ...module.WebModule) {
	webServer.modules = append(webServer.modules, module...)
}

// AddWebUi 添加前端页面访问
func (webServer *WebServer) AddWebUi() {
	uiFs, err := deeptest.GetUiFileSys()
	if err != nil {
		return
	}

	webServer.app.HandleDir("/", http.FS(uiFs), iris.DirOptions{
		IndexName: "index.html",
		ShowList:  false,
		SPA:       true,
	})
}

// AddUploadStatic 添加上传文件访问
func (webServer *WebServer) AddUploadStatic() {
	fsOrDir := iris.Dir(filepath.Join(dir.GetCurrentAbPath(), filepath.Join(webServer.staticPath, "upload")))
	webServer.addStatic("/upload", fsOrDir)
}

// AddTestStatic 添加测试文件访问
func (webServer *WebServer) AddTestStatic() {
	fsOrDir := iris.Dir(filepath.Join(dir.GetCurrentAbPath(), filepath.Join(webServer.staticPath, "test")))
	webServer.addStatic("/test", fsOrDir)
}

func (webServer *WebServer) addStatic(requestPath string, fsOrDir interface{}, opts ...iris.DirOptions) {
	webServer.app.HandleDir(requestPath, fsOrDir, opts...)
}

// GetModules 获取模块
func (webServer *WebServer) GetModules() []module.WebModule {
	return webServer.modules
}

// GetTestAuth 获取测试验证客户端
func (webServer *WebServer) GetTestAuth(t *testing.T) *tests.Client {
	var once sync.Once
	once.Do(
		func() {
			client = tests.New(str.Join("http://", webServer.addr), t, webServer.app)
			if client == nil {
				t.Fatalf("client is nil")
			}
		},
	)

	return client
}

// GetTestLogin 测试登录web服务
func (webServer *WebServer) GetTestLogin(t *testing.T, url string, res tests.Responses, datas ...map[string]interface{}) *tests.Client {
	client := webServer.GetTestAuth(t)
	err := client.Login(url, res, datas...)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// Init 加载模块
func initWebSocket(irisApp *iris.Application) {
	websocketCtrl := handler.NewWebsocketCtrl()
	injectWebsocketModule(websocketCtrl)

	mvc.New(irisApp)

	websocketAPI := irisApp.Party(consts.WsPath)
	m := mvc.New(websocketAPI)
	m.Register(
		&service.PrefixedLogger{Prefix: ""},
	)
	m.HandleWebsocket(websocketCtrl)
	websocketServer := websocket.New(
		middleware.DefaultUpgrader,
		//gorilla.Upgrader(gorillaWs.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}),
		m)

	websocketAPI.Get("/", websocket.Handler(websocketServer))
}

// Init 加载模块
func (webServer *WebServer) InjectModule() {
	var g inject.Graph
	g.Logger = logrus.StandardLogger()

	cron := cron.NewServerCron()
	cron.Init()
	indexModule := v1.NewIndexModule()

	// inject objects
	if err := g.Provide(
		&inject.Object{Value: dao.GetDB()},
		&inject.Object{Value: cron},
		&inject.Object{Value: indexModule},
	); err != nil {
		logrus.Fatalf("provide usecase objects to the Graph: %v", err)
	}
	err := g.Populate()
	if err != nil {
		logrus.Fatalf("populate the incomplete Objects: %v", err)
	}

	webServer.AddModule(indexModule.Party())
}
