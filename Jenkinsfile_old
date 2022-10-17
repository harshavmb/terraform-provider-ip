@Library(['pipeline-toolbox', 'iac-pipeline-shared-lib']) _

pipeline {
    agent any

    stages { 

        stage('Setup') {
            artifactName = env.BITBUCKET_REPOSITORY
            checkoutGit()            
            baseVersion = '1.0'
            newVersion = newVersion(baseVersion)
            registry = 'dockerhub.rnd.amadeus.net:5000'
            baseImage = 'maven:3.6.3-jdk-11'
            // pageId for https://rndwww.nce.amadeus.net/confluence/display/IBSDC/IaC+Release+Notes
            //releaseNotesOptions = ['spaceKey': 'IBSDC', 'parentPageId': 1654937615]
        }

        stage('Increment Tag && gorelease to artifactory') {
            steps {
              script {
                docker.image("dockerhub.rnd.amadeus.net/docker-production/iac/terraform-automation-azr:2.7.2").inside("-u iacuser") {
                    when(env.BRANCH_NAME == "master" || isReleasedBranch()) {
                        pushNewVersionTag(newVersion, baseVersion, releaseNotesOptions)
                        withCredentials([
                         string(credentialsId: 'GPG_FINGERPRINT', variable: 'GPG_FINGERPRINT'),
                         file(credentialsId: 'ash-gpg-key', variable: 'ASH_GPG_KEY'),
                         usernamePassword(credentialsId: '	goreleaser-artifactory-creds', usernameVariable: 'ARTIFACTORY_PRODUCTION_USERNAME', passwordVariable: 'ARTIFACTORY_PRODUCTION_SECRET')
                        ]) {
                            sh 'cp "${ASH_GPG_KEY}" .'
                            sh 'gpg --import ash-gpg-key && rm -rf ash-gpg-key'
                            sh 'wget -q -O /tmp/goreleaser.tar.gz https://github.com/goreleaser/goreleaser/releases/download/v1.11.4/goreleaser_Linux_x86_64.tar.gz'
                            sh 'tar -xf /tmp/goreleaser.tar.gz --directory /tmp/'
                            sh '/tmp/goreleaser release --rm-dist'
                            sh 'ls -lart dist/'
                        }
                    }                                    
                }           
              }
            }                        
        }
    }
    
}