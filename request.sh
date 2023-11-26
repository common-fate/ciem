go run cmd/cli/main.go entities put --file entities.json
go run cmd/cli/main.go policy apply --file example.cedar --id example
go run cmd/cli/main.go policy apply --file review.cedar --id review
# go run cmd/cli/main.go request gcp project --id demo-project --role roles/editor