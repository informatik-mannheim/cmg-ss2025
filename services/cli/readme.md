## Start CLI

1. Navigate to `services/cli` in your terminal 
run `make build-image`

2. run `make run`. You're now prompted to type (`>`)

## Use CLI

At all times, type `help` or `<command>` for instructions.

1. Login ``login --<your_secret>``
2. Create a job ``create-job --job-name <value> --creation-zone <value> 
--image-name <value> --image-version <value> --parameters <value>``
3. Get job outcome `` get-job-outcome --id <value>``
4. Get job `get-job --id <value>`