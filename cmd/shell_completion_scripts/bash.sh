# Get value after -f option
impulsar_get_f_value() {
  local args=("$@")  # All args as array
  for (( i = 1; i < ${#args[@]}; i++ )); do
    if [[ "${args[i-1]}" == "-f" && -n "${args[i]}" ]]; then
      echo "${args[i]}"  # Return value after -f
      return 0
    fi
  done
  return 1  # -f not found
}

# Generate completion for impulsar jobs
impulsar_get_jobs() {
  local provided_opt_f=$1

  local file_option="-f impulsar.yml"
  if [[ -n "$provided_opt_f" ]]; then
    file_option="-f $provided_opt_f"
  fi

  local jobs
  jobs=$(impulsar show-jobs $file_option 2>/dev/null)
  COMPREPLY=($(compgen -W "$jobs" -- "$cur"))
}

# Main completion function for impulsar
_impulsar_completion() {
  local cur prev opts
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"

  opts="run gen show-jobs version"

  case "${COMP_WORDS[1]}" in
    run)
      # If the current word starts with '-', complete with options for 'run'
      if [[ "$cur" == -* ]]; then
        COMPREPLY=($(compgen -W "-f -e -dryrun -dump-jobs" -- "$cur"))
        return 0
      fi

      # Handle the -f option specifically for file completion
      if [[ "$prev" == "-f" ]]; then
        COMPREPLY=($(compgen -f -- "$cur"))
        return 0
      fi

      # Handle job completion after other options
      if [[ "${COMP_WORDS[*]}" =~ "-f" ]]; then
        local f_value
        f_value=$(impulsar_get_f_value "${COMP_WORDS[@]}")
        impulsar_get_jobs "$f_value"
      else
        impulsar_get_jobs
      fi
      ;;
    gen)
      # If the current word starts with '-', complete with options for 'gen'
      if [[ "$cur" == -* ]]; then
        COMPREPLY=($(compgen -W "-o" -- "$cur"))
        return 0
      fi

      # Handle the -o option specifically
      if [[ "$prev" == "-o" ]]; then
        COMPREPLY=($(compgen -f -- "$cur"))
        return 0
      fi

      # Completion for the 'gen' subcommand
      COMPREPLY=($(compgen -W "bash zsh powershell" -- "$cur"))
      ;;
    show-jobs)
      # If the current word starts with '-', complete with options for 'run'
      if [[ "$cur" == -* ]]; then
        COMPREPLY=($(compgen -W "-f" -- "$cur"))
        return 0
      fi

      # Handle the -f option specifically for file completion
      if [[ "$prev" == "-f" ]]; then
        COMPREPLY=($(compgen -f -- "$cur"))
        return 0
      fi

      # Completion for the 'show-jobs' subcommand
      COMPREPLY=($(compgen -W "-f" -- "$cur"))
      ;;
    version)
      # No additional arguments for 'version'
      ;;
    *)
      # General completion for the first argument (subcommand)
      COMPREPLY=($(compgen -W "$opts" -- "$cur"))
      ;;
  esac
}

# Register the completion function
complete -F _impulsar_completion impulsar
