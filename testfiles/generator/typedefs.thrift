namespace go gentest

typedef binary Binary
typedef string String
typedef i32    Int32
typedef map<i32,string> Map1
typedef map<i32,Map1> Map2

struct St {
	1: Binary b,
	2: String s,
	3: Int32 i
}

service Svc {
  Binary ping(),
  Map2 call(
    1:string name,
    2:Map1 map1,
    3:Binary data)
}

