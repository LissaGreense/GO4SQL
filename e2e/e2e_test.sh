#!/bin/bash

e2e_failed=false

for test_file in e2e/test_files/*_test; do
  name="$(basename "${test_file%_test}")"
  output_file="./e2e/${name}_output"
  expected_output="e2e/expected_outputs/${name}_expected_output"
  expected_exit_file="e2e/expected_outputs/${name}_expected_exit"

  ./GO4SQL -file "$test_file" > "$output_file"
  actual_exit=$?

  expected_exit=0
  if [ -f "$expected_exit_file" ]; then
    expected_exit="$(tr -d '[:space:]' < "$expected_exit_file")"
  fi

  # Exit code is checked separately because Evaluate discards earlier output
  # on error, so a failing fixture can still match an empty expected output.
  if ! diff "$output_file" "$expected_output"; then
    echo "E2E test for: {$test_file} failed: stdout differs from $expected_output"
    e2e_failed=true
  fi
  if [ "$actual_exit" != "$expected_exit" ]; then
    echo "E2E test for: {$test_file} failed: expected exit code $expected_exit, got $actual_exit"
    e2e_failed=true
  fi

  rm "$output_file"
done

if [ "$e2e_failed" = true ]; then
  echo "E2E tests failed."
  exit 1
else
  echo "All E2E tests passed."
fi
