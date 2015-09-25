namespace go gentest

enum MyEnum {
	FIRST = 1,
	SECOND = 2
}

const map<MyEnum, string> STRINGY = {
	MyEnum.FIRST: "1st",
	MyEnum.SECOND: "2nd",
}

const i32 Fst = MyEnum.FIRST;
