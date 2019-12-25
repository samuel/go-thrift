// super-stripped down scribe binary

namespace go gentest

enum ResultCode
{
  OK,
  TRY_LATER
}

struct LogEntry
{
  1:  string category,
  2:  string message
}

exception FailedException
{
  1: string reason
}

service scribe {
  ResultCode Log(1: list<LogEntry> messages) throws (1: FailedException f)
  LogEntry Echo(1: LogEntry messages) throws (1: FailedException f)
}
