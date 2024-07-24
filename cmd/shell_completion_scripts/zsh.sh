#compdef impulsar

_impulsar_jobs() {
  local jobs
  jobs=($(impulsar show-jobs))
  _describe 'jobs' jobs
}

compdef _impulsar_jobs impulsar
