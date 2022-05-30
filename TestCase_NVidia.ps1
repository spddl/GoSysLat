Param(
    [int]$Id
)

# first time call: .\TestCase.ps1 -Id -1
# 511.79

$ComPort = 3
$Duration = "45s"
$RestartNeeded = $false

$GoSysLatPath = "$PSScriptRoot/GoSysLat.exe"
$LogFolder = "LogFolder"
# $BatchFolder = "$PSScriptRoot/TestCases/"
# $ApplyPath = Join-Path -Path $BatchFolder -ChildPath "0_Apply.exe"
Start-Transcript -Path "$PSScriptRoot\TestCase.log" -Append | Out-Null
Set-Location "$PSScriptRoot"

$width = 1280
$height = 720

# $One = @('1_Aspect_ratio.exe', '1_Fullscreen.exe', '1_Integer_scaling.exe', '1_no_scaling.exe')
$One = @('1_Aspect_ratio.exe', '1_Fullscreen.exe', '1_no_scaling.exe')
$Two = @('2_Perform_scaling_on_Display.exe', '2_Perform_scaling_on_GPU.exe')
$Three = @('3_Override_the_scaling_mode_set_by_games_and_programs_OFF.exe', '3_Override_the_scaling_mode_set_by_games_and_programs_ON.exe')

$TestCases = @()

for ($i = 0; $i -lt 3; $i++) {
    # 5 times
    foreach ($1 in $One) {
        foreach ($2 in $Two) {
            foreach ($3 in $Three) {
                $TestCases += @{
                    Name = "$($1.substring(2, $1.Length-6)), $($2.substring(2, $2.Length-6)), $($3.substring(2, $3.Length-6))"
                    Test = @($1, $2, $3)
                }
            }
        }
    }
}

if ($id -ne -1) {
    # Write-Host "$id/$($TestCases.Count) => $($TestCases[$id].Name)"
	if ($TestCases[$id].Name -eq "no_scaling, Perform_scaling_on_Display, Override_the_scaling_mode_set_by_games_and_programs_OFF" -or
		$TestCases[$id].Name -eq "Aspect_ratio, Perform_scaling_on_Display, Override_the_scaling_mode_set_by_games_and_programs_OFF"
	) {
		Write-Host 'reposition the sensor';
		$null = $Host.UI.RawUI.ReadKey('NoEcho,IncludeKeyDown');
	}
	# Start-Sleep -Seconds 3

    # .\GoSysLat.exe -d3d9 -fullscreen -port 4 -width 800 -height 600

	#Write-Host 'press any key to start';
	#$null = $Host.UI.RawUI.ReadKey('NoEcho,IncludeKeyDown');

    # Test
    Start-Process -FilePath $GoSysLatPath -ArgumentList "-d3d9", "-fullscreen", "-port $ComPort", "-time $Duration", "-width $width", "-height $height", "-name `"$($TestCases[$id].Name)`"" , "-logs `"$LogFolder`"" -Wait
} else {
	#Start-Process -FilePath "C:\Program Files\WindowsApps\NVIDIACorp.NVIDIAControlPanel_8.1.962.0_x64__56jybvy8sckqj\nvcplui.exe"
}

$Id += 1 # check if there is a test case next to it

if ($null -ne $TestCases[$id]) {
    # new test environment is being prepared
	Write-Host "$id/$($TestCases.Count) => $($TestCases[$id].Name)"

    $nvcplui = Start-Process -FilePath "C:\Program Files\WindowsApps\NVIDIACorp.NVIDIAControlPanel_8.1.962.0_x64__56jybvy8sckqj\nvcplui.exe" -passthru
	$nvcplui.WaitForExit()
	
	<#
    Start-Sleep -Seconds 6
    foreach ($testFiles in $TestCases[$id].Test) {
        $testFile = Join-Path -Path $BatchFolder -ChildPath $testFiles
        if (Test-Path -Path $testFile -PathType Leaf) {
			Write-Host $testFile
            Start-Process -FilePath $testFile -Wait -NoNewWindow
			Start-Sleep -Seconds 2
        }
    }
	Start-Sleep -Seconds 2
	Write-Host $ApplyPath
    Start-Process -FilePath $ApplyPath -Wait -NoNewWindow
	Start-Sleep -Seconds 2
    #>
	#Stop-Process -Name "nvcplui"

    if ($RestartNeeded) {
        New-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\RunOnce" -Name "GoSysLat" -Value "Powershell -File $PSScriptRoot\$($MyInvocation.MyCommand.Name) -Id $Id"

        Start-Process -FilePath shutdown -ArgumentList "/r", "/t 0" -Wait
    }
    else {
        . $MyInvocation.MyCommand.Path -Id $Id
    }
}

Stop-Transcript | Out-Null
# Start-Process -FilePath shutdown -ArgumentList "/s", "/t 0" -Wait