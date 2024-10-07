#!/bin/bash

# Define cronjob
script_path="$(dirname "$(realpath "$0")")/generate.sh"
cronjob="* * * * * $script_path"

(crontab -l 2>/tmp/generate.log; echo "$cronjob") | crontab -
