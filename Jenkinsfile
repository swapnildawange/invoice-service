pipeline {
    agent any 
    stages {
        stage('Checkout Codebase') {
            steps {
                checkout scm : [$class: 'GitSCM',branches:[[name:'*/develop']],
                userRemoteConfigs:[[credentialsId:'df4f5f85-8ea4-459c-a0f3-491ad36e9659',url:"https://github.com/swapnildawange/invoice-service.git"]]]
            }
        }
        stage('Build'){
            steps{
                echo 'Building CodeBase'
            }
        }

        stage('Test'){
            steps{
                echo 'Running Tests on changes'
            }
        }

        stage('Deploy'){
            steps{
                echo 'Done!'
            }
        }
    }
    // post {
    //     always {
    //         "go version"
    //     }
    // }
}   