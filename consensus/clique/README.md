# Run Private nodes

To run the codebase of feat/ren branch follow the steps

### Step 1
Install Go. Visit https://go.dev/doc/install for OS specific info.



### Step 2
Clone go-ethereum repo using
```
git clone https://github.com/gdsoumya/go-ethereum.git
```
and switch to feat/ren branch

### Step 3
Run
1.  ```
    make geth
    ```
2.	```
    make bootnode
    ```
3.  ```
    bash ./genesis/bootnode/start.sh
    ```
### Step 4 (opt)
Create a genesis folder in root directory
```
mkdir genesis
```
and 3 sub folders  inside genesis
```
mkdir node1 node2 node3
``` 

### Step 5(opt)
Create a genesis file in each of the nodes
```
{
"config": {
"chainId": 826978,
"homesteadBlock": 0,
"eip150Block": 0,
"eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
"eip155Block": 0,
"eip158Block": 0,
"byzantiumBlock": 0,
"constantinopleBlock": 0,
"petersburgBlock": 0,
"istanbulBlock": 0,
"clique": {
"period": 10,
"dnr": "0xe37D748D059eFCd8d92834548273Ce673dDeC691",
"epochBlock": 32048511,
"api" : "https://kovan.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
"initialValidators": ["0x5b5496828980b6bbd15579d348eaa1aa9b6d8cb3"]
}
},
"nonce": "0x1E9057F",
"timestamp": "0x628731a1",
"extraData": "0x00000000000000000000000000000000000000000000000000000000000000005b5496828980b6bbd15579d348eaa1aa9b6d8cb30000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
"gasLimit": "0x47b760",
"difficulty": "0x1",
"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
"coinbase": "0x0000000000000000000000000000000000000000",
"alloc": {
"0000000000000000000000000000000000000000": {
"balance": "0x1"
},
"91f79893e7b923410ef1aeba6a67c6fab07d800c": {
"balance": "0x200000000000000000000000000000000000000000000000000000000000000"
}
},
"number": "0x0",
"gasUsed": "0x0",
"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
"baseFeePerGas": null
}
```
### Step 6 (if 4 and 5 steps were skipped)
Download the genesis folder and include it in the directory
[https://drive.google.com/file/d/1npGSHhRr1_7jAjEbHrtZ6idGPTH_EJjw/view?usp=sharing](https://drive.google.com/file/d/1npGSHhRr1_7jAjEbHrtZ6idGPTH_EJjw/view?usp=sharing "https://drive.google.com/file/d/1npGSHhRr1_7jAjEbHrtZ6idGPTH_EJjw/view?usp=sharing")

### Step 7
Head to remix.ethereum.org and then navigate to deploy and run
1. Select "Injected Web3" as environment
2. Switch to Kovan Testnet on metamask
3. Provide at address as 0xe37D748D059eFCd8d92834548273Ce673dDeC691 and click the button
4. Use epoch function to produce an epoch and copy the latest block number
5. Update ren.json with latest block number and nonce as the hexadecimal of the blocknumbbr as 0x[nonce]

### Step 8
Run
```
bash ./genesis/bootnode/start.sh
```
### Step 9
Run nodes on your network
1. cd into node1 folder and run
```
../../build/bin/geth --datadir ./ init ren.json && ../../build/bin/geth --datadir ./ --syncmode 'full' --port 30311 --http --http.addr '0.0.0.0' --http.corsdomain "*" --http.port 8502 --http.api 'personal,clique,eth,net,web3,txpool,miner' --bootnodes 'enode://cec4db42fe455051e343c62bfee312b56fbe43b023be91c67446be283729fd1b11e4b8f193f3a1cc08d43fc21a6f2f6d0d09b9029b76c75a68423b19a4bef904@127.0.0.1:30310' --networkid 826978 --miner.gasprice '1' --allow-insecure-unlock -unlock 5b5496828980b6bbd15579d348eaa1aa9b6d8cb3 --password password.txt --mine
```
2. cd into node2 folder and run
```
../../build/bin/geth --datadir ./ init ren.json && ../../build/bin/geth --datadir ./ --syncmode 'full' --port 30312 --http --http.addr '0.0.0.0' --http.corsdomain "*" --http.port 8503 --http.api 'personal,clique,eth,net,web3,txpool,miner' --bootnodes 'enode://cec4db42fe455051e343c62bfee312b56fbe43b023be91c67446be283729fd1b11e4b8f193f3a1cc08d43fc21a6f2f6d0d09b9029b76c75a68423b19a4bef904@127.0.0.1:30310' --networkid 826978 --miner.gasprice '1' --allow-insecure-unlock -unlock a53a2B40039Ec8986Bfb1170e9f146DbC3e3de83 --password password.txt --mine
```
3. cd into node3 folder and run
```
../../build/bin/geth --datadir ./ init ren.json && ../../build/bin/geth --datadir ./ --syncmode 'full' --port 30313 --http --http.addr '0.0.0.0' --http.corsdomain "*" --http.port 8504 --http.api 'personal,eth,net,web3,txpool,miner,clique' --bootnodes 'enode://cec4db42fe455051e343c62bfee312b56fbe43b023be91c67446be283729fd1b11e4b8f193f3a1cc08d43fc21a6f2f6d0d09b9029b76c75a68423b19a4bef904@127.0.0.1:30310' --networkid 826978 --miner.gasprice '1' --allow-insecure-unlock -unlock 98d958B7862fF36D9f28F2483C68D96Fa11b6031 --password password.txt --mine
```
### Step 10
Head to remix.ethereum.org  then to deploy and run tab
1. Register or deregister nodes and crun epoch to edit signers map
2. Change will be reflected in the private nodes in about 10 block span
