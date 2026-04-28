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
                go test $(go list ./... | findstr /V tests/functional)
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
                bat '''
                cd PaymentService
                docker build -t payment-service:latest .
                cd ..
                '''
            }
        }

        // ================= PICKUP =================
        stage('Build Pickup Service') {
            steps {
                bat '''
                cd PickupService
                docker build -t pickup-service:latest .
                cd ..
                '''
            }
        }

        // ================= PAYMENT FUNCTIONAL =================
        stage('Functional Test Payment') {
            steps {
                bat '''
                start /b go run PaymentService/main.go
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
                start /b go run PickupService/main.go
                timeout /t 5

                curl -X GET http://localhost:8082/pickup
                '''
            }
        }

        stage('Push Image') {
            steps {
                bat '''
                docker tag payment-service:latest nadzalla/payment-service:latest
                docker push nadzalla/payment-service:latest
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
