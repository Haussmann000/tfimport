// cmd/tfimport/main.go
package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/Haussmann000/tfimport/internal/di"
)

func main() {
	// コマンドラインフラグの定義
	resourceTypes := flag.String("resource-types", "vpc,s3", "Importするリソースタイプをカンマ区切りで指定 (e.g., vpc,s3)")
	resourceName := flag.String("resource-name", "", "Importするリソースの名前を前方一致で指定")
	flag.Parse()

	ctx := context.Background()

	// DIコンテナを使ってアプリケーションを構築
	app, err := di.BuildApp(ctx)
	if err != nil {
		log.Fatalf("failed to build application: %v", err)
	}

	runOptions := di.RunOptions{
		ResourceTypes: strings.Split(*resourceTypes, ","),
		ResourceName:  *resourceName,
	}

	// アプリケーションを実行
	if err := app.Run(ctx, runOptions); err != nil {
		log.Fatalf("application returned an error: %v", err)
	}
}
