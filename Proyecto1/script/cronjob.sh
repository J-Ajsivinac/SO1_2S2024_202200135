#!/bin/bash

# Define cronjob
script_path="$(pwd)/generate.sh"
cronjob="* * * * * $script_path"

(crontab -l 2>/dev/null; echo "$cronjob") | crontab -
