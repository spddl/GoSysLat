Param(
    [int]$Id
)

# first time call: .\TestCase.ps1 -Id -1

$ComPort = 3
$Duration = "60s"
$RestartNeeded = $true

$GoSysLatPath = "$PSScriptRoot/GoSysLat.exe"
$LogFolder = "LogFolder"
$BatchFolder = "$PSScriptRoot/TestCases/"
Start-Transcript -Path "$PSScriptRoot\TestCase.log" -Append | Out-Null
Set-Location "$PSScriptRoot"

$TestCases = @()
for ($i = 0; $i -lt 15; $i++) { # 15 times
    # https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/bcdedit--set
    <# $TestCases += @{
        Name  = 'tscsyncpolicy Enhanced'
        Batch = 'tscsyncpolicy_enhanced.bat'
    }
    $TestCases += @{
        Name  = 'tscsyncpolicy Legacy'
        Batch = 'tscsyncpolicy_legacy.bat'
    }
    $TestCases += @{
        Name  = 'tscsyncpolicy default'
        Batch = 'tscsyncpolicy_default.bat'
    } #>
	$TestCases += @{
        Name  = 'disabledynamictick yes'
        Batch = 'disabledynamictick_yes.bat'
    }
    $TestCases += @{
        Name  = 'disabledynamictick no'
        Batch = 'disabledynamictick_no.bat'
    }
    <# $TestCases += @{
        Name  = 'useplatformtick yes'
        Batch = 'useplatformtick_yes.bat'
    }
    $TestCases += @{
        Name  = 'useplatformtick deletevalue'
        Batch = 'useplatformtick_deletevalue.bat'
    } #>
	<# $TestCases += @{
        Name  = 'Services Disable'
        Batch = 'Services_Disable.bat'
    }
    $TestCases += @{
        Name  = 'Services Enable'
        Batch = 'Services_Enable.bat'
    } #>
}

Write-Host "$id/$($TestCases.Count) => $($TestCases[$id].Name)"
Write-Host "Start-Process -FilePath $GoSysLatPath -ArgumentList `"-d3d9`", `"-fullscreen`", `"-port $ComPort`", `"-time $Duration`", `"-name $($TestCases[$id].Name)`" , `"-logs $LogFolder`" -Wait"

if ($id -ne -1) {
    Start-Sleep -Seconds 45
    # Write-Host -NoNewLine 'Press any key to continue...';
    # $null = $Host.UI.RawUI.ReadKey('NoEcho,IncludeKeyDown');

    # Test
    Start-Process -FilePath $GoSysLatPath -ArgumentList "-d3d9", "-fullscreen", "-port $ComPort", "-time $Duration", "-name `"$($TestCases[$id].Name)`"" , "-logs `"$LogFolder`"" -Wait
}

$Id += 1 # check if there is a test case next to it

if ($null -ne $TestCases[$id]) {
    # new test environment is being prepared
    $TestCase = Join-Path -Path $BatchFolder -ChildPath $TestCases[$id].Batch
    if (Test-Path -Path $TestCase -PathType Leaf) {
        Start-Process -FilePath $TestCase -Wait # Test environment is created
    }

    if ($RestartNeeded) {
        New-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\RunOnce" -Name "GoSysLat" -Value "Powershell -File $PSScriptRoot\$($MyInvocation.MyCommand.Name) -Id $Id"

        Start-Process -FilePath shutdown -ArgumentList "/r", "/t 0" -Wait
    } else {
        . $MyInvocation.MyCommand.Path -Id $Id
    }
}

Stop-Transcript | Out-Null
Start-Process -FilePath shutdown -ArgumentList "/s", "/t 0" -Wait