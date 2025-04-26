Param (
    [Parameter(Mandatory=$true,Position=0)]
    [ValidateNotNullOrEmpty()]
    [ValidateSet('clean','lint','test','build')]
    [string]$Command
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
        $DayString = $TimeSpan.Days;
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

function Exec {

    Param (
        [Parameter(Mandatory=$true,Position=0)]
        [string]$Command
    )

    $TaskName = (Get-PSCallStack)[1].Command;
    Write-Host " > ${TaskName}: $Command";
    Invoke-Expression "$Command";

}

#######
# Tasks
#######

function Clean-Winres-Files {

    $SysoFiles = (
        "rsrc_windows_386.syso",
        "rsrc_windows_amd64.syso"
    )

    foreach ($File in $SysoFiles) {
        Exec "Remove-Item -Path $File -Force -ErrorAction SilentlyContinue";
    }

}

function Clean {

    Exec "Remove-Item -Path 'output' -Recurse -Force -ErrorAction SilentlyContinue";
    Clean-Winres-Files;

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

    begin {
        Check-If-Installed "Inno Setup" "iscc";
    }

    process {
        Exec "go-winres make --in ./winres.json";
        Exec "go build -o ./output/ezstore.exe";
        Exec "iscc /Q 'setup.iss'";
    }

    end {
        Clean-Winres-Files;
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
