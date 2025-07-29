[System.Diagnostics.CodeAnalysis.SuppressMessageAttribute(
    "PSUseDeclaredVarsMoreThanAssignments", "",
    Justification="Declared variables in Before* blocks available and used inside It blocks."
)]
param()

BeforeAll {
    . $PSCommandPath.Replace('.Tests.ps1','.ps1');

    $PackageInstalledRegexp = 'Package ([a-zA-Z0-9.]+) v?([\d.]+) installed.$';
    $ColorRegexp = '\x1b\[[0-9;]*m';

    Import-ModuleSafe -Name "Appx" -UseWindowsPowerShell;
}

Describe "Install subcommand (<arch>)" -Skip:$SkipInstallTests -ForEach $Targets {

    Context "positive tests" -Tag "Positive" {

        BeforeEach {
            $Before = Get-PackageFullName;
        }

        It "successfully install '<name>' v<version>" -ForEach @(
            @{
                Name = "Tree CLI"
                FullName = "PeterEtelej.TreeCLI"
                Version = "1.1.0.0"
                Id = "9mvsm3j7zj7c"
                PackageId = "vvysxk2z46ddc"
            },
            @{
                Name = "Wikipedia"
                FullName = "WikimediaFoundation.Wikipedia"
                Version = "1.0.0.0"
                Id = "9wzdncrfhwm4"
                PackageId = "54ggd3ev8bvz6"
            },
            @{
                Name = "VPN Proxy: Fast & Unlimited"
                FullName = "59992Roob.BestProxyFastUnlimitedVPNfunctionality"
                Version = "1.0.20.0"
                Id = "9pntscmcg01j"
                PackageId = "bzvrdnc3w98g4"
            }
        ) {
            $Output, $Code = Invoke-EzstoreInstall $Path $Id $Version;

            $Code | Should -Be 0;
            $Output | Should -AssertOutput -LineNum -2 -ShouldMatch $PackageInstalledRegexp;

            Assert-PackageInstalled -Name $FullName -Version $Version -PackageId $PackageId | Should -Be $true;
        }

        It "successfully install without output color" {
            $FullName = "PeterEtelej.TreeCLI";
            $Id = "9mvsm3j7zj7c";
            $Version = "1.1.0.0";
            $PackageId = "vvysxk2z46ddc"

            try {
                $OldValue = $Env:NO_COLOR; $Env:NO_COLOR = "1";
                $Output, $Code = Invoke-EzstoreInstall $Path $Id $Version;
            } finally {
                $Env:NO_COLOR = $OldValue;
            }

            $Code | Should -Be 0;
            $Output | Should -AssertOutput -Not -LineNum 0 -ShouldMatch $ColorRegexp;

            Assert-PackageInstalled -Name $FullName -Version $Version -PackageId $PackageId | Should -Be $true;
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

    Context "negative tests" -Tag "Negative" {

        It "fails to install unexisted app" {
            $Id = "f1o2o3b4a5r6";
            $Expected = '[ERR] Finished with error: can not fetch product info: product with id "' + $Id + '" and';
            $Expected += ' locale "en-US" not found';

            $Output, $Code = Invoke-EzstoreInstall $Path $Id "1.0.0.0";

            $Code | Should -Be 1;
            $Output | Should -AssertOutput -LineNum -1 -Script { $_ -replace $ColorRegexp; } -ShouldBeExactly $Expected;
        }

    }

}
