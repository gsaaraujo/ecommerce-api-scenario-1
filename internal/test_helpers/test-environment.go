package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type TestEnvironment struct {
	baseUrl                string
	client                 *http.Client
	awsConfig              aws.Config
	pgxPool                *pgxpool.Pool
	redisClient            *redis.Client
	rabbitmqConn           *amqp091.Connection
	postgresContainerUrl   string
	localstackContainerUrl string
	wiremockContainerUrl   string
	redisContainerUrl      string
	rabbitmqContainerUrl   string
}

func NewTestEnvironment() *TestEnvironment {
	return &TestEnvironment{}
}

func (t *TestEnvironment) Start() {
	t.postgresContainerUrl = utils.GetOrThrow(NewPostgresContainer()).url
	t.localstackContainerUrl = utils.GetOrThrow(NewLocalstackContainer()).url
	t.wiremockContainerUrl = utils.GetOrThrow(NewWiremockContainer()).url
	t.redisContainerUrl = utils.GetOrThrow(NewRedisContainer()).url
	t.rabbitmqContainerUrl = utils.GetOrThrow(NewRabbitmqContainer()).url

	_ = os.Setenv("AWS_REGION", "us-east-1")
	_ = os.Setenv("AWS_ACCESS_KEY_ID", "test")
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	_ = os.Setenv("AWS_ENDPOINT_URL", t.localstackContainerUrl)
	_ = os.Setenv("AWS_SECRET_MANAGER_NAME", "secret-us-east-1-local-app")
	_ = os.Setenv("TERN_MIGRATIONS_PATH", "../migrations")
	_ = os.Setenv("ZIPCODE_URL", t.wiremockContainerUrl)

	t.awsConfig = utils.GetOrThrow(config.LoadDefaultConfig(context.TODO()))

	t.createSecrets()

	utils.ThrowOnError(t.runMigrations())

	t.pgxPool = utils.GetOrThrow(pgxpool.New(context.Background(), t.postgresContainerUrl))
	t.rabbitmqConn = utils.GetOrThrow(amqp091.Dial(t.rabbitmqContainerUrl))

	redisParsedUrl := utils.GetOrThrow(redis.ParseURL(t.redisContainerUrl))
	t.redisClient = redis.NewClient(&redis.Options{
		Addr:     redisParsedUrl.Addr,
		Password: redisParsedUrl.Password,
	})

	t.startHttpServer()
}

func (t *TestEnvironment) createSecrets() {
	secretsClient := secretsmanager.NewFromConfig(t.awsConfig)

	utils.GetOrThrow(secretsClient.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
		Name: aws.String("secret-us-east-1-local-app"),
		SecretString: aws.String(fmt.Sprintf(`
			{
				"REDIS_URL": "%s",
				"POSTGRES_URL": "%s",
				"RABBITMQ_URL": "%s",
				"MERCADO_PAGO_ACCESS_KEY": "",
				"ZIPCODE_TOKEN": "a7416146283d464294cebea38d5cb5ff",
				"ACCESS_TOKEN_SIGNING_KEY": "81c4a8d5b2554de4ba736e93255ba633"
			}
		`, t.redisContainerUrl, t.postgresContainerUrl, t.rabbitmqContainerUrl)),
	}))
}

func (t *TestEnvironment) runMigrations() error {
	urlParsed := utils.GetOrThrow(url.Parse(t.postgresContainerUrl))

	_ = os.Setenv("PGUSER", "postgres")
	_ = os.Setenv("PGPASSWORD", "postgres")
	_ = os.Setenv("PGHOST", urlParsed.Hostname())
	_ = os.Setenv("PGPORT", urlParsed.Port())
	_ = os.Setenv("PGDATABASE", "postgres")

	migrations := ""
	if _, ok := os.LookupEnv("TERN_MIGRATIONS_PATH"); !ok {
		return errors.New("TERN_MIGRATIONS_PATH environment variable not found")
	}

	migrations = os.Getenv("TERN_MIGRATIONS_PATH")
	cmd := exec.Command("tern", "migrate", "-m", migrations)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error: %s, output: %s", err.Error(), string(output))
	}

	return nil
}

func (t *TestEnvironment) startHttpServer() {
	httpServer := internal.NewHttpServer()
	httpServer.Ready()
	server := httptest.NewServer(httpServer.Echo())

	t.baseUrl = server.URL
	t.client = server.Client()
}

func (s *TestEnvironment) BaseUrl() string {
	return s.baseUrl
}

func (s *TestEnvironment) WiremockContainerUrl() string {
	return s.wiremockContainerUrl
}

func (s *TestEnvironment) Client() *http.Client {
	return s.client
}

func (s *TestEnvironment) PgxPool() *pgxpool.Pool {
	return s.pgxPool
}

func (s *TestEnvironment) RedisClient() *redis.Client {
	return s.redisClient
}

func (s *TestEnvironment) RabbitmqConn() *amqp091.Connection {
	return s.rabbitmqConn
}
