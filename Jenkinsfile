pipeline {
    agent any

    environment {
        GIT_REPO = 'https://github.com/malyshchyk/fantasy-backend.git'
        GO_VERSION = '1.22' // Adjust to your Go version
    }

    stages {
        stage('Clone Repository') {
            steps {
                git branch: 'main', url: "${env.GIT_REPO}"
            }
        }

        stage('Setup Go Environment') {
            steps {
                sh '''
                wget https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz
                sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
                export PATH=$PATH:/usr/local/go/bin
                '''
            }
        }

        stage('Build') {
            steps {
                sh 'go build ./...'
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./...'
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: '**/bin/*', allowEmptyArchive: true
            junit '**/test-reports/*.xml'
        }
    }
}
