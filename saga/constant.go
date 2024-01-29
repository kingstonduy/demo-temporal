package shared

import "time"

const TASKQUEUE = "money-transfer-service-task-queue"
const WORKFLOW = "money-transfer-service-workflow"

const MONEY_TRANSFER_SERVICE_HOST_PORT = "localhost:7201"
const LIMITATION_SERVICE_HOST_PORT = "localhost:7202"
const T24_SERVICE_HOST_PORT = "localhost:7203"
const NAPAS_SERVICE_HOST_PORT = "localhost:7204"

const POSTGRES_URL = "jdbc:postgresql://localhost:5432/postgres"
const POSTGRES_HOST = "localhost"
const POSTGRES_PORT = 5432
const POSTGRES_USER = "postgres"
const POSTGRES_PASSWORD = "changeme"
const POSTGRES_DBNAME = "postgres"

<<<<<<< HEAD
const SERVICE_TIMEOUT = time.Second * 2
const CLIENT_TIMEOUT = time.Second * 10

const RETRYABLE_ERROR = "1"
const NONRETRYABLE_ERROR = "0"
=======
const TIMEOUT = time.Second * 2
>>>>>>> parent of 5bf56e9 (refactor service sleep)
