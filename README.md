# wiregock
Small and very fast and stable implementation of Wiremock with Goland and MongoDB.



## Configuration

| config file  | env     | default  |  description |
|---|---|---|---|
|              | CONFIG | config.yml | path to configuration file (supports YAML, JSON and TOML) |
| server.host | SERVER_HOST | localhost | server host  |
| server.host | SERVER_POST | 8080   server port |
| mongo.url | MONGO_URL | mongodb://localhost:27017  | MongoDB connection string |
| mongo.db | MONGO_DB | local | MongoDB database name |
| mongo.collection | MONGO_COLLECTION | mocks  | MongoDB collection |
| mongo.caFile | MONGO_CA |   | path to CA certificate |
| mongo.certFile | MONGO_CERT |   | path to public client certificate |
| mongo.keyFile | MONGO_KEY |   | path to private client key |
| log.encoding | LOG_ENCODING | json  | storage format for logs |
| log.output | LOG_OUTPUTPATH | stdout,/tmp/logs  | output pipelines for logs |
| log.erroutput | LOG_OUTPUTERRORPATH | stderr  | error pipelines for logs |
