https://blog.friendsofgo.tech/posts/introduccion-a-grpc/
Protocol BUFFERS:

- Es una forma de diseñar una API y la otra es REST
- Es el protocolo que utiliza el framework gRPC sobre RPC y HTTP/2
- Es un mecanismo de serializacion de datos mas simple y mas liviana qque XML
- Utiliza buffer de protocolo y HTTP2 para la transmision de datos
- Necesitaremos una sintaxis y un compilador
- No es formato JSON ni tampoco XML, sino que es formato Protobuf
- gRPC permite a los desarrolladores crear API de alto rendimiento para arquitecturas de microservicios en centros de datos distribuidos. Es más adecuado para sistemas internos que requieren transmisión en tiempo real y grandes cargas de datos.
  > sintaxis para definir la estructura de nuestro datos
  > compilador para hacer la serializacion y desarializacion de los datos

> Conexion a la DB

```
func ConnectionToDatabase(ctx context.Context , dbConnectionString string)(*pgxpool,error){
  var dbPool *pgxpool.Pool
  var err

  var retryCount := 0
  for retryCount < 5 {
    dbPool , err := pgxpool.Connect(ctx,dbConnectionString )
    if(err == nil){
      break;
    }

    time.Sleep(5 * time.second)
    retryCount+++
  }

  if err != nil{
    log.Print("Error al conectarse a la DB.")
  }

  log.Print("Conectado a la DB")
  return dbPool, nil
}
```
