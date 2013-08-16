namespace go scribe

enum ResultCode {
	OK,
	TRY_LATER
}

struct LogEntry {
	1: required string category;
	2: required string message;
}

struct ScribeLogResponse {
	required ResultCode value;
}

struct ScribeLogRequest {
	required list<LogEntry> messages;
}

service Scribe {
	ResultCode Log(list<LogEntry> messages)
}
