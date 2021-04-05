
version="$(jq -r .version < ./release/release.json)"

go build -o evergreen-prompt

tar -zcvf evergreen-prompt-"${version}".tar.gz evergreen-prompt && rm evergreen-prompt

gh release create v"${version}" evergreen-prompt-"${version}".tar.gz --title "Evergreen Prompt v${version}" --draft --notes-file ./release/RELEASE_NOTES.md
