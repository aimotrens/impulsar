#/usr/bin/env bash
_dothis_completions()
{
  COMPREPLY=($(compgen -W "$(impulsar --show-jobs)" "${COMP_WORDS[1]}"))
}

complete -F _dothis_completions impulsar
