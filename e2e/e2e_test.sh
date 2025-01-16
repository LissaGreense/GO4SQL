#!/bin/bash

e2e_failed=false

for test_file in e2e/test_files/*_test; do
  output_file="./e2e/$(basename "${test_file/_test/_output}")"
  ./GO4SQL -file "$test_file" > "$output_file"
  expected_output="e2e/expected_outputs/$(basename "${test_file/_test/_expected_output}")"
  diff "$output_file" "$expected_output"
  if [ $? -ne 0 ]; then
    echo "E2E test for: {$test_file} failed"
    e2e_failed=true
  fi
  rm "./$output_file"
done

if [ "$e2e_failed" = true ]; then
  echo "E2E tests failed."
  exit 1
else
  echo "All E2E tests passed."
fi
