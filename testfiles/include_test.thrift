include "./a/shared.thrift"

struct S {
  1: shared.AStruct s
  2: string c = shared.EXAMPLE_CONSTANT
}
 
