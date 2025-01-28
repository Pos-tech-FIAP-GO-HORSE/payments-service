# Definindo a variável do arquivo ZIP para reutilização
ZIP_FILE=function.zip

# Terraform commands
terraform:
	terraform init
	terraform plan
	terraform apply

# Deploy commands
deploy:
	cd cmd && GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
	zip -r $(ZIP_FILE) cmd internal go.mod go.sum
	aws s3 cp function.zip s3://payments-service/function.zip --region us-east-1
	aws lambda update-function-code --function-name payments-service-lambda --s3-bucket payments-service --s3-key function.zip --region us-east-1
