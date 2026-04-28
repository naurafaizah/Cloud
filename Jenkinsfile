pipeline {
    agent any

    environment {
        IMAGE_NAME = "payment-service"
        IMAGE_TAG = "latest"
    }

    stages {

        stage('Checkout Repo') {
            steps {
                deleteDir()
                git branch: 'main', url: 'https://github.com/naurafaizah/Cloud.git'
            }
        }

        stage('Unit Test') {
            steps {
                bat '''
                go list ./... | findstr /V functional > packages.txt
                for /f %i in (packages.txt) do go test %i
                '''
            }
        }

        stage('Vet') {
            steps {
                bat 'go vet ./...'
            }
        }

        // ================= PAYMENT =================
        stage('Build Docker Image') {
            steps {
                dir('PaymentService') {
                    bat 'docker build -t payment-service:latest .'
                }
            }
        }

        // ================= PICKUP =================
        stage('Build Pickup Service') {
            steps {
                dir('PickupService') {
                    bat 'docker build -t pickup-service:latest .'
                }
            }
        }

        // ================= PAYMENT FUNCTIONAL =================
        stage('Functional Test Payment') {
            steps {
                bat '''
                cd PaymentService
                start /b go run .
                cd ..
                ping 127.0.0.1 -n 6 > nul

                curl -X POST http://localhost:8081/payment ^
                -H "Content-Type: application/json" ^
                -d "{\\"amount\\":10000,\\"paid\\":10000}"
                '''
            }
        }

        // ================= PICKUP FUNCTIONAL =================
        stage('Functional Test Pickup') {
            steps {
                bat '''
                cd PickupService
                start /b go run .
                cd ..
                ping 127.0.0.1 -n 6 > nul

                curl -X POST http://localhost:8082/pickup ^
                -H "Content-Type: application/json" ^
                -d "{\\"order_id\\":\\"ORD1\\",\\"payment_status\\":\\"paid\\",\\"weight\\":2}"
                '''
            }
        }

        stage('Push Image') {
            steps {
                bat 'docker images'
                bat '''
                docker tag payment-service:latest naurafaizah/payment-service:latest
                docker push naurafaizah/payment-service:latest
                '''
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                bat 'kubectl apply -f k8s/'
            }
        }

        stage('Verify Deployment') {
            steps {
                bat 'kubectl get pods && kubectl get svc'
            }
        }
    }
}
