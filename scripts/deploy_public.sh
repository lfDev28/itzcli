# these come from .travis.yml script cmd
TOKEN=$1
TAG=$2

echo $TAG

GENERATE_POST_BODY() {
  cat <<EOF
{
  "tag_name": "$TAG",
  "target_commitish": "main",
  "name": "$TAG",
  "body": "release version: $TAG",
  "draft": false,
  "prerelease": false
}
EOF
}

API_RESPONSE_STATUS=$(
  curl \
    --header 'Accept: application/vnd.github.v3+json' \
    --header "Authorization: token $TOKEN" \
    --header 'Content-Type: application/json' \
    --data "$(GENERATE_POST_BODY)" \
    -s \
    https://api.github.com/repos/lfDev28/itzcli/releases
)
echo API_RESPONSE:
echo "$API_RESPONSE_STATUS"
RELEASE_ID=$(echo $API_RESPONSE_STATUS | jq '.id')
echo $RELEASE_ID

BIN_DARWIN="itzcli-darwin-amd64.tar.gz"
BIN_LINUX="itzcli-linux-amd64.tar.gz"

# Construct url
GH_ASSET_DARWIN="https://uploads.github.com/repos/lfDev28/itzcli/releases/$RELEASE_ID/assets?name=$(basename $BIN_DARWIN)"
GH_ASSET_LINUX="https://uploads.github.com/repos/lfDev28/itzcli/releases/$RELEASE_ID/assets?name=$(basename $BIN_LINUX)"

curl --data-binary @"$BIN_DARWIN" -H "Authorization: token $TOKEN" -H "Content-Type: application/octet-stream" $GH_ASSET_DARWIN
curl --data-binary @"$BIN_LINUX" -H "Authorization: token $TOKEN" -H "Content-Type: application/octet-stream" $GH_ASSET_LINUX