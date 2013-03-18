#!/bin/bash

LAST_RET_CODE=0

if [ $LAST_RET_CODE == 0 ]; then
  bundle install
  LAST_RET_CODE=$?
fi

if [ $LAST_RET_CODE == 0 ]; then
  librarian-puppet install --verbose
  LAST_RET_CODE=$?
fi

if [ $LAST_RET_CODE == 0 ]; then
  bundle exec rake $@
else
  exit $LAST_RET_CODE
fi
