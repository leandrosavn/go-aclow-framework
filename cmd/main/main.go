package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-aclow-framework/pkg/main/api"
	"github.com/go-aclow-framework/pkg/main/api/config"
	"github.com/go-aclow-framework/pkg/main/auth"
	pConfig "github.com/go-aclow-framework/pkg/main/config"
	"github.com/go-aclow-framework/pkg/main/core"
	"github.com/go-aclow-framework/pkg/main/types"
	"github.com/go-aclow-framework/pkg/main/utils"
	"github.com/joho/godotenv"
	"github.com/lfigueiredo82/aclow"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var (
	corsAllowHeaders     = "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization"
	corsAllowMethods     = "HEAD,GET,POST,PUT,DELETE,OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "false"
)

func CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)

		next(ctx)
	}
}

func main() {

	godotenv.Load()
	os.Setenv("TZ", "America/Sao_Paulo")
	os.Setenv("LOG_APP", "MAIN")

	startOpt := aclow.StartOptions{
		Debug: false,
		Local: true,
	}

	var app = &aclow.App{}

	app.Start(startOpt)

	//connectOnMongo(app)

	app.RegisterModule("api", api.Nodes())

	app.RegisterModule("auth", auth.Nodes())

	app.RegisterModule("core", core.Nodes())

	router := routing.New()
	StartRouter(app, router)

	panic(fasthttp.ListenAndServe(":8090", CORS(router.HandleRequest)))

}

func StartRouter(app *aclow.App, r *routing.Router) {

	routeConfig := config.ConfigRoutes()
	restConfig := routeConfig.Rest["routes"]
	for route, address := range restConfig {
		var method string = strings.TrimSpace(route[0:7])
		var path string = strings.TrimSpace(route[7:])
		fmt.Println("Registering Route...", method, path, address)
		registerRoute(app, r, method, path, address)
	}

}

func registerRoute(application *aclow.App, r *routing.Router, method string, path string, address string) {
	r.To(method, path, func(ctx *routing.Context) error {

		authorization := string(ctx.Request.Header.Peek("Authorization"))
		var headers = make(map[string]interface{}, 0)
		headers["address"] = address
		if authorization != "" {
			headers["token"] = strings.Trim(strings.ReplaceAll(authorization, "Bearer ", ""), "")
		}

		reply, err := application.Call("api@calling-flow", aclow.Message{Header: headers, Body: ctx.PostBody()})
		if err != nil {
			if err, ok := err.(*utils.MessageError); ok {

				ctx.SetContentType("application/json")
				ctx.Response.Header.Set("Content-Type", "application/json;charset=UTF-8")
				ctx.SetStatusCode(err.StatusCode)
				fmt.Fprintf(ctx, err.ErrorMessage)
			}
			return nil
		}

		isFileType, fType := isFile(reply.Body)
		if isFileType {
			var fileType types.AbstractFile
			if fType == "xlsx" {
				fileType = reply.Body.(types.XlsxFileType)
			} else {
				fileType = reply.Body.(types.GenericFileType)
			}
			fileName := utils.NameOf(fileType.GetFileName())
			ctx.Response.Header.Set("Content-Description", "File Transfer")
			ctx.Response.Header.Set("Content-Type", "application/octet-stream")
			ctx.Response.Header.Set("Content-Disposition", "attachment; filename="+fileName)
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.Response.SendFile(fileType.GetFileName())
		} else {
			ctx.SetContentType("application/json")
			ctx.Response.Header.Set("Content-Type", "application/json;charset=UTF-8")
			ctx.SetStatusCode(fasthttp.StatusOK)
			fmt.Fprintf(ctx, reply.Body.(string))
		}

		return err

	})

}

func isFile(data interface{}) (bool, string) {
	_, ok := interface{}(data).(interface {
		FileType() string
	})
	if ok {
		return ok, data.(types.AbstractFile).FileType()
	}
	return ok, ""
}

var mongoTries = 0

func connectOnMongo(app *aclow.App) {
	log.Println("Connecting on MongoDB...")
	opt := &options.ClientOptions{}
	opt = opt.ApplyURI(pConfig.MongoDbDsn())
	opt = opt.SetMaxPoolSize(pConfig.MongoDbMaxPoolSize())
	client, err := mongo.NewClient(opt)
	if err != nil {
		log.Panic(err)
	}

	err = client.Connect(context.TODO())

	if err != nil {
		time.Sleep(time.Second * 1)
		log.Println(err.Error())
		mongoTries++
		if mongoTries > 10 {
			log.Panic(err)
		}
		log.Println("Trying again...")
		connectOnMongo(app)
	} else {
		app.Resources["mongo"] = client
	}
}
