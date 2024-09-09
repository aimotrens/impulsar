#compdef impulsar

# Get value after -f option
impulsar_get_f_value() {
  local args=("$@")  # All args as array
  for (( i = 1; i <= ${#args[@]}; i++ )); do
    if [[ "${args[i-1]}" == "-f" && -n "${args[i]}" ]]; then
      echo "${args[i]}" # Return value after -f
      return 0
    fi
  done
  return 1 # -f not found
}

# Generate completion for impulsar jobs
impulsar_get_jobs() {
  local provided_opt_f=$1

  local file_option=(-f impulsar.yml)
  if [[ -n "$provided_opt_f" ]]; then
    file_option=(-f "$provided_opt_f")
  fi

  local jobs
  jobs=("${(@f)$(impulsar show-jobs "${file_option[@]}" 2>/dev/null)}")
  _describe 'job' jobs
}

# Main completion function for impulsar
impulsar_completion() {
  local curcontext="$curcontext" state line

  _arguments -C \
    '1:subcommand:->subcmds' \
    '*::args:->args'

  case $state in
    (subcmds)
      local -a subcmds
      subcmds=(
        'run:run a job'
        'gen:generate shell completion'
        'show-jobs:show available jobs'
        'version:show version'
      )
      _describe 'subcommand' subcmds
      ;;
    (args)
      case $line[1] in
        (run)
          local f_value
          f_value=$(impulsar_get_f_value "${line[@]}")

          _arguments -s -C \
            '-f[impulsar file]:file:_files' \
            '-e[set environment variable]:name=value:' \
            '-dryrun[dryrun, only show execution plan]' \
            '-dump-jobs[dump parsed jobs to impulsar-dump.yml]' \
            '*::job:->jobs'

          case $state in
            (jobs)
              impulsar_get_jobs $f_value
              ;;
          esac
          ;;
        (gen)
          _arguments -s -C \
            '-o[output completion file]:file:_files' \
            '*::shell type:(bash zsh powershell)' \
          ;;
        (show-jobs)
          _arguments -s -C \
            '-f[impulsar file]:file:_files'
          ;;
        (version)
          # No additional arguments for 'version'
          ;;
      esac
      ;;
  esac
}

# Register the completion function with compdef
compdef impulsar_completion impulsar
