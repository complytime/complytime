#!/bin/bash
# SPDX-License-Identifier: Apache-2.0

#printf '    %q\n' "$0" "$@"
#set -x
#exec 2>&1

git ls-files -z | grep -z -f <(
    echo '\.go$'
    echo '\.sh$'
) | while IFS='' read -r -d '' file; do
    line_comment_token='//'
    case "$file" in
        *.go) line_comment_token='//';;
        *.sh) line_comment_token='#';;
    esac

    #echo "checking $file"
    # Check the top 3 lines for the license comment
    if ! head -n3 -- "$file" |
        grep -qFx \
            -e '// SPDX-License-Identifier: Apache-2.0' \
            -e '# SPDX-License-Identifier: Apache-2.0'; then
        printf '\e[0;31mLicense comment missing from file: \e[0;33m%s\e[0m\n' "$file"
        # Add the comment, with a blank line after
        git show HEAD:"$file" | awk -v line_comment_token="$line_comment_token" '
            # Skip shebang if it exists
            NR==1 && /^#!/ { print; next }
            !added {
                print line_comment_token " SPDX-License-Identifier: Apache-2.0"
                if(length) { print "" }
                if(!length) { getline; if(length) { print "" } }
                print
                added=1
                next
            }
            1' >"$file"
        printf '\e[0;31mPlease add and re-commit\e[0m\n'
        exit 1
    fi
done
