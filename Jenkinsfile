pipeline {
    agent any
    stages {
        stage('CI Build and Test') {
            when {
                branch 'PR-*'
            }
            steps {
                dir ('/home/jenkins/go/src/github.com/jenkins-x/jx-release-version') {
                    checkout scm
                    sh "make"
                    sh "./bin/jx-release-version-linux"
                }
            }
        }
    
        stage('Build and Release') {
            environment {
                GH_CREDS = credentials('jx-pipeline-git-github-github')
            }
            when {
                branch 'master'
            }
            steps {
                dir ('/home/jenkins/go/src/github.com/jenkins-x/jx-release-version') {
                    git "https://github.com/jenkins-x/jx-release-version"
                    
                    sh "GITHUB_ACCESS_TOKEN=$GH_CREDS_PSW make release"
                }
            }
        }
    }
}
