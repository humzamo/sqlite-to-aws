# sqlite-to-aws

A small script to take an SQLite database and upload the contents to a specified AWS S3 bucket as JSON.

## Usage

Run the script by running the following command:

```
go run main.go <database_name> <table_name>
```

where `<database_name>` and `<table_name>` are the respective database file and table names with the data to be uploaded.
Ensure that the database file is in the `cmd/main` location.
The `CreateSampleTable()` function in the `data` package can also be used to create a sample SQLite database for the script.

## Requirements

- Go
- Internet connection to use the AWS SDK
- AWS credentials for the specified bucket (see the [AWS documentation](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-credentials) for more details)

## Assumptions

- This script assumes a relatively simple table schema as defined in `models.go` get the row data. For another schema, these structs can be updated in the `data` package.
- The details for the region, bucket, and path have been fixed in the `awsupload` package. These have been set using resources made in my AWS account for testing purposes.
- The script assumes there are AWS credentials available on the local machine. These could be configurable elsewhere (e.g. using GitHub secrets).

## Further Improvements

Some further improvements could be made given additional business logic and time, such as:

- A wider array of inputs could be used in the command-line arguments (e.g. to specify an S3 bucket).
- The format for the table data and JSON to be uploaded has been kept simple and could be expanded if required.
- A more sophisticated Golang package to create command line scripts could be used (e.g. [Cobra](https://cobra.dev)).
- Interfaces could be used in the `awsupload` and `data` packages; this would also help with unit and end-to-end testing of the script.
