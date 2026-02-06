#!/bin/bash
# bash completion for go-version

_go_version() {
    local cur prev commands
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    commands="file ldflags show version help"

    case "${prev}" in
        go-version)
            COMPREPLY=( $(compgen -W "${commands}" -- "${cur}") )
            return 0
            ;;
        file)
            COMPREPLY=( $(compgen -W "-o -output -v -version -t -timestamp -h" -- "${cur}") )
            return 0
            ;;
        ldflags)
            COMPREPLY=( $(compgen -W "-p -package -v -version -static -shell -h" -- "${cur}") )
            return 0
            ;;
        -o|-output)
            COMPREPLY=( $(compgen -f -- "${cur}") )
            return 0
            ;;
        *)
            ;;
    esac

    COMPREPLY=( $(compgen -W "${commands}" -- "${cur}") )
}

complete -F _go_version go-version
