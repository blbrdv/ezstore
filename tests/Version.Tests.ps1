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

    Context "positive tests" -Tag "Positive" {

        It "returned correct app version (<_>)" -ForEach @(
            "-v"
            "--version"
        ) {
            $Output, $Code = Invoke-Ezstore $Path @($_);

            $Code | Should -Be 0;
            $Output | Should -AssertOutput -LineNum 0 -ShouldMatch "ezstore v$Version";
        }

    }

}
