pipeline {
    agent {
        docker { image 'hashicorp/terraform:0.15.5'
                 label 'docker'
                 args '--entrypoint='
                 }
    }
    environment {
        // this slows down the process but makes it available to all stages so I can use it in the message. Move it to the plan stage in order to speed it up but it'll break the message 
                OUT = """${sh(
                        returnStdout: true,
                        script: 'terraform plan'
                    )}""" 
    }
    stages {
        stage('test') {
            steps {
                sh 'ls -la'
                sh 'echo "here come the terratest tests"'
            }
        }
        stage('validate') {
            steps {
                sh 'terraform init'
                sh 'terraform validate'
            }
        }
        stage('plan') {
            steps {
                sh 'terraform plan'
            }
        }
        stage('Confirmation') {
            input {
                message "Confirm the changes? "
                ok "Submit"
                submitter "franco,federico,alexis,hugo"
                submitterParameter "WHO"
                parameters {
                    string(name: 'ANSWER', description: 'the script only accepts "yes" to procede')
                }
            }
            steps {
                sh 'if [[ "${ANSWER}" == "yes" ]] ; then echo "Confirmation was provided by ${WHO} at $(date -Iseconds)" ; else exit 2; fi'
            }
        }
        stage('apply') {
            steps {
                sh 'terraform apply -auto-approve'
            }
        }
    }
}
