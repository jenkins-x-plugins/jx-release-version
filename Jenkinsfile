#!/usr/bin/groovy
@Library('github.com/fabric8io/fabric8-pipeline-library@master')
def dummy
goNode{
  dockerNode{
    ws{
      if (env.BRANCH_NAME.startsWith('PR-')) {
        def buildPath = "/home/jenkins/go/src/github.com/rawlingsj/semver-release-number"
        sh "mkdir -p ${buildPath}"

        dir(buildPath) {
            container(name: 'go') {
                stage ('build binary'){
                    // it looks like using checkout scm looses the tags
                    sh "git clone https://github.com/rawlingsj/semver-release-number ."
                    sh "git fetch origin pull/${env.CHANGE_ID}/head:test"
                    sh "git checkout test"
                    sh "make"
                }
            }
        }
      } else if (env.BRANCH_NAME.equals('master')) {
        def v = goRelease{
          githubOrganisation = 'rawlingsj'
          dockerOrganisation = 'fabric8'
          project = 'semver-release-number'
        }
      }
    }
  }
}

