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
	var resourceTypes, resourceName, clusterName, serviceName, securityGroupID, dbClusterIdentifier, dbInstanceIdentifier, bucketName string
	flag.StringVar(&resourceTypes, "resource-types", "", "aws resource type. s3, vpc, ecs, elbv2, iam, security_group, rds")
	flag.StringVar(&resourceName, "resource-name", "", "aws resource name (for vpc, elbv2, iam, rds parameter group)")
	flag.StringVar(&bucketName, "bucket-name", "", "s3 bucket name")
	flag.StringVar(&clusterName, "cluster-name", "", "ecs cluster name")
	flag.StringVar(&serviceName, "service-name", "", "ecs service name")
	flag.StringVar(&securityGroupID, "security-group-id", "", "comma separated security group ids")
	flag.StringVar(&dbClusterIdentifier, "db-cluster-identifier", "", "rds db cluster identifier")
	flag.StringVar(&dbInstanceIdentifier, "db-instance-identifier", "", "rds db instance identifier")

	flag.Parse()

	if resourceTypes == "" {
		log.Fatal("resource-types is required")
	}

	ctx := context.Background()
	app, err := di.BuildApp(ctx)
	if err != nil {
		log.Fatalf("failed to build app: %v", err)
	}

	options := di.RunOptions{
		ResourceTypes:        strings.Split(resourceTypes, ","),
		ResourceName:         resourceName,
		BucketName:           bucketName,
		ClusterName:          clusterName,
		ServiceName:          serviceName,
		SecurityGroupID:      securityGroupID,
		DBClusterIdentifier:  dbClusterIdentifier,
		DBInstanceIdentifier: dbInstanceIdentifier,
	}

	if err := app.Run(ctx, options); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
