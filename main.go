package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/cloudrunv2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")
		mongodbUri := conf.RequireSecret("mongodbUri")
		jwtSecret := conf.RequireSecret("jwtSecret")
		redisAddress := conf.RequireSecret("redisAddress")
		redisPassword := conf.RequireSecret("redisPassword")
		gmailEmail := conf.RequireSecret("gmailEmail")
		gmailEmailPassword := conf.RequireSecret("gmailEmailPassword")
		payosClientID := conf.RequireSecret("payosClientID")
		payosApiKey := conf.RequireSecret("payosApiKey")
		payosChecksum := conf.RequireSecret("payosChecksum")
		dockerImageLink := conf.RequireSecret("dockerImageLink")

		service, err := cloudrunv2.NewService(ctx, "backend-service", &cloudrunv2.ServiceArgs{
			Ingress:     pulumi.String("INGRESS_TRAFFIC_ALL"),
			LaunchStage: pulumi.String("GA"),
			Location:    pulumi.String("asia-southeast1"),
			Name:        pulumi.String("go-backend-service"),
			Project:     pulumi.String("bebanbong"),
			Template: &cloudrunv2.ServiceTemplateArgs{
				Containers: cloudrunv2.ServiceTemplateContainerArray{
					&cloudrunv2.ServiceTemplateContainerArgs{
						Envs: cloudrunv2.ServiceTemplateContainerEnvArray{
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name: pulumi.String("API_ADDR"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name: pulumi.String("FE_ADDR"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("MONGODB_URI"),
								Value: mongodbUri,
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("JWT_SECRET"),
								Value: jwtSecret,
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("REDIS_ADDR"),
								Value: redisAddress,
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("REDIS_PASSWORD"),
								Value: redisPassword,
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("GMAIL_EMAIL"),
								Value: gmailEmail,
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("GMAIL_EMAIL_PASSWORD"),
								Value: gmailEmailPassword,
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("PAYOS"),
								Value: pulumi.String("https://api-merchant.payos.vn/"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("PAYOS_CLIENT_ID"),
								Value: payosClientID,
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("PAYOS_API_KEY"),
								Value: payosApiKey,
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("PAYOS_CHECKSUM"),
								Value: payosChecksum,
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("CLOUDINARY_URL"),
								Value: pulumi.String("cloudinary://<api_key>:<api_secret>@<cloud_name>"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("SMTP_HOST"),
								Value: pulumi.String("smtp.gmail.com"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("SMTP_PORT"),
								Value: pulumi.String("587"),
							},
							&cloudrunv2.ServiceTemplateContainerEnvArgs{
								Name:  pulumi.String("GIN_MODE"),
								Value: pulumi.String("release"),
							},
						},
						Image: dockerImageLink,
						Ports: cloudrunv2.ServiceTemplateContainerPortArray{
							&cloudrunv2.ServiceTemplateContainerPortArgs{
								ContainerPort: pulumi.Int(8080),
								Name:          pulumi.String("http1"),
							},
						},
						Resources: &cloudrunv2.ServiceTemplateContainerResourcesArgs{
							CpuIdle: pulumi.Bool(true),
							Limits: pulumi.StringMap{
								"cpu":    pulumi.String("1000m"),
								"memory": pulumi.String("512Mi"),
							},
						},
						StartupProbe: &cloudrunv2.ServiceTemplateContainerStartupProbeArgs{
							FailureThreshold: pulumi.Int(1),
							PeriodSeconds:    pulumi.Int(240),
							TcpSocket: &cloudrunv2.ServiceTemplateContainerStartupProbeTcpSocketArgs{
								Port: pulumi.Int(8080),
							},
							TimeoutSeconds: pulumi.Int(240),
						},
					},
				},
				MaxInstanceRequestConcurrency: pulumi.Int(80),
				ServiceAccount:                pulumi.String("743300054599-compute@developer.gserviceaccount.com"),
				Timeout:                       pulumi.String("300s"),
			},
			Traffics: cloudrunv2.ServiceTrafficArray{
				&cloudrunv2.ServiceTrafficArgs{
					Percent: pulumi.Int(100),
					Type:    pulumi.String("TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"),
				},
			},
		}, pulumi.Protect(true))
		if err != nil {
			return err
		}

		_, err = cloudrunv2.NewServiceIamMember(ctx, "backend-service-public-access", &cloudrunv2.ServiceIamMemberArgs{
			Project:  pulumi.String("bebanbong"),
			Location: pulumi.String("asia-southeast1"),
			Name:     service.Name,
			Role:     pulumi.String("roles/run.invoker"),
			Member:   pulumi.String("allUsers"),
		})
		if err != nil {
			return err
		}
		ctx.Export("uri", service.Uri)
		return nil
	})
}
