#!/bin/sh
set -eu

base_url="${BASE_URL:-http://localhost:8080}"
book_path="${BOOK_PATH:-examples/pride-and-prejudice/pride-and-prejudice.epub}"
book_b64_file="$(mktemp)"
trap 'rm -f "$book_b64_file"' EXIT INT TERM

token="$(
  curl -fsS "$base_url/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin","twoFactor":"123456"}' |
    jq -r .token
)"

base64 < "$book_path" | tr -d '\n' > "$book_b64_file"

publication_id="$(
  jq -cn --rawfile file "$book_b64_file" '{
      query: "mutation UploadPublication($title: String!, $file: Upload!) { uploadPublication(title: $title, file: $file) { id } }",
      variables: {
        title: "Pride and Prejudice",
        file: $file
      }
    }' |
    curl -fsS "$base_url/graphql" \
      -H "Authorization: Bearer $token" \
      -H "Content-Type: application/json" \
      --data-binary @- |
    jq -er .data.uploadPublication.id
)"

license_id="$(
  jq -cn --arg publicationID "$publication_id" '{
      query: "mutation CreateLicense($publicationID: ID!, $userID: ID!, $passphrase: String!, $hint: String!) { createLicense(publicationID: $publicationID, userID: $userID, passphrase: $passphrase, hint: $hint) { id } }",
      variables: {
        publicationID: $publicationID,
        userID: "reader-01",
        passphrase: "open-sesame",
        hint: "demo"
      }
    }' |
    curl -fsS "$base_url/graphql" \
      -H "Authorization: Bearer $token" \
      -H "Content-Type: application/json" \
      --data-binary @- |
    jq -er .data.createLicense.id
)"

echo "publication_id=$publication_id"
echo "license_id=$license_id"
echo "license_url=$base_url/licenses/$license_id.lcpl"
