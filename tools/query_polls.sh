#!/usr/bin/env bash

ENDPOINT="http://localhost:8000"

POLL_ID="$1"
DISCRIMINATOR="$2"

EAV=$(cat <<END
{
  ":pollId": {"S": "$POLL_ID"},
  ":discriminator": {"S": "$DISCRIMINATOR"}
}
END
)

aws --endpoint-url "$ENDPOINT" --output json dynamodb query \
  --table-name polls \
  --key-condition-expression "PollId = :pollId AND begins_with(Discriminator, :discriminator)" \
  --expression-attribute-values "$EAV"