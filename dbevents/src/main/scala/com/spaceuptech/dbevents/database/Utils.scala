package com.spaceuptech.dbevents.database

import java.nio.ByteBuffer
import java.nio.charset.StandardCharsets
import java.util.{Calendar, HashMap, Properties}
import java.util.concurrent.{Callable, ExecutorService}
import java.util.function.Consumer

import akka.actor.typed.ActorRef
import com.mongodb.client.model.changestream.{ChangeStreamDocument, FullDocument, OperationType}
import com.mongodb.{MongoClient, MongoClientURI}
import com.spaceuptech.dbevents.database.Database.ChangeRecord
import com.spaceuptech.dbevents.{DatabaseSource, Global}
import io.debezium.engine.format.Json
import io.debezium.engine.DebeziumEngine
import org.bson.{BsonDocument, Document}
import org.json4s._
import org.json4s.jackson.JsonMethods._

object Utils {
  def startMongoWatcher(projectId: String, dbAlias: String, conn: String, dbName: String, executorService: ExecutorService, actor: ActorRef[Database.Command]): java.util.concurrent.Future[_] = {
    val connectionString = new MongoClientURI(conn)
    val mongoClient = new MongoClient(connectionString)
    val db = mongoClient.getDatabase(dbName)

    val offsetStore = new KubeOffsetBackingStore()
    offsetStore.setName(s"$projectId-$dbAlias")
    offsetStore.start()

    val key = StandardCharsets.UTF_8.encode("resume-token")
    val resumeToken = offsetStore.get(java.util.Arrays.asList(Array(key))).get().get(key)

    var t = Calendar.getInstance().getTime

    val consumer: Consumer[ChangeStreamDocument[Document]] = (doc) => {
      // Check if 60 minutes have elapsed since last timer
      val cal = Calendar.getInstance()
      cal.setTime(t)
      cal.add(Calendar.MINUTE, 60)

      if (Calendar.getInstance().getTime.after(cal.getTime)) {
        t = Calendar.getInstance().getTime

        val value = StandardCharsets.UTF_8.encode(doc.getResumeToken.toJson)
        offsetStore.set(java.util.Collections.singletonMap(key, value), null)
      }

      doc.getOperationType match {
        case OperationType.INSERT =>
          actor ! ChangeRecord(
            payload = ChangeRecordPayload(
              op = Some("c"),
              before = None,
              after = Some(mongoDocumentToMap(doc.getFullDocument)),
              source = getMongoSource(projectId, dbAlias, doc)
            ),
            project = projectId,
            dbAlias = dbAlias,
            dbType = "mongo"
          )

        case OperationType.UPDATE | OperationType.REPLACE =>
          actor ! ChangeRecord(
            payload = ChangeRecordPayload(
              op = Some("u"),
              before = Option(mongoDocumentKeyToMap(doc.getDocumentKey)),
              after = Some(mongoDocumentToMap(doc.getFullDocument)),
              source = getMongoSource(projectId, dbAlias, doc)
            ),
            project = projectId,
            dbAlias = dbAlias,
            dbType = "mongo"
          )

        case OperationType.DELETE =>
          actor ! ChangeRecord(
            payload = ChangeRecordPayload(
              op = Some("d"),
              before = Option(mongoDocumentKeyToMap(doc.getDocumentKey)),
              after = None,
              source = getMongoSource(projectId, dbAlias, doc)
            ),
            project = projectId,
            dbAlias = dbAlias,
            dbType = "mongo"
          )

        case _ =>
          println(s"Invalid operation type (${doc.getOperationType.getValue}) received")
      }
    }

    executorService.submit(new Callable[Unit] {
      override def call(): Unit = {

        var w = db.watch().fullDocument(FullDocument.UPDATE_LOOKUP)
        if (resumeToken != null) {
          w = w.resumeAfter(mongoByteBufferToBsonDocument(resumeToken))
        }
        w.forEach(consumer)
      }
    })
  }

  def mongoByteBufferToBsonDocument(data: ByteBuffer): BsonDocument = {
    BsonDocument.parse(new String(data.array(), "UTF-8"))
  }

  def mongoDocumentKeyToMap(find: BsonDocument): Map[String, Any] =  {
    val id = find.getObjectId("_id").getValue.toHexString
    Map("_id" -> id)
  }

  def mongoDocumentToMap(doc: Document): Map[String, Any] =  {
    implicit val formats: DefaultFormats.type = org.json4s.DefaultFormats

    val jsonString = doc.toJson
    parse(jsonString).extract[Map[String, Any]]
  }

  def getMongoSource(projectId: String, dbAlias: String, doc: ChangeStreamDocument[Document]): ChangeRecordPayloadSource = {
    ChangeRecordPayloadSource(
      name = s"${projectId}_$dbAlias",
      ts_ms = doc.getClusterTime.getTime * 1000,
      table = doc.getNamespace.getCollectionName
    )
  }

  def startDebeziumEngine(source: DatabaseSource, executorService: ExecutorService, actor: ActorRef[Database.Command]): DebeziumStatus = {
    // Create the engine configuration object
    val props = source.dbType match {
      case "mysql" => prepareMySQLConfig(source)
      case "postgres" => preparePostgresConfig(source)
      case "sqlserver" => prepareSQLServerConfig(source)
      case _ => throw new IllegalArgumentException
    }

    // Create a new debezium engine

    val engine = DebeziumEngine.create(classOf[Json]).using(props).notifying(record => {
      implicit val formats: DefaultFormats.type = org.json4s.DefaultFormats
      // Extract the change feed value
      val jsonString = record.value()

      // Marshal the string only if the json string is not null
      if (jsonString != null) {
        // Parse the json value and forward it to our actor
        val payload = parse(jsonString).extract[ChangeRecordPayload]
        actor ! ChangeRecord(payload, source.project, source.dbAlias, source.dbType)
      }

    }).build()

    // Run the engine asynchronously
    val future = executorService.submit(engine)
    DebeziumStatus("", future, engine)
  }

  def prepareMySQLConfig(source: DatabaseSource): Properties = {
    val props = io.debezium.config.Configuration.empty().asProperties()
    props.setProperty("name", generateConnectorName(source))
    props.setProperty("connector.class", "io.debezium.connector.mysql.MySqlConnector")
    props.setProperty("offset.storage", getOffsetStorageClass)
    props.setProperty("offset.storage.file.filename", "./offsets.dat")
    props.setProperty("offset.flush.interval.ms", "60000")
    props.setProperty("converter.schemas.enable", "false")
    /* begin connector properties */
    props.setProperty("database.hostname", source.config.getOrElse("host", "localhost"))
    props.setProperty("database.port", source.config.getOrElse("port", "3306"))
    props.setProperty("database.user", source.config.getOrElse("user", "root"))
    props.setProperty("database.password", source.config.getOrElse("password", "my-secret-pw"))
    props.setProperty("database.include.list", source.config.getOrElse("db", "test"))
    props.setProperty("database.server.name", s"${generateConnectorName(source)}_connector")
    props.setProperty("database.ssl.mode", source.config.getOrElse("sslMode", "disabled"))
    props.setProperty("database.history", getDatabaseHistoryStorageClass)
    props.setProperty("database.history.file.filename", "./dbhistory.dat")

    props
  }

  def prepareSQLServerConfig(source: DatabaseSource): Properties = {
    val name = generateConnectorName(source)

    val props = io.debezium.config.Configuration.empty().asProperties()
    props.setProperty("name", name)
    props.setProperty("connector.class", "io.debezium.connector.postgresql.PostgresConnector")
    props.setProperty("offset.storage", getOffsetStorageClass)
    props.setProperty("offset.storage.file.filename", "./offsets.dat")
    props.setProperty("offset.flush.interval.ms", "60000")
    props.setProperty("converter.schemas.enable", "false")
    /* begin connector properties */
    props.setProperty("database.hostname", source.config.getOrElse("host", "localhost"))
    props.setProperty("database.port", source.config.getOrElse("port", "1433"))
    props.setProperty("database.user", source.config.getOrElse("user", "sa"))
    props.setProperty("database.password", source.config.getOrElse("password", "mypassword"))
    props.setProperty("database.dbname", source.config.getOrElse("db", "test"))
    props.setProperty("database.server.name", s"${generateConnectorName(source)}_connector")

    props
  }



  def preparePostgresConfig(source: DatabaseSource): Properties = {
    val name = generateConnectorName(source)

    val props = io.debezium.config.Configuration.empty().asProperties()
    props.setProperty("name", name)
    props.setProperty("connector.class", "io.debezium.connector.postgresql.PostgresConnector")
    props.setProperty("offset.storage", getOffsetStorageClass)
    props.setProperty("offset.storage.file.filename", "./offsets.dat")
    props.setProperty("offset.flush.interval.ms", "60000")
    props.setProperty("converter.schemas.enable", "false")
    /* begin connector properties */
    props.setProperty("plugin.name", "pgoutput")
    props.setProperty("slot.name", name)
    props.setProperty("publication.name", name)
    props.setProperty("database.hostname", source.config.getOrElse("host", "localhost"))
    props.setProperty("database.port", source.config.getOrElse("port", "5432"))
    props.setProperty("database.user", source.config.getOrElse("user", "postgres"))
    props.setProperty("database.password", source.config.getOrElse("password", "mysecretpassword"))
    props.setProperty("database.dbname", source.config.getOrElse("db", "postgres"))
    props.setProperty("schema.include.list", source.config.getOrElse("schema", "test"))
    props.setProperty("database.server.name", s"${generateConnectorName(source)}_connector")
    props.setProperty("database.sslmode", source.config.getOrElse("sslMode", "disable"))

    props
  }

  def getOffsetStorageClass: String = {
    Global.storageType match {
      case "k8s" => "com.spaceuptech.dbevents.database.KubeOffsetBackingStore"
      case _ => "org.apache.kafka.connect.storage.FileOffsetBackingStore"
    }
  }

  def getDatabaseHistoryStorageClass: String = {
    Global.storageType match {
      case "k8s" => "com.spaceuptech.dbevents.database.KubeDatabaseHistory"
      case _ => "io.debezium.relational.history.FileDatabaseHistory"
    }
  }

  def generateConnectorName(source: DatabaseSource): String = {
    s"${source.project}_${source.dbAlias}"
  }
}
