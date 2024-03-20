package bootstrap

import (
	"os"
	"time"

	"github.com/Zawa-ll/raffle/conf"
	"github.com/Zawa-ll/raffle/cron"
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/sessions"
)

type Configurator func(bootstrapper *Bootstrapper)

type Bootstrapper struct {
	*iris.Application // BS instance has all properties and method of an iris.Application instnace
	AppName           string
	AppOwner          string
	AppSpawnDate      time.Time // Stores the date and time when the application was started

	Sessions *sessions.Sessions
}

// Generate a new Bootstrapper Applicatio
func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper {
	b := &Bootstrapper{ // base on the created Struct, create a bootstrap instance
		Application:  iris.New(),
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnDate: time.Now(),
	}

	for _, cfg := range cfgs {
		cfg(b) // each cfg is of type Configurator Function Type
	}

	return b
}

// set up for rendering views in Iris web application
func (b *Bootstrapper) SetupViews(viewDir string) { // viewDir specifies where the template can be found
	htmlEngine := iris.HTML(viewDir, ".html").Layout("shared/layout.html") // setting base layout for all templates
	// Reload Template Again at Each Time
	htmlEngine.Reload(true)

	// Time Formatting Functions
	htmlEngine.AddFunc("FormUnixtimeShort", func(t int) string {
		dt := time.Unix(int64(t), int64(0))
		return dt.Format(conf.SysTimeformShort)
	})

	htmlEngine.AddFunc("FormUnixtime", func(t int) string {
		dt := time.Unix(int64(t), int64(0))
		return dt.Format(conf.SysTimeform)
	})

	b.RegisterView(htmlEngine)
}

// SetupSessions initializes the sessions
// For storing information across multiple HTTP requests made by the same client
func (b *Bootstrapper) SetupSessions(expires time.Duration, cookieHashKey, cookieBlockKey []byte) {
	b.Sessions = sessions.New(sessions.Config{
		Cookie:   "SECRET_SESS_COOKIE_" + b.AppName,
		Expires:  expires,
		Encoding: securecookie.New(cookieHashKey, cookieBlockKey),
	})
}

// SetupErrorHandlers prepares the http error handlers
// `(context.StatusCodeNotSuccessful`,  which defaults to < 200 || >= 400 but you can change it).
func (b *Bootstrapper) SetupErrorHandlers() {
	b.OnAnyErrorCode(func(ctx iris.Context) {
		// Error Handler
		err := iris.Map{
			"app":     b.AppName,
			"status":  ctx.GetStatusCode(),
			"message": ctx.Values().GetString("message"),
		}

		if jsonOutput := ctx.URLParamExists("json"); jsonOutput {
			ctx.JSON(err)
			return
		}

		// ViewData(k, v) passed error details to view template
		ctx.ViewData("Err", err)       // Passes the error details (err) to the view template under the key "Err"
		ctx.ViewData("Title", "Error") // Sets the page title to "Error", making it available in the view template under the key "Title".
		ctx.View("shared/error.html")  // Render the view using error.html template
	})
}

// Configure accepts configurations and runs them inside the Bootstraper's context.
// Receiver Function
func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}

// Initiate planned mission services
func (b *Bootstrapper) setupCron() {
	// Service class applications
	if conf.RunningCrontabService {
		// start cron job for reseting all gift prize data every 5 minutes
		//  and distribute prizes into the pool every minute
		cron.ConfigueAppOneCron()
	}
	cron.ConfigueAppAllCron()
}

const (
	// StaticAssets is the root directory for public assets like images, css, js.
	StaticAssets = "./public"
	// Favicon is the relative 9to the "StaticAssets") favicon path for our app.
	Favicon = "/favicon.ico"
)

// Bootstrap prepares our application.
//
// Returns itself.
func (b *Bootstrapper) Bootstrap() *Bootstrapper {
	b.SetupViews("./views")
	b.SetupSessions(24*time.Hour,
		[]byte("the-big-and-secret-fash-key-here"),
		[]byte("lot-secret-of-characters-big-too"),
	)
	b.SetupErrorHandlers()

	// static files
	b.Favicon(StaticAssets + Favicon)
	b.HandleDir(StaticAssets[1:], StaticAssets)
	indexHtml, err := os.ReadFile(StaticAssets + "/index.html")
	if err == nil {
		b.StaticContent(StaticAssets[1:]+"/", "text/html",
			indexHtml)
	}
	// Don't leave out the "/" at the end of the catalog
	iris.WithoutPathCorrectionRedirection(b.Application)

	// crontab
	b.setupCron()

	// middleware, after static files
	b.Use(recover.New())
	b.Use(logger.New())

	return b
}

// Listen starts the http server with the specified "addr".
func (b *Bootstrapper) Listen(addr string, cfgs ...iris.Configurator) {
	b.Run(iris.Addr(addr), cfgs...)
}
