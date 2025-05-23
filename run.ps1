Param (
    [Parameter(Mandatory=$true,Position=0)]
    [ValidateNotNullOrEmpty()]
    [ValidateSet('clean','format','check','test','deps','build','rebuild')]
    [string]$Command
)

$global:SysoFiles = @(
    "rsrc_windows_386.syso",
    "rsrc_windows_amd64.syso"
);

$global:BuildDirs = @(
    "output",
    "release"
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

    return (git describe --tags --abbrev=0).Replace("v", "");

}

function Get-File-Version {

    Param (
        [Parameter(Mandatory=$true,Position=0)]
        [string]$Version
    )

    $Value = "0";

    $LastTag = git describe --tags --abbrev=0;

    $Count = Invoke-Expression "git log $LastTag..HEAD --oneline" | Measure-Object -Line | %{$_.Lines};

    if ( $Count -ne "" ) {
        $Value = $Count;
    }

    return "${Version}.${Value}";

}

function Exec {

    Param (
        [Parameter(Mandatory=$true,Position=0)]
        [string]$Name,
        [Parameter(Mandatory=$true,Position=1)]
        [scriptblock]$Command
    )

    $TaskName = (Get-PSCallStack)[1].Command;
    Write-Host " > ${TaskName}: $Name";

    $global:LASTEXITCODE = 0;

    & $Command;

    if ( $global:LASTEXITCODE -ne 0 ) {
        throw "Command exited with code $global:LASTEXITCODE";
    }

}

#######
# Tasks
#######

function Remove-Winres-Files {

    foreach ($File in $global:SysoFiles) {
        $Path = ".\cmd\$File";
        Exec "Removing $Path" { Remove-Item -Path $Path -Force -ErrorAction SilentlyContinue };
    }

}

function Clean {

    foreach ($Dir in $global:BuildDirs) {
        Exec "Removing $Dir" { Remove-Item -Path $Dir -Recurse -Force -ErrorAction SilentlyContinue };
    }

    Remove-Winres-Files;

}

function Format {

    Exec "Formatting go files" { go fmt .\... };

}

function Check {

    Check-If-Installed "Staticcheck" "staticcheck";

    Exec "Checking issues" { go vet .\... };
    Exec "Checking codestyle" { staticcheck .\... };
    Exec "Checking format" {
            $Location = (Get-Location | %{$_.Path}) + "\";
            $Files = Get-Childitem –Path . -Include *.go -Recurse -ErrorAction SilentlyContinue |
                %{$_.FullName.Replace($Location,'')};
            $Result = gofmt -l $Files;

            if ($null -ne $Result -And ($Result | Measure-Object -Line | %{$_.Lines} -gt 0)) {
                Write-Host "Code need formatting:";
                Write-Host $Result;
                Write-Host;
                $global:LASTEXITCODE = 1;
            }
        }

}

function Test {

    Exec "Running tests" { go test .\internal\... };

}

function Dependencies {

    Exec "Check dependencies for update" {
            $Paths = go list -m -u -f '{{if not (or .Indirect .Main)}}{{with .Update}}{{$.Path}} {{$.Version}} {{.Version}}{{end}}{{end}}' all;

            $Deps = @();
            foreach ($Path in $Paths) {
                $Result = $Path -match "^(\S+) v(\S+) v(\S+)$";
                if ( $Result ) {
                    $Current = [version]$Matches[2];
                    $New = [version]$Matches[3];
                    if ( $New.Major -gt $Current.Major -Or ($New.Major -eq $Current.Major -And $New.Minor -gt $Current.Minor) ) {
                        $Deps += $Path
                    }
                }
            }

            if ( $Deps.Length -gt 0 ) {
                Write-Host "Dependencies need updates:";
                foreach ($Dep in $Deps) {
                    Write-Host $Dep;
                }
                Write-Host;

                $global:LASTEXITCODE = 1;
            }
        }

}

function Build {

    Check-If-Installed "go-winres" "go-winres";
    Check-If-Installed "7-Zip" "7z";
    Check-If-Installed "Inno Setup" "iscc";

    $ProductVersion = Get-Product-Version;
    $FileVersion = Get-File-Version $ProductVersion;

    Write-Host "Building project version $ProductVersion ($FileVersion)";

    Exec "Embedding resources" {
            go-winres make --in ".\winres.json" --product-version $ProductVersion --file-version $FileVersion
        };

    foreach ($File in $global:SysoFiles) {
        $Target = ".\cmd\$File";
        Exec "Moving $File to $Target" {
                Move-Item -Path $File -Destination $Target -Force -ErrorAction SilentlyContinue
            };
    }

    Exec "Compiling exe" { go build -ldflags="-X main.version=$ProductVersion" -o ".\output\bin\ezstore.exe" ".\cmd" };

    Remove-Winres-Files;

    $DistFiles = Get-ChildItem -Path .\cmd -Exclude "*.go" | %{$_.Name};
    foreach ($File in $DistFiles) {
        $Path = ".\cmd\$File"
        Exec "Copying $Path to output" { Copy-Item -Path $Path -Destination ".\output" }
    }

    Exec "Archiving files" {
            7z a -bso0 -bd -sse ".\release\ezstore-portable.7z" ".\output\bin\ezstore.exe" ".\output\README.txt" ".\output\update.ps1"
        };

    Exec "Compiling installer" { iscc /Q "setup.iss" /DPV=$ProductVersion /DFV=$FileVersion };

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
        'check' {
            Check;
            break;
        }
        'test' {
            Test;
            break;
        }
        'deps' {
            Dependencies;
            break;
        }
        'build' {
            Build;
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
    Write-Host "Success $duration";
}
else {
    Write-Host "Failed with code $ExitCode $duration";
}

exit $ExitCode;
