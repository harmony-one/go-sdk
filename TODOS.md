After talking with John, clear that the tasks are:

1. Will need to make an Account with local alias as a nice to have but not really needed
   1.5) Can do struct embedded to add the alias feature

2. Also need to accept address with plain 0x, nice to have with bech32 address

3. Do a signRawTransaction, send it off over the RPC correctly, do expose the shard from and to IDs,
   already have example from JS SDK which is great

4. Need to have travis testing

5. Need to have e2e testing examples

6. Big problem on nil reference on the Account that you can't find in local store, need to stop that
   earlier before a panic goes off
