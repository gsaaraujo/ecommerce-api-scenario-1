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

func (t *TestEnvironment) Start() error {
	err := t.startContainers()
	if err != nil {
		return err
	}

	_ = os.Setenv("AWS_REGION", "us-east-1")
	_ = os.Setenv("AWS_ACCESS_KEY_ID", "test")
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	_ = os.Setenv("AWS_ENDPOINT_URL", t.localstackContainerUrl)
	_ = os.Setenv("AWS_SECRET_MANAGER_NAME", "secret-us-east-1-local-app")
	_ = os.Setenv("TERN_MIGRATIONS_PATH", "../migrations")
	_ = os.Setenv("ZIPCODE_URL", t.wiremockContainerUrl)

	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	t.awsConfig = awsConfig

	err = t.createSecrets()
	if err != nil {
		return err
	}

	err = t.runMigrations()
	if err != nil {
		return err
	}

	pgxPool, err := pgxpool.New(context.Background(), t.postgresContainerUrl)
	if err != nil {
		return err
	}

	t.pgxPool = pgxPool

	redisParsedUrl, err := redis.ParseURL(t.redisContainerUrl)
	if err != nil {
		return err
	}

	t.redisClient = redis.NewClient(&redis.Options{
		Addr:     redisParsedUrl.Addr,
		Password: redisParsedUrl.Password,
	})

	rabbitmqConn, err := amqp091.Dial(t.rabbitmqContainerUrl)
	if err != nil {
		return err
	}

	t.rabbitmqConn = rabbitmqConn

	t.startHttpServer()
	return nil
}

func (t *TestEnvironment) startContainers() error {
	postgresContainer, err := NewPostgresContainer()
	if err != nil {
		return err
	}

	t.postgresContainerUrl = postgresContainer.Url()

	localstackContainer, err := NewLocalstackContainer()
	if err != nil {
		return err
	}

	t.localstackContainerUrl = localstackContainer.Url()

	wiremockContainer, err := NewWiremockContainer()
	if err != nil {
		return err
	}

	t.wiremockContainerUrl = wiremockContainer.Url()

	redisContainer, err := NewRedisContainer()
	if err != nil {
		return err
	}

	t.redisContainerUrl = redisContainer.Url()

	rabbitmqContainer, err := NewRabbitmqContainer()
	if err != nil {
		return err
	}

	t.rabbitmqContainerUrl = rabbitmqContainer.Url()
	return nil
}

func (t *TestEnvironment) createSecrets() error {
	secretsClient := secretsmanager.NewFromConfig(t.awsConfig)

	_, err := secretsClient.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
		Name: aws.String("secret-us-east-1-local-app"),
		SecretString: aws.String(fmt.Sprintf(`
			{
				"REDIS_URL": "%s",
				"POSTGRES_URL": "%s",
				"RABBITMQ_URL": "%s",
				"ZIPCODE_TOKEN": "a7416146283d464294cebea38d5cb5ff",
				"ACCESS_TOKEN_SIGNING_KEY": "81c4a8d5b2554de4ba736e93255ba633"
			}
		`, t.redisContainerUrl, t.postgresContainerUrl, t.rabbitmqContainerUrl)),
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *TestEnvironment) runMigrations() error {
	urlParsed, err := url.Parse(t.postgresContainerUrl)
	if err != nil {
		return err
	}

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
