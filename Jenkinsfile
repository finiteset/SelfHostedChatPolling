pipeline {
    agent any

    tools {
        tool name: 'Go 1.8', type: 'go'
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