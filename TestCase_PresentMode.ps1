
$ComPort = 3
# $Duration = "30s"
$Duration = "10m"

$GoSysLatPath = "$PSScriptRoot/GoSysLat.exe"
$LogFolder = "LogFolder"
Set-Location "$PSScriptRoot"

$path = "HKCU:\System\GameConfigStore"

Function Set-RegistryData {
    param(
        [string]$p,
        [string]$n,
        [int]$v
    )
    if (!(Test-Path -Path $p)) {
        New-item -Path $p -Force
    }
    Set-ItemProperty -Path $p -Name $n -Value $v -Force
}

for ($i = 0; $i -lt 2; $i++) {
    Start-Sleep -s 3
    if($i % 2 -eq 0) { # FSO / Hardware: Independent Flip
        Write-Host "# FSO / Hardware: Independent Flip"
        Set-RegistryData -p $path -n "GameDVR_FSEBehaviorMode" -v 0
        Set-RegistryData -p $path -n "GameDVR_FSEBehavior" -v 0
        Set-RegistryData -p $path -n "GameDVR_HonorUserFSEBehaviorMode" -v 0
        Set-RegistryData -p $path -n "GameDVR_DXGIHonorFSEWindowsCompatible" -v 0
        Start-Sleep -s 1
        Start-Process -FilePath $GoSysLatPath -ArgumentList "-d3d9", "-fullscreen", "-port $ComPort", "-time $Duration", "-name FSO" , "-logs `"$LogFolder`"" -Wait
    } else { # FSE / Hardware: Legacy Flip
        Write-Host "# FSE / Hardware: Legacy Flip"
        Set-RegistryData -p $path -n "GameDVR_FSEBehaviorMode" -v 2
        Set-RegistryData -p $path -n "GameDVR_FSEBehavior" -v 2
        Set-RegistryData -p $path -n "GameDVR_HonorUserFSEBehaviorMode" -v 1
        Set-RegistryData -p $path -n "GameDVR_DXGIHonorFSEWindowsCompatible" -v 1
        Start-Sleep -s 1
        Start-Process -FilePath $GoSysLatPath -ArgumentList "-d3d9", "-fullscreen", "-port $ComPort", "-time $Duration", "-name FSE" , "-logs `"$LogFolder`"" -Wait
    }
}




