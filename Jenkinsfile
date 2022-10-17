/* groovylint-disable-next-line UnusedVariable,VariableName,CompileStatic */
@Library(['pipeline-toolbox', 'iac-pipeline-shared-lib']) _

node {
    try {

        stage('Setup') {
            artifactName = env.BITBUCKET_REPOSITORY
            checkoutGit()
            terraformAutoImageVersion = '2.7.2'
            baseTerraformAutoImage = "docker-production/iac/terraform-automation-azr:${terraformAutoImageVersion}"
            baseVersion = '1.0'
            newVersion = newVersion(baseVersion)
            registry = 'dockerhub.rnd.amadeus.net'
            baseImage = 'maven:3.6.3-jdk-11'
            // pageId for https://rndwww.nce.amadeus.net/confluence/display/IBSDC/IaC+Release+Notes
            //releaseNotesOptions = ['spaceKey': 'IBSDC', 'parentPageId': 1654937615]
        }

        stage('QA tests and goreleaser to release locally (not pushing to artifactory yet)') {
            when(isPullRequest()) {
                docker.withRegistry("https://${registry}", 'IZ_USER') {
                    withCredentials([
                         string(credentialsId: 'GPG_FINGERPRINT', variable: 'GPG_FINGERPRINT'),
                         file(credentialsId: 'ash-gpg-key', variable: 'ASH_GPG_KEY'),
                         usernamePassword(credentialsId: '	goreleaser-artifactory-creds', usernameVariable: 'ARTIFACTORY_PRODUCTION_USERNAME', passwordVariable: 'ARTIFACTORY_PRODUCTION_SECRET')
                        ]) {
                            sh '''
                              make fmtcheck
                              make test
                              goreleaser release --snapshot --rm-dist
                        '''
                        }                    
                }
            }
        }

        stage ('Increment Tag && gorelease to artifactory') {
          withCredentials([
                string(credentialsId: 'GPG_FINGERPRINT', variable: 'GPG_FINGERPRINT'),
                file(credentialsId: 'ash-gpg-key', variable: 'ASH_GPG_KEY'),
                usernamePassword(credentialsId: '	goreleaser-artifactory-creds', usernameVariable: 'ARTIFACTORY_PRODUCTION_USERNAME', passwordVariable: 'ARTIFACTORY_PRODUCTION_SECRET')
          ]) {  
            when(env.BRANCH_NAME == "master" || isReleasedBranch()) {
                pushNewVersionTag(newVersion, baseVersion, releaseNotesOptions)
                docker.withRegistry("https://${registry}", 'IZ_USER') {
                    docker.image(baseTerraformAutoImage).inside {
                        sh '''
                          cp "${ASH_GPG_KEY}" .
                          gpg --import ash-gpg-key && rm -rf ash-gpg-key
                          wget -q -O /tmp/goreleaser.tar.gz https://github.com/goreleaser/goreleaser/releases/download/v1.11.4/goreleaser_Linux_x86_64.tar.gz
                          tar -xf /tmp/goreleaser.tar.gz --directory /tmp/
                          wget -q -O /tmp/tf-provider-registry-api-generator.zip https://repository.adp.amadeus.net/generic-production-iac/binaries/tf-provider-registry-api-generator/1.0.1/tf-provider-registry-api-generator_1.0.1_linux_amd64.zip
                          unzip -o /tmp/tf-provider-registry-api-generator.zip -d /tmp/ && mv /tmp/tf-provider-registry-api-generator* /tmp/tf-provider-registry-api-generator
                          /tmp/goreleaser release --rm-dist                                                   
                        '''
                    }
                }
            }
          }
        }
    } catch (err) {
        echo "Caught: ${err}"
        currentBuild.result = 'FAILURE'
    } finally {
        echo 'Done'
    }
}