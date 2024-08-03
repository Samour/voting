#!/usr/bin/env bash

ENDPOINT="http://localhost:8000"

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

create_table() {
  TABLE_NAME="$1"
  FILE_NAME="$SCRIPT_DIR/CreateTable-$TABLE_NAME.json"
  ATTRIBUTE_DEFINITIONS=$(jq '.AttributeDefinitions' "$FILE_NAME")
  KEY_SCHEMA=$(jq '.KeySchema' "$FILE_NAME")

  AWS_ACCESS_KEY_ID=dummy AWS_SECRET_ACCESS_KEY=dummy AWS_DEFAULT_REGION=local aws --endpoint-url "$ENDPOINT" \
    dynamodb create-table \
    --table-name "$TABLE_NAME" \
    --attribute-definitions "$ATTRIBUTE_DEFINITIONS" \
    --key-schema "$KEY_SCHEMA" \
    --billing-mode "PAY_PER_REQUEST" >/dev/null
  
  echo "Created table $TABLE_NAME"
}

create_table polls
create_table users
create_table username-password-credentials

echo "All tables created"
