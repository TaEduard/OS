

const aesjs = require('aes-js')

const normalizePassword = (password) => {
  const pwLength = password.length
  const chars = []
  const splitChars = () => {
    for (let i = 0; i < password.length; i++) {
      chars.push(password.charAt(i))
    }
  }
  if (pwLength < 32) {
    splitChars()
    for (let i = chars.length; i < 32; i++) {
      chars.push('0')
    }
    return chars.join('')
  } else if (pwLength > 32) {
    splitChars()
    return chars.slice(0, 32).join('')
  }
  return password
}

function encrypt(text, password) {
  try {
    password = aesjs.utils.utf8.toBytes(normalizePassword(password))
    const aesCounter = new aesjs.ModeOfOperation.ctr(password, new aesjs.Counter(5))
    const encryptedBytes = aesCounter.encrypt(text)
    let encryptedText = aesjs.utils.hex.fromBytes(encryptedBytes)
    encryptedText = `SHA256 encryption by TaEd\n` + encryptedText
    return encryptedText
  } catch (error) {
    console.log(error)
    return null
  }
}

function decrypt(encryptedText, password) {
  try {
    password = aesjs.utils.utf8.toBytes(normalizePassword(password))
    encryptedText = encryptedText.replace('SHA256 encryption by TaEd\n', '')
    const encryptedBytes = aesjs.utils.hex.toBytes(encryptedText)
    const aesCounter = new aesjs.ModeOfOperation.ctr(password, new aesjs.Counter(5));
    const decryptedBytes = aesCounter.decrypt(encryptedBytes)
    return decryptedBytes
  } catch (error) {
    console.log(error)
    return null
  }
}
module.exports = {
  encrypt,
  decrypt
}