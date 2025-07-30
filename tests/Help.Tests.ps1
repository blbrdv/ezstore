[System.Diagnostics.CodeAnalysis.SuppressMessageAttribute(
    "PSUseDeclaredVarsMoreThanAssignments", "",
    Justification="Declared variables in Before* blocks available and used inside It blocks."
)]
param()

Describe "Help flag (<arch>)" -ForEach $Targets {

    Context "positive tests" -Tag "Positive" {

        BeforeAll {
            $HelpText = Get-Content -Path "$PSScriptRoot\..\cmd\README.txt" -WarningAction 'SilentlyContinue';
        }

        It "returned correct help text (<_>)" -ForEach @(
            "-h"
            "--help"
        ) {
            $Output, $Code = Invoke-Ezstore $Path @($_);

            $Code | Should -Be 0;
            $Output.Count | Should -Not -Be 0;

            $Output += ""; # Adding last empty line
            $Text = $Output -join [Environment]::NewLine;

            $Text | Should -BeExactly $HelpText;
        }

    }

}
