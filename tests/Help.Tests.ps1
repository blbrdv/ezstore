Describe "Help flag (<arch>)" -ForEach $Targets {

    BeforeAll {
        $HelpText = Get-Content -Raw -Path "$PSScriptRoot\..\cmd\README.txt" -WarningAction 'SilentlyContinue';
        $HelpText = $HelpText -join [Environment]::NewLine;
    }

    It "Returned correct help text (<_>)" -ForEach @(
        "-h"
        "--help"
    ) {
        $Output, $Code = Invoke-Ezstore $Path $_;

        $Code | Should -Be 0;
        $Output.Count | Should -Not -Be 0;

        $Output += ""; # Adding last empty line
        $Text = $Output -join [Environment]::NewLine;

        $Text | Should -BeExactly $HelpText;
    }

}
