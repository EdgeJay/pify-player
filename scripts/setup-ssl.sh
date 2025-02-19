#!/bin/bash

data_path="$PWD/certs"

helpFunction()
{
   echo ""
   echo "Usage: $0 -d domain.name -e email_address"
   echo -e "\t-d Domain name to generate SSL certs for"
   exit 1 # Exit script after printing help
}

while getopts "d:" opt
do
   case "$opt" in
      d ) domain="$OPTARG" ;;
      ? ) helpFunction ;; # Print helpFunction in case parameter is non-existent
   esac
done

# Print helpFunction in case parameters are empty
if [ -z "$domain" ]
then
   echo "Some or all of the parameters are empty";
   helpFunction
fi

# Begin script in case all parameters are correct
echo "Domain: $domain"

if ! [ -x "$(command -v mkcert)" ]; then
  echo 'Error: mkcert is not installed. Run `brew install mkcert`' >&2
  exit 1
fi

if [ -d "$data_path" ]; then
  read -p "Existing data found for $domain. Continue and replace existing certificate? (y/N) " decision
  if [ "$decision" != "Y" ] && [ "$decision" != "y" ]; then
    exit
  fi
fi

mkdir -p $data_path

mkcert -cert-file "$data_path/$domain.pem" -key-file "$data_path/$domain.key.pem" $domain localhost
