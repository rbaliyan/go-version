# fish completion for go-version

# Disable file completion by default
complete -c go-version -f

# Commands
complete -c go-version -n "__fish_use_subcommand" -a "file" -d "Generate a .version file"
complete -c go-version -n "__fish_use_subcommand" -a "ldflags" -d "Generate ldflags for go build"
complete -c go-version -n "__fish_use_subcommand" -a "show" -d "Display current git information"
complete -c go-version -n "__fish_use_subcommand" -a "version" -d "Show go-version version"
complete -c go-version -n "__fish_use_subcommand" -a "help" -d "Show help"

# file subcommand options
complete -c go-version -n "__fish_seen_subcommand_from file" -s o -l output -d "Output file path" -r
complete -c go-version -n "__fish_seen_subcommand_from file" -s v -l version -d "Version string" -r
complete -c go-version -n "__fish_seen_subcommand_from file" -s t -l timestamp -d "Build timestamp" -r
complete -c go-version -n "__fish_seen_subcommand_from file" -s h -d "Show help"

# ldflags subcommand options
complete -c go-version -n "__fish_seen_subcommand_from ldflags" -s p -l package -d "Package path" -r
complete -c go-version -n "__fish_seen_subcommand_from ldflags" -s v -l version -d "Version string" -r
complete -c go-version -n "__fish_seen_subcommand_from ldflags" -l static -d "Output static values"
complete -c go-version -n "__fish_seen_subcommand_from ldflags" -l shell -d "Output shell substitutions"
complete -c go-version -n "__fish_seen_subcommand_from ldflags" -s h -d "Show help"
