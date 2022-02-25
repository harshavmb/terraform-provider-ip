#!/bin/bash
set -e

## extract the provider name from the project name
provider=$(echo $4 | rev | cut -d- -f1 | rev)

## download the versions file if exists
response=$(curl -k -Is https://repository.rnd.amadeus.net/artifactory/generic-production-iac/terraform/providers/v1/amadeus/$provider/versions | head -1 | cut -d ' ' -f2)

## create version from the passed one
function create_version() {
   cat <<EOF > new_version
{
 "versions":
 [
   {
     "version": "$1",
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
  echo "{\"versions\": $(jq -s '.[0].versions + .[1].versions | unique' versions new_version) }" > all_versions
}

## This will create the version for the first time
## merge with new_version twice
function create_initial_version() {
  echo "{\"versions\": $(jq -s '.[0].versions + .[1].versions | unique' new_version new_version) }" > all_versions
}

## Upload the final file to the artifactory
function upload_to_artifactory() {
  curl -k --user $1:$2 --upload-file all_versions -X PUT "https://repository.rnd.amadeus.net/artifactory/generic-production-iac/terraform/providers/v1/amadeus/$provider/versions"
}

if [[ $response -eq "200" ]]
then
   echo "versions file exists, downloading it"
   curl -k -s -o versions https://repository.rnd.amadeus.net/artifactory/generic-production-iac/terraform/providers/v1/amadeus/$provider/versions
   create_version $1
   update_version
   upload_to_artifactory $2 $3
elif [[ $response -eq "404" ]]
then
   echo "versions file does not exist, it's a new provider"
   create_version $1
   create_initial_version
   upload_to_artifactory $2 $3
else
   echo "Some error while downloading versions file from artifactory"
   exit 1
fi

echo $?