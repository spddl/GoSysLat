Param(
    [int]$Id
)

# first time call: .\TestCase.ps1 -Id -1

$ComPort = 3
$Duration = "60s"
$RestartNeeded = $true

$GoSysLatPath = "$PSScriptRoot/GoSysLat.exe"
$LogFolder = "LogFolder"
# $BatchFolder = "$PSScriptRoot/TestCases/"
Start-Transcript -Path "$PSScriptRoot\TestCase.log" -Append | Out-Null
Set-Location "$PSScriptRoot"

$TestCases = @()
for ($i = 0; $i -lt 5; $i++) { # 5 times
    $TestCases += @{ Value = "VRROptimizeEnable=0"; Name = 'VRROptimizeEnable Off' }
    $TestCases += @{ Value = "VRROptimizeEnable=1"; Name = 'VRROptimizeEnable On' }
}

if ($id -ne -1) {
	Write-Host "$id/$($TestCases.Count) => $($TestCases[$id].Name)"
    Start-Sleep -Seconds 45
    # Write-Host -NoNewLine 'Press any key to continue...';
    # $null = $Host.UI.RawUI.ReadKey('NoEcho,IncludeKeyDown');

    # Test
    Start-Process -FilePath $GoSysLatPath -ArgumentList "-d3d9", "-fullscreen", "-port $ComPort", "-time $Duration", "-name `"$($TestCases[$id].Name)`"" , "-logs `"$LogFolder`"" -Wait
}

$Id += 1 # check if there is a test case next to it

if ($null -ne $TestCases[$id]) {
    # new test environment is being prepared
    Set-ItemProperty -Path "HKCU:\Software\Microsoft\DirectX\UserGpuPreferences" -Name "DirectXUserGlobalSettings" -Value $TestCases[$id].Value

    if ($RestartNeeded) {
        # Autostart
        New-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\RunOnce" -Name "GoSysLat" -Value "Powershell -File $PSScriptRoot\$($MyInvocation.MyCommand.Name) -Id $Id"

        Start-Process -FilePath shutdown -ArgumentList "/r", "/t 0" -Wait
    } else {
        . $MyInvocation.MyCommand.Path -Id $Id
    }
}

Stop-Transcript | Out-Null
Start-Process -FilePath shutdown -ArgumentList "/s", "/t 0" -Wait