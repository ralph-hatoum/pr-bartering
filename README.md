# Bartering protocol

Ralph HATOUM - PR INSA Lyon 2023-2024

Projet still under development

## General description

This project aims at building and IPFS pinning overlay to ensure data replication and efficient storage use over a set of nodes. Unlike other pinning services such as Filecoin that rely on a cryptocurrency and blockchain to allow nodes to reserve storage space on each other, we aim to implement a bartering system.

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

### peers-connect
Functions used to interact with other peers in the network

### utils
General functions to perform operations such as printing lists or handling errors