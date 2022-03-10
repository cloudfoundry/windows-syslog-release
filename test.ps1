$date=Get-Date

while($true){
    start-sleep -Seconds 1
    Get-EventLog -LogName Security -After $date | ConvertTo-Json -Compress >> /var/vcap/sys/log/test/Security.log
    Get-EventLog -LogName System -After $date | ConvertTo-Json -Compress >> /var/vcap/sys/log/test/System.log
    Get-EventLog -LogName Application -After $date | ConvertTo-Json -Compress >> /var/vcap/sys/log/test/Application.log
    
    $date=Get-Date
}
