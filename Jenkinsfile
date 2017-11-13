#!/usr/bin/groovy
@Library('github.com/fabric8io/fabric8-pipeline-library@master')
def dummy
goNode{
  dockerNode{
    ws{
      if (env.BRANCH_NAME.startsWith('PR-')) {
        goCI{
          githubOrganisation = 'rawlingsj'
          dockerOrganisation = 'fabric8'
          project = 'semver-release-number'
          makeTarget = 'clean test cross'
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

