# Bartering protocol

Ralph HATOUM - PR INSA Lyon 2023-2024

Projet still under development

## General description

This project aims at building and IPFS pinning overlay to ensure data replication and efficient storage use over a set of nodes. Unlike other pinning services such as Filecoin that rely on a cryptocurrency and blockchain to allow nodes to reserve storage space on each other, we aim to implement a bartering system.

## Quickstart for local dev
This is a guide to launch the protocol locally.

### Running IPFS
Bartering relies on IPFS. So, we need to have an IPFS client running on your computer. We therefore need to install IPFS. There are several ways to do this, here are two : 
* [IPFS Desktop App](https://docs.ipfs.tech/install/ipfs-desktop/), which will have a GUI 
* [IPFS Kubo](https://docs.ipfs.tech/install/command-line/), which is a CLI tool

I advise using the second option (Kubo) as this is what the protocol was built with. Follow the instructions to install it and, once done, verify your installation with the following command :

```
ipfs version
```
which should ouput the installed IPFS version if everything went well. 

Once this is successfully done, we can run IPFS with the following command : 
```
ipfs daemon
```
IPFS is now running your machine ! Keep the terminal in which you started this command open. If you want to reuse the same terminal, you can also launch the command and run it in the background :
```
ipfs daemon &
```

### Launching the bartering bootstrap

The protocol relies on a bootstrap system. Basically, whenever a new node connects to the network, it will contact the bootstrap to get a list of peers. We have to launch this bootstrap first. In a new terminal, run the following commands, from the root of this repository :
```
cd bartering-node
go run bootstrap.go 127.0.0.1
```

### Launching bartering

You can now run the protocol, again in a new terminal (we need both IPFS and the bootstrap running still ! ). Enter the following command :
```
go run main.go 127.0.0.1
```
Your node is now up and running. By default, it will watch and send to the network any file that you create or add in the data/ folder. Try it !

## Libraries description

### api-ipfs
Functions used to interact with the IPFS daemon

### bartering-api
Functions related to the bartering process

### bootstrap-connect
Functions used to connect and interact with the network's bootstrap node

### bootstrap-node
Daemon to run on the bootstrap node

### functions
General funtions used in the node daemon

### failure-simulation
Functions to make nodes fail - used to test the protocol's behavior

### peers-connect
Functions used to interact with other peers in the network

### storage-requests
Functions and datatypes used to represent storage requests and deal with them

### storage-testing
Functions used to request proof of storage from other peers

### datastructures
All datastructures used declared in a single library to prevent circular dependencies

### utils
General functions to perform operations such as printing lists or handling errors
