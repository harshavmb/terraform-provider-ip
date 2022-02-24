#!/bin/bash

## create a json metadata file of tf plugin version
shasum = $(sha256sum "dist/$3.zip")

## extract the provider name from the project name
provider = $(echo $3 | rev | cut -d- -f1 | rev)

## with the above data, now we create metadata file to upload to artifactory
cat <<EOF > dist/$4_$3_$1_$2_metadata
{
"protocols": [
"5.1"
],
"os": "$1",
"arch": "$2",
"filename": "$4_$3_$1_$2.zip",
"download_url": "https://repository.adp.amadeus.net/generic-production-iac/providers/$provider/$BITBUCKET_REPOSITORY/$3/$4_$3_$1_$2.zip",
"shasums_url": "https://repository.adp.amadeus.net/generic-production-iac/providers/$provider/$BITBUCKET_REPOSITORY/$3/$4_$3_SHA256SUMS",
"shasums_signature_url": "https://repository.adp.amadeus.net/generic-production-iac/providers/$provider/$BITBUCKET_REPOSITORY/$3/$4_$3_SHA256SUMS.sig",
"shasum": "$shasum",
"signing_keys": {
"gpg_public_keys": [
{
"key_id": "$GPG_FINGERPRINT",
"ascii_armor": "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\n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\nSS6vklk8KN0VJ0axDcZndwfKt6Rl\n=ZG6A\n-----END PGP PUBLIC KEY BLOCK-----",
"trust_signature": "",
"source": "Amadeus",
"source_url": "https://www.amadeus.com/"
}
]
}
}
EOF

## upload the created file to artifactory
curl --user $ARTIFACTORY_PRODUCTION_USERNAME:$ARTIFACTORY_PRODUCTION_SECRET --data-binary @dist/$4_$3_$1_$2_metadata -X PUT "https://$ARTIFACTORY_URL/artifactory/generic-production-iac/terraform/providers/v1/amadeus/$provider/$3/download/$1/$2"