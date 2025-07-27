[System.Diagnostics.CodeAnalysis.SuppressMessageAttribute(
    "PSUseDeclaredVarsMoreThanAssignments", "",
    Justification="Declared variables in Before* blocks available and used inside It blocks."
)]
param()

BeforeAll {
    . $PSCommandPath.Replace('.Tests.ps1','.ps1');

    $SkipInstallTests = $false;
    $PackageInstalledRegexp = 'Package ([a-zA-Z0-9.]+) v?([\d.]+) installed.$';
    $ColorRegexp = '\x1b\[[0-9;]*m';

    Import-ModuleSafe -Name "Appx" -UseWindowsPowerShell;

    function Get-PackageFullName {
        return Get-AppxPackage | ForEach-Object { $_.PackageFullName; };
    }

    # skipping 386 and ARM architectures on Windows 11 ARM64 on Github VMs due to this issue
    # https://learn.microsoft.com/en-us/windows/release-health/status-windows-11-21h2#2819msgdesc
    if ( $null -ne $Env:GITHUB_ACTION ) {
        $OSArch = [System.Runtime.InteropServices.RuntimeInformation,mscorlib]::OSArchitecture.ToString().ToLower();
        $OSBuild = [Environment]::OSVersion.Version.BuildNum;

        $Is32bitApp = ( $Arch -eq "386" ) -or ( $Arch -eq "arm" );
        $IsArm64Win11 = ( $OSBuild -ge 22000 ) -and ( $OSArch -eq "arm64" )

        if ( $IsArm64Win11 -and $Is32bitApp ) {
            $SkipInstallTests = $True;
        }
    }
}

Describe "Install subcommand (<arch>)" -ForEach $Targets {

    Context "positive tests" -Tag "Positive" -Skip:$SkipInstallTests {

        BeforeEach {
            $Before = Get-PackageFullName;
        }

        It "successfully install '<name>' v<version>" -ForEach @(
            @{
                Name = "Tree CLI"
                Id = "9mvsm3j7zj7c"
                Version = "1.1.0.0"
            },
            @{
                Name = "Wikipedia"
                Id = "9wzdncrfhwm4"
                Version = "1.0.0.0"
            },
            @{
                Name = "VPN Proxy: Fast & Unlimited"
                Id = "9pntscmcg01j"
                Version = "1.0.20.0"
            }
        ) {
            $Output, $Code = Invoke-EzstoreInstall $Path $Id $Version;

            $Code | Should -Be 0;
            $Output.Count | Should -Not -Be 0;
            $Output | Select-Object -Last 2 | Select-Object -First 1 | Should -Match $PackageInstalledRegexp;
        }

        It "successfully install without output color" {
            try {
                $OldValue = $Env:NO_COLOR; $Env:NO_COLOR = "1";
                $Output, $Code = Invoke-EzstoreInstall $Path "9mvsm3j7zj7c" "1.1.0.0";
            } finally {
                $Env:NO_COLOR = $OldValue;
            }

            $Code | Should -Be 0;
            $Output.Count | Should -Not -Be 0;
            $Output[0] | Should -Not -Match $ColorRegexp;
        }

        AfterEach {
            Get-PackageFullName | Where-Object { $Before -NotContains $_; } | ForEach-Object {
                try {
                    Remove-AppxPackage $_;
                } catch {
                    Write-Warning $_;
                }
            }
        }

    }

    Context "negative tests" -Tag "Negative" -Skip:$SkipInstallTests {

        It "fails to install unexisted app" {
            $Id = "f1o2o3b4a5r6";
            $Expected = '[ERR] Finished with error: can not fetch product info: product with id "' + $Id + '" and locale "en-US" not found';

            $Output, $Code = Invoke-EzstoreInstall $Path $Id "1.0.0.0";

            $Code | Should -Be 1;
            $Output.Count | Should -Not -Be 0;
            ($Output | Select-Object -Last 1) -replace $ColorRegexp | Should -BeExactly $Expected;
        }

    }

}
