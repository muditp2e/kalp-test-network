#!/bin/bash

export FABRIC_CA_CLIENT_HOME=$PWD
export USER="clientapp2"
export USER_PWD="clientapppw2"
export HLF_POC1_IP="localhost"
export HLF_POC2_IP="localhost"
export TLS_CA_URL="$HLF_POC1_IP"
export TLS_PORT="7080"
export RELATIVE_PATH_TO_TLS_ROOT_CERT="tls-root-cert/tls-ca-cert.pem"
export TLS_ADMIN_MSP_DIR="tls-ca/tls-admin/msp"

#register
# ./fabric-ca-client register --id.name $USER --id.secret $USER_PWD -u https://"$TLS_CA_URL":"$TLS_PORT"  --tls.certfiles $RELATIVE_PATH_TO_TLS_ROOT_CERT --mspdir $TLS_ADMIN_MSP_DIR

#fabric-ca-client register --caname ca-org1 --id.name $USER --id.secret $USER_PWD -u https://"$TLS_CA_URL":"$TLS_PORT"  --tls.certfiles "${PWD}/organizations/fabric-ca/org1/ca-cert.pem"  --tls.certfiles $RELATIVE_PATH_TO_TLS_ROOT_CERT --mspdir $TLS_ADMIN_MSP_DIR

fabric-ca-client register --caname ca-org1 --id.name $USER --id.secret $USER_PWD --tls.certfiles "./organizations/fabric-ca/org1/ca-cert.pem"  

# --mspdir $TLS_ADMIN_MSP_DIR

# export USER="clientapp1"
# export USER_PWD="clientapppw1"
# export TLS_CA_URL="$HLF_POC1_IP"
# export TLS_PORT="7080"
# export RELATIVE_PATH_TO_TLS_ROOT_CERT="tls-root-cert/tls-ca-cert.pem"
# export SIGNING_PROFILE="tls"
# export USER_MSP_DIR="tls-ca/clientapp/msp"
# export TLS_CA_CSR_HOSTS="localhost,clientapp1.p2epl.com,127.0.0.1,$HLF_POC1_IP,$HLF_POC2_IP"
# export c="IN"
# export st="Uttar Pradesh"
# export l="7th Floor FC-19 Sec-16A FilmCity Noida"
# export o="P2E Pro Pvt. Ltd."

# #enroll
# # ./fabric-ca-client enroll -u https://"$USER":"$USER_PWD"@"$TLS_CA_URL":"$TLS_PORT" --enrollment.profile $SIGNING_PROFILE  --csr.hosts $TLS_CA_CSR_HOSTS   --tls.certfiles $RELATIVE_PATH_TO_TLS_ROOT_CERT  --mspdir $USER_MSP_DIR --csr.names "C=$c,ST=$st,L=$l,O=$o"

# fabric-ca-client enroll -u https://"$USER":"$USER_PWD"@"$TLS_CA_URL":"$TLS_PORT" --caname ca-org1 --enrollment.profile $SIGNING_PROFILE --csr.hosts $TLS_CA_CSR_HOSTS    --tls.certfiles $RELATIVE_PATH_TO_TLS_ROOT_CERT  --mspdir $USER_MSP_DIR --csr.names "C=$c,ST=$st,L=$l,O=$o"

# # -M "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls"

# #fabric-ca-client enroll -u https://"$USER":"$USER_PWD"@"$TLS_CA_URL":"$TLS_PORT" --caname ca-org1 -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/ca-cert.pem"  --enrollment.profile $SIGNING_PROFILE  --mspdir $USER_MSP_DIR 



# # fabric-ca-client enroll -u https://user1:user1pw@localhost:7054 --caname ca-org1 -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/ca-cert.pem"

# #rename private key
# mv $USER_MSP_DIR/keystore/* $USER_MSP_DIR/keystore/key.pem

