Param(
    [int]$Id
)

# first time call: .\TestCase.ps1 -Id -1

$ComPort = 3
$Duration = "15s"
$RestartNeeded = $true

$GoSysLatPath = "$PSScriptRoot/GoSysLat.exe"
$LogFolder = "LogFolder"
# $BatchFolder = "$PSScriptRoot/TestCases/"
Start-Transcript -Path "$PSScriptRoot\TestCase.log" -Append | Out-Null
Set-Location "$PSScriptRoot"

$TestCases = @()
for ($i = 0; $i -lt 5; $i++) {
    # 5 times
    $TestCases += @{ Value = 0x2A; Name = '0x2A, Short, Fixed, High foreground boost' } 
    $TestCases += @{ Value = 0x29; Name = '0x29, Short, Fixed, Medium foreground boost.' }
    $TestCases += @{ Value = 0x28; Name = '0x28, Short, Fixed, No foreground boost.' }
    
    $TestCases += @{ Value = 0x26; Name = '0x26, Short, Variable, High foreground boost.' }
    $TestCases += @{ Value = 0x25; Name = '0x25, Short, Variable, Medium foreground boost.' }
    $TestCases += @{ Value = 0x24; Name = '0x24, Short, Variable, No foreground boost.' }

    $TestCases += @{ Value = 0x1A; Name = '0x1A, Long, Fixed, High foreground boost.' }
    $TestCases += @{ Value = 0x19; Name = '0x19, Long, Fixed, Medium foreground boost.' }
    $TestCases += @{ Value = 0x18; Name = '0x18, Long, Fixed, No foreground boost.' }

    $TestCases += @{ Value = 0x16; Name = '0x16, Long, Variable, High foreground boost.' }
    $TestCases += @{ Value = 0x15; Name = '0x15, Long, Variable, Medium foreground boost.' }
    $TestCases += @{ Value = 0x14; Name = '0x14, Long, Variable, No foreground boost.' }
}

if ($id -ne -1) {
    Write-Host "$id/$($TestCases.Count) => $($TestCases[$id].Name)"
    Start-Sleep -Seconds 35

    # Write-Host -NoNewLine 'Press any key to continue...'
    # $null = $Host.UI.RawUI.ReadKey('NoEcho,IncludeKeyDown')

    # Test
    Start-Process -FilePath $GoSysLatPath -ArgumentList "-d3d9", "-fullscreen", "-port $ComPort", "-time $Duration", "-name `"$($TestCases[$id].Name)`"" , "-logs `"$LogFolder`"" -Wait
}

$Id += 1 # check if there is a test case next to it

if ($null -ne $TestCases[$id]) {
    # new test environment is being prepared
    Set-ItemProperty -Path "HKLM:\SYSTEM\ControlSet001\Control\PriorityControl" -Name "Win32PrioritySeparation" -Value $TestCases[$id].Value -Type DWord

    if ($RestartNeeded) {
        New-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\RunOnce" -Name "GoSysLat" -Value "Powershell -File $PSScriptRoot\$($MyInvocation.MyCommand.Name) -Id $Id"

        Start-Process -FilePath shutdown -ArgumentList "/r", "/t 0" -Wait # restart
    }
    else {
        . $MyInvocation.MyCommand.Path -Id $Id
    }
}

Stop-Transcript | Out-Null
Start-Process -FilePath shutdown -ArgumentList "/s", "/t 0" -Wait # shutdown