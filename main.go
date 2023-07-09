package main

import (
	"fmt"
	"github.com/hub1989/keycloak-grpc-service/grpc/controller"
	"github.com/hub1989/keycloak-grpc-service/grpc/logger"
	"github.com/hub1989/keycloak-grpc-service/keycloak"
	"github.com/hub1989/keycloak-grpc-service/otel_config"
	user "github.com/hub1989/keycloak-protobuf/golang/keycloak"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)

	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}
}

func main() {

	grpcPort := os.Getenv("GRPC_PORT")

	if grpcPort == "" {
		panic("grpc port is not set")
	}

	tp, err := otel_config.TracerProvider(os.Getenv("TRACING_URL"))
	if err != nil {
		log.Error(err)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))

	configureGrpc(grpcPort)
}

func configureGrpc(grpcPort string) {
	lis, _ := net.Listen("tcp", grpcPort)
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(otelgrpc.UnaryServerInterceptor(), logger.ServerLogger))

	httpClient := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	configuration := keycloak.DefaultKeycloakConfiguration{
		BaseURL: os.Getenv("KEYCLOAK_URL"),
		Realm:   os.Getenv("KEYCLOAK_REALM"),
		Client:  httpClient,
	}
	clientService := keycloak.DefaultClientService{Configuration: configuration}
	credentialService := keycloak.DefaultCredentialService{Configuration: configuration}
	groupService := keycloak.DefaultGroupService{Configuration: configuration}
	roleService := keycloak.DefaultRoleService{Configuration: configuration, ClientService: clientService}
	userService := keycloak.DefaultUserService{Configuration: configuration}

	user.RegisterUserServiceServer(s, &controller.UserController{
		CredentialService: credentialService,
		UserService:       userService,
	})

	user.RegisterRoleServiceServer(s, &controller.RoleController{
		RoleService:       roleService,
		CredentialService: credentialService,
	})

	user.RegisterGroupServiceServer(s, &controller.GroupController{
		CredentialService: credentialService,
		GroupService:      groupService,
	})

	user.RegisterClientServiceServer(s, &controller.ClientController{
		ClientService:     clientService,
		CredentialService: credentialService,
	})

	if os.Getenv("ENV") == "development" {
		log.Info("in development environment, reflection is enabled")
		reflection.Register(s)
	}

	log.Info(fmt.Sprintf("running grpc on port %s", grpcPort))
	log.Fatal(s.Serve(lis))
}
