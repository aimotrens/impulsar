# Impulsar Completion
Register-ArgumentCompleter -CommandName impulsar -Native -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    $jobs = $(impulsar --show-jobs).Split("\n")

    $jobs `
    | Where-Object { $_ -like "${wordToComplete}*" } `
    | ForEach-Object { [CompletionResult]::new($_, $_, [CompletionResultType]::ParameterValue, $_) }
}