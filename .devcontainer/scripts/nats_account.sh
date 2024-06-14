#!/bin/bash
 # script to dump creds for use in our app

sudo chown -R vscode /nsc

 echo "Dumping NATS user creds file"
 nsc generate creds -a serverservice -n USER > /tmp/user.creds

 echo "Dumping NATS sys creds file"
 nsc generate creds -a SYS -n sys > /tmp/sys.creds
