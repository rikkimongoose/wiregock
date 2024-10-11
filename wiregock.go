package wiregock

const (
    productName = "WireGock"
    productVersion = "0.1.0"
)

const anyMathods = [...]string{ "GET", "HEAD", "OPTIONS", "TRACE", "PUT", "DELETE", "POST", "PATCH", "CONNECT" }

import (
    "fmt"
    "flag"
    "github.com/ilyakaznacheev/cleanenv"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "github.com/google/uuid"
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/adaptor"
    "context"
    "crypto/tls"
    "crypto/x509"
    "os"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    actuator "github.com/sinhashubham95/go-actuator"
    "strings"
)

type AppConfig struct {
    Server struct {
        Host string `yaml:"host", envconfig:"SERVER_HOST", env-default:"localhost", env-description:"server host"`
        Port string `yaml:"port", envconfig:"SERVER_PORT", env-default:"8080", env-description:"server port"`
    } `yaml:"server"`
    Mongo struct {
        url string `yaml:"url", envconfig:"MONGO_URL", env-description:"MongoDB connection string"`
        database string `yaml:"db", envconfig:"MONGO_DB", env-description:"MongoDB database"`
        collection string `yaml:"collection", envconfig:"MONGO_COLLECTION", env-description:"MongoDB collection"`
        caFile string `yaml:"caFile", envconfig:"MONGO_CA", env-description:"path to CA certificate"`
        certFile string `yaml:"certFile", envconfig:"MONGO_CERT", env-description:"path to public client certificate"`
        keyFile string `yaml:"keyFile", envconfig:"MONGO_KEY", env-description:"path to private client key"`
    } `yaml:"mongo"`
    Log struct {
        Encoding string `yaml:"encoding", env-default:"json", envconfig:"LOG_ENCODING", env-description:"storage format for logs"`
        OutputPaths []string `yaml:"output", env-default:"stdout,/tmp/logs", envconfig:"LOG_OUTPUTPATH", env-description:"output pipelines for logs"`
        ErrorOutputPaths []string `yaml:"erroutput", env-default:"stderr", envconfig:"LOG_OUTPUTERRORPATH", env-description:"error pipelines for logs"`
    } `yaml:"log"`
}

var logger *zap.Logger
var config *AppConfig

func main() {
    var err error
    cfgPath := flag.String("CONFIG", "config.yml", "Path to application config file")

    err := cleanenv.ReadConfig(cfgPath, &config)
    if err != nil {
        panic(fmt.Sprintf("Unable to load config file %s", cfgPath))
    }

    zc := &zap.Config{
        Encoding:         config.Log.Encoding,
        OutputPaths:      config.Log.OutputPaths,
        ErrorOutputPaths: config.Log.ErrorOutputPaths
    }
    logger, err := zc.Build(zap.AddStacktrace(zapcore.PanicLevel)) // Или NewDevelopment для более подробного логирования
    if err != nil {
        panic(err) // Не удалось создать логгер
    }
    defer logger.Sync() // все асинхронные логи будут записаны перед выходом

    opts := options.Client()
                .ApplyURI(config.MongoDatabase.url)
                .SetTLSConfig(mongoTlsConfig(
                        &MongoTlsConfigInput(
                            caFile: config.Mongo.caFile,
                            certFile: config.Mongo.certFile,
                            keyFile: config.Mongo.keyFile
                        )
                    )
                )
    ctx := context.TODO()
    client, err := mongo.Connect(ctx, opts)
    var resultPing bson.M
    if err := client.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Decode(&resultPing); err != nil {
        logger.Error(err)
    }
    logger.Info("Pinged your deployment. You successfully connected to MongoDB!")

    wiremockCollection := client.Database(config.MongoDatabase.database).Collection(config.MongoDatabase.collection)
    cursor, err := wiremockCollection.Find(ctx, bson.M{})
    if err != nil {
        logger.Fatal(err)
    }
    defer cur.Close(ctx)

    server := fiber.New()
    actuatorConfig := &actuator.Config{
        Endpoints: []int{
            actuator.Env,
            actuator.Info,
            actuator.Metrics,
            actuator.Ping,
            actuator.Shutdown,
            actuator.ThreadDump,
        },
        Name: productName,
        Port: config.Server.Port,
        Version: productVersion,
    }

    server.Get("/actuator", adaptor.HTTPHandlerFunc(actuator.GetActuatorHandler(actuatorConfig)))

    for cursor.Next(ctx) {
        var wiremock Mock
        if err = cursor.Decode(&wiremock); err != nil {
            logger.Fatal(err)
        }

        methodNames := loadMethods(wiremock["request"]["method"])

        server.Add(methodNames, wiremock["request"]["urlPath"], func(c fiber.Ctx) error {
            condition = parseCondition(wiremock["request"])
            if !condition.check(c) {
                return fiber.StatusNotFound
            }

            for key, value := range wiremock["response"]["headers"] {
                c.Set(key, value)
            }
            return c.Status(wiremock["response"]["code"]).SendString(wiremock["response"]["body"])
        })
    }
    serverPath = fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
    logger.Info("Establishing server at URL", zap.String("url", serverPath))
    server.Listen(serverPath)
}

func loadMethods(methodNames string) string[] {
    if strings.Compare(methodNames, "ANY") {
        return anyMathods
    }
    return strings.Split(methodNames, ",")
}


type MongoTlsConfigInput struct {
    caFile, certFile, keyFile string
}

func mongoTlsConfig(input *MongoTlsConfigInput) *tls.Config {
    caCertPool := x509.NewCertPool()
    if input.caFile == nil || input.certFile == nil || input.keyFile == nil {
        return &tls.Config {
            RootCAs: caCertPool,
            ClientAuth: tls.NoClientCert,
            ClientCAs: nil,
            InsecureSkipVerify: true,
            Certificates: []tls.Certificate{},
        }
    }
    caFile := input.caFile
    certFile := input.certFile
    keyFile := input.keyFile

    // Loads CA certificate file
    caCert, err := os.ReadFile(caFile)
    if err != nil {
        logger.Error(err)
    }
    if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
        logger.Error("Error: CA file must be in PEM format")
    }
    // Loads client certificate files
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        logger.Error(err)
    }
    // Instantiates a Config instance
    return &tls.Config{
        RootCAs:      caCertPool,
        Certificates: []tls.Certificate{cert},
    }
}