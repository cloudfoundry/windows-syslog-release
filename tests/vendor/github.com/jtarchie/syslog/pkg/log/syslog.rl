%%{
  machine syslog_rfc5424;

  nil           = "-";
  nonzero_digit = "1".."9";
  printusascii  = "!".."~";
  sp            = " ";

  utf8_string = any*;
  bom         = 0xEF 0xBB 0xBF;
  msg_utf8    = bom utf8_string;
  msg_any     = any*;
  msg         = (msg_any | msg_utf8) >mark %message;

  sd_name         = printusascii{1,32} -- ("=" | sp | "]" | 34 | '"');
  param_value     = utf8_string >mark %paramvalue;
  param_name      = sd_name >mark %paramname;
  sd_id           = sd_name >mark %sdid;
  sd_param        = param_name '="' param_value :>> '"';
  sd_element      = "[" sd_id ( sp sd_param )* "]";
  structured_data = nil | sd_element{1};

  time_hour      = digit{2};
  time_minute    = digit{2};
  time_second    = digit{2};
  time_secfrac   = "." digit{1,6};
  time_numoffset = ("+" | "-") time_hour ":" time_minute;
  time_offset    = "Z" | time_numoffset;
  partial_time   = time_hour ":" time_minute ":" time_second time_secfrac?;
  full_time      = partial_time time_offset;
  date_mday      = digit{2};
  date_month     = digit{2};
  date_fullyear  = digit{4};
  full_date      = date_fullyear "-" date_month "-" date_mday;
  timestamp      = nil | (full_date "T" full_time) >mark %timestamp;

  msg_id   = nil | printusascii{1,32} >mark %msgid;
  proc_id  = nil | printusascii{1,128} >mark %procid;
  app_name = nil | printusascii{1,48} >mark %appname;

  hostname = nil | printusascii{1,255} >mark %hostname;
  version  = (nonzero_digit digit{0,2}) >mark %version;
  prival   = digit{1,3} >mark %priority;
  pri      = "<" prival ">";
  header   = pri version sp timestamp sp hostname sp app_name sp proc_id sp msg_id;

  syslog_msg = header sp structured_data (sp msg)?;
  main := syslog_msg;

}%%
