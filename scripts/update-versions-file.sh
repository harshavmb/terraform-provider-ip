#!/bin/bash
set -e

provider=$(echo $PROJECT_NAME | rev | cut -d- -f1 | rev)
## create version from the passed one
function create_version() {
  cat <<EOF >new_version
{
 "versions":
 [
   {
     "version": "$VERSION",
     "protocols": ["4.0", "5.1"],
     "platforms": [
       {"os": "darwin", "arch": "amd64"},
       {"os": "darwin", "arch": "arm64"},
       {"os": "freebsd", "arch": "386"},
       {"os": "freebsd", "arch": "amd64"},
       {"os": "freebsd", "arch": "arm64"},
       {"os": "linux", "arch": "386"},
       {"os": "linux", "arch": "amd64"},
       {"os": "linux", "arch": "arm"},
       {"os": "linux", "arch": "arm64"},     
       {"os": "windows", "arch": "386"},
       {"os": "windows", "arch": "amd64"},
       {"os": "windows", "arch": "arm"} 
     ]
   }
 ]
}
EOF
}

## This will update the version by merging the latest one with the existing one
function update_version() {
  echo "{\"versions\": $(jq -s '.[0].versions + .[1].versions | unique' versions new_version) }" >all_versions
}

## This will create the version for the first time
## merge with new_version twice
function create_initial_version() {
  echo "{\"versions\": $(jq -s '.[0].versions + .[1].versions | unique' new_version new_version) }" >all_versions
}

## Upload the final file to the artifactory
function upload_to_artifactory() {
  curl -k --user $ARTIFACTORY_USER:$ARTIFACTORY_PASSWORD --upload-file all_versions -X PUT "${ARTIFACTORY_URL}/artifactory/generic-production-iac/terraform/providers/v1/amadeus/$provider/versions"
}

## Check whether new_version file exists in the current dir
## This crappy logic is to avoid concurrent runs to push same versions file onto artifactory
## if we use commercial version of goreleaser, this can be avoided with post hooks (after)
## neverthless, this is for improvement. It first checks whether new_version file exists, if yes it skips all execution..
## so that only one build runs parallelly
if [[ "$LINUX_AMD64" == "linux_amd64" ]]; then
  ## extract the provider name from the project name

  ## download the versions file if exists  
  response=$(curl -k -Is ${ARTIFACTORY_URL}/artifactory/generic-production-iac/terraform/providers/v1/amadeus/$provider/versions | head -1 | cut -d ' ' -f2)  

  if [[ $response -eq "200" ]]; then
    echo "versions file exists, downloading it"
    curl -k -s -o versions ${ARTIFACTORY_URL}/artifactory/generic-production-iac/terraform/providers/v1/amadeus/$provider/versions
    create_version
    update_version
    upload_to_artifactory
  elif [[ $response -eq "404" ]]; then
    echo "versions file does not exist, it's a new provider"
    create_version
    create_initial_version
    upload_to_artifactory
  else
    echo "Some error while downloading versions file from artifactory"
    exit 1
  fi
else
  echo "not linux_amd64, skipping rest of the execution"
fi

echo $?
