# End to end test Refresco

This e2e test is meant to prevent regression in the austria juice email pipeline in case changes are made to the mail server, the email pipeline or associated configuration.

## Test dependencies
As an e2e test, dependencies are mocked as much as possible. The import API is mocked with a local http server that receives the POST message and can perform assertions on it.
A local SMTP server is started to check error messages. This email server is configured in the docker-compose file in the root directory.

## Test scenario: success
Test setup
1. A HTTP server is started on localhost:80
2. Pipeline is started
Test execution
3. The pipeline is triggered by sending
Test assertions
4. The incoming message on the local HTTP server is checked for the expected body that should be sent to the import API.
6. The SMTP server is checked for not having error messages
Test teardown
5. The pipeline is stopped


