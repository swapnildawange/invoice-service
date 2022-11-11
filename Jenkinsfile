pipeline {
    agent any 
    tools {
        go '1.19.3'
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
    }
    stages {
        stage('Checkout Codebase') {
            steps {
                checkout scm : [$class: 'GitSCM',branches:[[name:'*/develop']],
                userRemoteConfigs:[[credentialsId:'df4f5f85-8ea4-459c-a0f3-491ad36e9659',url:"https://github.com/swapnildawange/invoice-service.git"]]]
            }
        }
        stage('Pre Test') {
            steps {
                echo 'Installing dependencies'
                sh 'go version'
                sh 'go get -u golang.org/x/lint/golint'
            }
        }
        
        stage('Build') {
            steps {
                echo 'Compiling and building'
                sh 'go build ./cmd/main.go'
            }
        }

        stage('Test') {
            steps {
                withEnv(["PATH+GO=${GOPATH}/bin"]){
                    echo 'Running vetting'
                    sh 'go vet ./...'
                    echo 'Running test'
                    sh 'go test ./...'
                }
            }
        }
    }
    // post {
    //     always {
    //         "go version"
    //     }
    // }
}   