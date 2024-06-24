# go-telegraf-gui
This is a very simple web page to configure [telegraf OPCUA input plugin](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/opcua). Check the link for full documentation.
Keep in mind that this project simply makes adding nodes easier, and does not implement any OPC UA Browser functionality.

## Installation
Simply download this project, either use the big green <> Code button above (unzip if needed) or: 
```sh 
git clone https://github.com/Rareshp/go-telegraf-gui
```

Additionaly, you need to have go instealled. See [here](https://go.dev/doc/install). I used go 1.22.1 

## Usage 
First, navigate to where you downloaded the project, then run the server:
```sh 
go run main.go
```

Now simply navigate to `http://localhost:8080`. Follow the wizard by clicking the blue buttons.
The site will then generate text you can simply copy in your telegraf.conf
