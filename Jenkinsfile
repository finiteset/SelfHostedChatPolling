pipeline {
    agent any

    tools {
        go 'Go 1.8'
    }

    stages {
        stage('Build') {
            steps {
                sh 'go build'
            }
        }
        stage('Test') {
            steps {
                sh 'go test'
            }
        }
    }
}