# Shady

Very basic ransomware implementation that encrypts files, sends encryption keys to the API server, requests Monero from the user. Uses [MoneroPay](https://gitlab.com/moneropay/moneropay) to track payments sent by the victim. Provides an endpoint for the user to download keys to decrypt their files if the payment has arrived.

That's right. The only way I will ever touch Powershell is if it's gonna be straight up malware. Universities need to stop teaching proprietary software that nobody uses. Teach a [real](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html) shell, not a C# interpreter.

## Server
### POST /encrypt
### GET /decrypt/{id}
### POST /callback/{id}

## Client
