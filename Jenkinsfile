@Library("edge-jenkins-lib") _

pipeline {
    agent { 
        kubernetes {
            defaultContainer 'ucpe-jenkins-dind'
            yaml readTrusted('pod.yaml')
        }
    }
    stages {
        stage('Setup') {
          steps {
            script {
                dockerSetup()
            }
          }
        }
        stage('Build') {
            steps {
                script {
                    // Build and run services from docker/docker-compose.yml
                    buildServices{ logFailedServices = ["netapp-exporter"] }
                }
            }
        }
        stage('Publish') {
            when {
                expression { currentBuild.result == 'SUCCESS' }
            }
            steps {
                script {
                    // Tag and push images to registry
                    pushImages{ images = ["netapp-exporter"] }
                }
            }
        }
    }
}