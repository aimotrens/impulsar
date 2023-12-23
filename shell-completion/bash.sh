#/usr/bin/env bash
_impulsar_completions()
{
  COMPREPLY=($(compgen -W "$(impulsar --show-jobs)" "${COMP_WORDS[1]}"))
}

complete -F _impulsar_completions impulsar
