Param (
    [Parameter(Mandatory=$true,Position=0)]
    [ValidateNotNullOrEmpty()]
    [ValidateSet('clean','lint','test','build')]
    [string]$Command
)

$global:SysoFiles = @(
    "rsrc_windows_386.syso",
    "rsrc_windows_amd64.syso"
)

#######
# Utils
#######

function GetDuration {

    Param (
        [Parameter(Mandatory=$true,Position=0)]
        [System.TimeSpan]$TimeSpan
    )

    $Data = New-Object System.Collections.Generic.List[string]

    if ( $TimeSpan.Days -gt 0 ) {
        $DayString = $TimeSpan.Days.ToString();
        $Data.Add("${DayString}d");
    }

    if ( $TimeSpan.Hours -gt 0 ) {
        $HourString = $TimeSpan.Hours.ToString('00');
        $Data.Add("${HourString)}h");
    }

    if ( $TimeSpan.Minutes -gt 0 ) {
        $MinString = $TimeSpan.Minutes.ToString('00');
        $Data.Add("${MinString}m");
    }

    if ( $TimeSpan.Seconds -gt 0 ) {
        $SecString = $TimeSpan.Seconds.ToString('00');
        $Data.Add("${SecString}s");
    }

    $MillString = $TimeSpan.Milliseconds.ToString('000');
    $Data.Add("${MillString}ms");

    return "[" + (($Data) -join " ") + "]";

}

function Check-If-Installed {

    Param (
        [Parameter(Mandatory=$true,Position=0)]
        [string]$Name,
        [Parameter(Mandatory=$true,Position=1)]
        [string]$Command
    )

    if ( -not [bool] (Get-Command -ErrorAction Ignore -Type Application $Command) ) {
        Write-Host "Error: $Command not found in `$PATH. $Name must be installed.";
        exit 1;
    }

}

function Get-Product-Version {

    $Version =
        Select-String -LiteralPath "CHANGELOG.md" -Pattern "## \[([\d\.]+)\] - \d\d\d\d-\d\d-\d\d"
        | Select-Object -Index 0
        | %{$_.Matches.Groups[1].Value};

    if ( $Version -eq "" ) {
        Write-Host "Can not get version from CHANGELOG.md file...";
        exit 1;
    }

    return $Version;

}

function Get-File-Version {

    Param (
        [Parameter(Mandatory=$true,Position=0)]
        [string]$Version
    )

    $Value = "0"

    $LastTag = git describe --tags --abbrev=0;

    $Count = Invoke-Expression "git log $LastTag..HEAD --oneline"  | Measure-Object -Line | %{$_.Lines};

    if ( $Count -ne "" ) {
        $Value = $Count
    }

    return "${Version}.${Value}"

}

function Exec {

    Param (
        [Parameter(Mandatory=$true,Position=0)]
        [string]$Command
    )

    $TaskName = (Get-PSCallStack)[1].Command;
    Write-Host " > ${TaskName}: $Command";

    $global:LASTEXITCODE = 0
    Invoke-Expression "${Command}";

    if ( $LASTEXITCODE -ne 0 ) {
        throw "Command exited with code $LASTEXITCODE";
    }

}

#######
# Tasks
#######

function Remove-Winres-Files {

    foreach ($File in $global:SysoFiles) {
        Exec "Remove-Item -Path $File -Force -ErrorAction SilentlyContinue";
    }

}

function Clean {

    Exec "Remove-Item -Path 'output' -Recurse -Force -ErrorAction SilentlyContinue";
    Remove-Winres-Files;

}

function Lint {

    Check-If-Installed "Staticcheck" "staticcheck";

    Exec "go vet";
    Exec "staticcheck .";

}

function Test {

    Exec "go test ./...";

}

function Build {

    Check-If-Installed "go-winres" "go-winres";
    Check-If-Installed "Inno Setup" "iscc";

    $ProductVersion = Get-Product-Version;
    $FileVersion = Get-File-Version $ProductVersion;

    try {
        Exec "go-winres make --in ./winres.json --product-version $ProductVersion --file-version $FileVersion";
        Exec "go build -o ./output/ezstore.exe";
        Exec "iscc /Q 'setup.iss' /DPV='$ProductVersion' /DFV='$FileVersion'";
    }
    finally {
        $Code = $lastexitcode;
        Remove-Winres-Files;
        exit $Code;
    }

}

##################
# Script beginning
##################

Check-If-Installed "Golang" "go";

Write-Host "Starting..."

$sw = [System.Diagnostics.Stopwatch]::New();
$sw.Start();

switch ( $Command ) {
    'clean' {
        Clean;
        break;
    }
    'lint' {
        Lint;
        break;
    }
    'test' {
        Test;
        break;
    }
    'build' {
        Build;
        break;
    }
}

$sw.Stop();
$duration = GetDuration $sw.Elapsed;

Write-Host "Finished $duration";
