#compdef go-version

_go-version() {
    local -a commands
    commands=(
        'file:Generate a .version file'
        'ldflags:Generate ldflags for go build'
        'show:Display current git information'
        'version:Show go-version version'
        'help:Show help'
    )

    _arguments -C \
        '1: :->command' \
        '*: :->args'

    case $state in
        command)
            _describe -t commands 'go-version commands' commands
            ;;
        args)
            case $words[2] in
                file)
                    _arguments \
                        '(-o -output)'{-o,-output}'[Output file path]:file:_files' \
                        '(-v -version)'{-v,-version}'[Version string]:version:' \
                        '(-t -timestamp)'{-t,-timestamp}'[Build timestamp]:timestamp:' \
                        '-h[Show help]'
                    ;;
                ldflags)
                    _arguments \
                        '(-p -package)'{-p,-package}'[Package path]:package:' \
                        '(-v -version)'{-v,-version}'[Version string]:version:' \
                        '-static[Output static values]' \
                        '-shell[Output shell substitutions]' \
                        '-h[Show help]'
                    ;;
                show|version|help)
                    ;;
            esac
            ;;
    esac
}

_go-version "$@"
