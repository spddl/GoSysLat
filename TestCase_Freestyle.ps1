$ComPort = 3
$Duration = "60s"

$GoSysLatPath = "$PSScriptRoot/GoSysLat.exe"
$LogFolder = "LogFolder"
Set-Location "$PSScriptRoot"

$TestCaseName = Read-Host "TestCase Name?"
# Start-Process -FilePath $GoSysLatPath -ArgumentList "-d3d9", "-fullscreen", "-port $ComPort", "-time $Duration", "-name `"$TestCaseName`"" , "-logs `"$LogFolder`"" -Wait
Start-Process -FilePath $GoSysLatPath -ArgumentList "-ogl", "-fullscreen", "-port $ComPort", "-time $Duration", "-name `"$TestCaseName`"" , "-logs `"$LogFolder`"" -Wait
