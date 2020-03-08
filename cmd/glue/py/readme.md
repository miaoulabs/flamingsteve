# Glue Python Generated Client

```python
import fpixel_pb2
import fpixel_pb2_grpc

def run():
  connection = grpc.insecure_channel('localhost:8080')
  stub = fpixel_pb2_grpc.FlamePixelsStub(connection)

  response = stub.ListDisplays(fpixel_pb2.EmptyRequest())
  print("Displays: " + response.message)

  response = stub.ListDisplays(fpixel_pb2.EmptyRequest())
  print("Sensors: " + response.message)
```

- Reading: [GRPC Python Quick Start](https://grpc.io/docs/quickstart/python/)
