from __future__ import print_function
import logging
import grpc

import fpixel_pb2
import fpixel_pb2_grpc


def run():
  connection = grpc.insecure_channel('localhost:8080')
  stub = fpixel_pb2_grpc.FlamePixelsStub(connection)

  response = stub.ListDisplays(fpixel_pb2.EmptyRequest())
  print("=== Displays === \n" + '{}'.format(response))

  response = stub.ListSensors(fpixel_pb2.EmptyRequest())
  print("=== Sensors === \n" + '{}'.format(response))

if __name__ == '__main__':
  logging.basicConfig()
  run()
