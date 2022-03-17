$date=Get-Date

while($true){
	$newDate=Get-Date
 
        Get-EventLog -LogName Security -After $date -Before $newDate | ForEach-Object {
		$_ | ConvertTo-Json -Compress | Out-File -Append -Encoding utf8 -FilePath /var/vcap/sys/log/event_logger/Security.log
	}
	
	Get-EventLog -LogName Application -After $date -Before $newDate | ForEach-Object {
		$_ | ConvertTo-Json -Compress | Out-File -Append -Encoding utf8 -FilePath /var/vcap/sys/log/event_logger/Application.log 
	}
	
	Get-EventLog -LogName System -After $date -Before $newDate | ForEach-Object {
		$_ | ConvertTo-Json -Compress | Out-File -Append -Encoding utf8 -FilePath /var/vcap/sys/log/event_logger/System.log 
	}
    
        $date=$newDate
	start-sleep -Seconds 1
}
 
