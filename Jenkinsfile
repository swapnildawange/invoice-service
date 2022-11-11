pipeline {
    agent any 
    stages {
        stage('Checkout Codebase') {
            steps {
                checkout scm : [$class: 'GitSCM',branches:[[name:'*/develop']],
                userRemoteConfigs:[[credentialsId:'github-ssh-key',url:"git@github.com:swapnildawange/invoice-service.git"]]]
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
    post {
        always {
            "go version"
        }
    }
}   