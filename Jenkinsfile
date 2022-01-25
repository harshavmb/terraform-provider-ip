/* groovylint-disable-next-line UnusedVariable,VariableName,CompileStatic */
@Library(['pipeline-toolbox', 'iac-pipeline-shared-lib']) _

node {
    try {

        stage('Setup') {
            artifactName = env.BITBUCKET_REPOSITORY
            checkoutGit()
            terraformAutoImageVersion = 'latest'
            baseTerraformAutoImage = "iac/terraform-automation-base-ash:${terraformAutoImageVersion}"
            baseVersion = '1.0'
            newVersion = newVersion(baseVersion)
            registry = 'dockerhub.rnd.amadeus.net:5000'
            baseImage = 'maven:3.6.3-jdk-11'
            // pageId for https://rndwww.nce.amadeus.net/confluence/display/IBSDC/IaC+Release+Notes
            //releaseNotesOptions = ['spaceKey': 'IBSDC', 'parentPageId': 1654937615]
        }

        stage('QA tests and goreleaser to release locally (not pushing to artifactory yet)') {
            when(isPullRequest()) {
                docker.withRegistry("https://${registry}", 'IZ_USER') {
                    docker.image(baseTerraformAutoImage).inside {
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
          withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'goreleaser-artifactory-creds', usernameVariable: 'ARTIFACTORY_PRODUCTION_USERNAME', passwordVariable: 'ARTIFACTORY_PRODUCTION_SECRET']]) {  
            when(env.BRANCH_NAME == "main" || isReleasedBranch()) {
                pushNewVersionTag(newVersion, baseVersion, releaseNotesOptions)
                docker.withRegistry("https://${registry}", 'IZ_USER') {
                    docker.image(baseTerraformAutoImage).inside {
                        sh '''
                          goreleaser release --rm-dist
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