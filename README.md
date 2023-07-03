# ChefBook Backend Encryption Service

This microservice is responsible for encryption keys storing and management.

## Encryption Process

ChefBook encrypted recipes doesn't have to be non-sharable.
But sharing recipe must give access to this recipe only.
So, to achieve this, we can use simple algorithm:

<p align="center">
    <img src="./img/recipe_encryption.png"/>
</p>

### Vault

1. When user creates encrypted vault, client generates RSA keypair.
2. To synchronize data between devices, user have to store vault keys remotely. 
Private key must be accessible only for user.
3. For this purpose client generates AES based on user passphrase.
4. RSA encrypted with AES is uploaded to the server.

### Recipe

1. For every encrypted recipe client generates random AES.
2. AES encrypted with RSA Public Key is uploaded to the server.

So chain's first link is user passphrase. Each link can be changed,
which gives flexibility in sharing data.

## Recipe Key Sharing Flow

<p align="center">
    <img src="./img/recipe_key_sharing.png"/>
</p>