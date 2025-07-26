BeforeAll {

    Push-Location;
    Set-Location $PSScriptRoot\..;

    try {
        $Version = (git describe --tags --abbrev=0) -split [Environment]::NewLine;
        $Version = ($Version | Select-Object -First 1) -replace "v";
    } finally {
        Pop-Location;
    }

}

Describe "Version flag (<arch>)" -ForEach $Targets {

    It "Returned correct app version (<_>)" -ForEach @(
        "-v"
        "--version"
    ) {
        $Output, $Code = Invoke-Ezstore $Path $_;

        $Code | Should -Be 0;
        $Output.Count | Should -Be 1;
        $Output[0] | Should -BeExactly "ezstore v$Version";
    }

}
