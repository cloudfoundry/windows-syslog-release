%%{
  machine syslog_rfc5424;
}%%

class Syslog

	def initialize
  	%% write data;
  end

private def bytes_ref(data)
  return nil if data == '-'
  data
end

def parse(data)
  log, nanosecond,  = {}, 0

  # set defaults for state machine parsing
  cs, p, pe, eof = 0, 0, data.length, data.length

  # use to keep track start of value
  mark = 0

  # taken directly from https://tools.ietf.org/html/rfc5424#page-8
  %%{
    action mark      { mark = p }
    action version   { log[:version]  = data[mark..p].to_i }
    action priority  { log[:priority] = data[mark..p].to_i }
    action hostname  { log[:hostname] = bytes_ref(data[mark..p]) }
    action appname   { log[:appname]  = bytes_ref(data[mark..p]) }
    action procid    { log[:proc_id]  = bytes_ref(data[mark..p]) }
    action msgid     { log[:msg_id]   = bytes_ref(data[mark..p]) }
    action sdid      {
      log[:data] = {
        id: data[mark..p],
        properties: {}
      }
    }
    action paramname  { paramName = data[mark..p] }
    action paramvalue { log[:data][:properties][paramName] = data[mark..p] }

    action timestamp {
      log[:timestamp] = Time.utc(
        data[mark..mark+4].to_i,
        data[mark+5..mark+7].to_i,
        data[mark+8..mark+10].to_i,
        data[mark+11..mark+13].to_i,
        data[mark+14..mark+16].to_i,
        data[mark+17..p-1].to_f
      )
    }
    action message { log[:message] = data[mark..p] }

    include syslog_rfc5424 "../syslog.rl";
    write init;
    write exec;
  }%%

  if cs < syslog_rfc5424_first_final
    raise "error in msg at pos #{p}: #{data}"
  end

  return log
end
end
