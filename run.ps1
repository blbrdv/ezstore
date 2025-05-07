Param (
    [Parameter(Mandatory=$true,Position=0)]
    [ValidateNotNullOrEmpty()]
    [ValidateSet('clean','format','lint','test','check','build','rebuild')]
    [string]$Command
)

$global:SysoFiles = @(
    "rsrc_windows_386.syso",
    "rsrc_windows_amd64.syso"
);

#######
# Utils
#######

function Get-Duration {

    Param (
        [Parameter(Mandatory=$true,Position=0)]
        [System.TimeSpan]$TimeSpan
    )

    $Data = New-Object System.Collections.Generic.List[string];

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

    $Value = "0";

    $LastTag = git describe --tags --abbrev=0;

    $Count = Invoke-Expression "git log $LastTag..HEAD --oneline"  | Measure-Object -Line | %{$_.Lines};

    if ( $Count -ne "" ) {
        $Value = $Count;
    }

    return "${Version}.${Value}";

}

function Exec {

    Param (
        [Parameter(Mandatory=$true,Position=0)]
        [string]$Command
    )

    $TaskName = (Get-PSCallStack)[1].Command;
    Write-Host " > ${TaskName}: $Command";

    $global:LASTEXITCODE = 0;
    try {
        Invoke-Expression "${Command}";
    }
    catch {
        Write-Host "Invalid command ${Command}"
        $global:LASTEXITCODE = 1;
    }

    if ( $global:LASTEXITCODE -ne 0 ) {
        throw "Command exited with code $global:LASTEXITCODE";
    }

}

#######
# Tasks
#######

function Remove-Winres-Files {

    foreach ($File in $global:SysoFiles) {
        Exec "Remove-Item -Path .\cmd\$File -Force -ErrorAction SilentlyContinue";
    }

}

function Clean {

    Exec "Remove-Item -Path 'output' -Recurse -Force -ErrorAction SilentlyContinue";
    Exec "Remove-Item -Path 'release' -Recurse -Force -ErrorAction SilentlyContinue";
    Remove-Winres-Files;

}

function Format {

    Write-Host " > Format";
    $Location = (Get-Location | %{$_.Path}) + "\";
    $Files = Get-Childitem –Path . -Include *.go -Recurse -ErrorAction SilentlyContinue | %{$_.FullName.Replace($Location,'')};

    $Result = & {gofmt -l $Files};

    if ( $Result -ne "" ) {
        foreach ($File in $Result.split(" ")) {
            Write-Host $File;
        }
        $global:LASTEXITCODE = 1;
        throw "Format exited with code $global:LASTEXITCODE";
    }

}

function Lint {

    Check-If-Installed "Staticcheck" "staticcheck";

    Exec "go vet .\...";
    Exec "staticcheck .\...";

}

function Test {

    Exec "go test .\internal\...";

}

function Build {

    Check-If-Installed "go-winres" "go-winres";
    Check-If-Installed "7-Zip" "7z";
    Check-If-Installed "Inno Setup" "iscc";

    $ProductVersion = Get-Product-Version;
    $FileVersion = Get-File-Version $ProductVersion;

    Exec "go-winres make --in .\winres.json --product-version $ProductVersion --file-version $FileVersion";

    foreach ($File in $global:SysoFiles) {
        Exec "Move-Item -Path $File -Destination .\cmd\$File -Force -ErrorAction SilentlyContinue";
    }

    Exec "go build -ldflags='-X main.version=$ProductVersion' -o .\output\ezstore.exe .\cmd";

    Exec "7z a -bso0 -bd -sse .\release\ezstore-portable.7z .\output\ezstore.exe .\cmd\README.txt .\cmd\update.ps1"

    Exec "iscc /Q 'setup.iss' /DPV='$ProductVersion' /DFV='$FileVersion'";

    Remove-Winres-Files;

}

##################
# Script beginning
##################

Check-If-Installed "Golang" "go";

Write-Host "Starting...";

$sw = [System.Diagnostics.Stopwatch]::New();
$sw.Start();

$ExitCode = 0;

try {
    switch ( $Command ) {
        'clean' {
            Clean;
            break;
        }
        'format' {
            Format;
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
        'check' {
            Clean;
            Format;
            Lint;
            Test;
            break;
        }
        'rebuild' {
            Clean;
            Build;
            break;
        }
    }
}
catch {
    $ExitCode = $global:LASTEXITCODE;
}

$sw.Stop();
$duration = Get-Duration $sw.Elapsed;

if ( $ExitCode -eq 0 ) {
    Write-Host "Finished $duration";
}
else {
    Write-Host "Failed with code $ExitCode $duration";
}

exit $ExitCode;
