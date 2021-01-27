pipeline {
    agent any
    
    stages {
        stage("Build linux/amd64") {
            agent {
                docker { 
                    image 'golang:buster'
                    args '-u root --privileged'
                }
            }
            environment {
                GOOS = 'linux'
                GOARCH = 'amd64'
                BINNAME = 'agent'
            }
            steps {
                build_binary()
            }
        }
        stage("Build linux/386") {
            agent {
                docker { 
                    image 'golang:buster'
                }
            }
            environment {
                GOOS = 'linux'
                GOARCH = '386'
                BINNAME = 'agent'
            }
            steps {
                build_binary()
            }
        }
        stage("Build linux/arm") {
            agent {
                docker { 
                    image 'golang:buster'
                }
            }
            environment {
                GOOS = 'linux'
                GOARCH = 'arm'
                BINNAME = 'agent'
            }
            steps {
                build_binary()
            }
        }
        stage("Build linux/arm64") {
            agent {
                docker { 
                    image 'golang:buster'
                }
            }
            environment {
                GOOS = 'linux'
                GOARCH = 'arm64'
                BINNAME = 'agent'
            }
            steps {
                build_binary()
            }
        }
    }
}

def build_binary() {
    try {
        sh "mkdir -p release/$GOOS/$GOARCH"
        sh "go build -o release/$GOOS/$GOARCH/$BINNAME main.go"
    } catch (err) {
        echo "Caught: ${err}"
        currentBuild.result = 'FAILURE'
    }
    archiveArtifacts artifacts: 'release/**', fingerprint: true
    script {
        stash includes: 'release/**', name: 'release'
    }
}
