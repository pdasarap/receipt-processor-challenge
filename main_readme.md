Following are the Steps to run this webservice

Clone my repository:

`git clone https://github.com/pdasarap/receipt-processor-challenge.git`

Build the Docker image using following command in the terminal

`docker build -t receipt_processor .`

Run the Docker container:

`docker run -dp 8080:8080 receipt_processor`

Check the APIs using Postman or cUrl using command or any online curl editor

POST ON `http://localhost:8080/receipts/process`

A receipt ID is generated in response - Copy this ID and use in GetPoints URL. Here's an example

GET ON `http://localhost:8080/receipts/8bf92b8e-76e6-4b94-9a79-1247a3207008/points`