const utils = require('./utils')
const fs = require('fs')
const encrpytion = require('./encrpytion')
const LZString = require('lzjs')
var SnappyJS = require('snappyjs')
password = 'wzZfI5C3cfs3uhGPnV1GX7YtzX3lgPtu'

let args = utils.minimist(process.argv.slice(2), {
  strings: ['remove', 'decompress', 'decrypt', 'password', 'compress', 'encrypt', 'help', 'file',],
  alias: {
    f: 'file',
    c: 'compress',
    x: 'decompress',
    e: 'encrypt',
    d: 'decrypt',
    p: 'password',
    r: 'remove',
    h: 'help'
  },
  default: {
    f: '',
    c: false,
    e: false,
    d: false,
    x: false,
    r: false,
    p: '',
    h: false
  }
})
if (args.help === true || args.h === true) {
  return utils.helpMessage("File")
}
if (utils.mandatoryCheck(["f"], args) === 0) return 0

if (args.p !== '') {
  password = args.p
}
if (!fs.existsSync(args.f)) {
  console.log(`File not found!`)
  return 1
}

if (args.e) {
  let file = fs.readFileSync(args.f)
  encryptedFile = encrpytion.encrypt(file, password)
  fs.writeFileSync(args.f + `.enc`, encryptedFile)
  if (args.r)
    fs.unlink(args.f, (err) => {
      if (err) throw err
    })
  console.log(`File: "${args.f}" successfully encrypted!`)
  return 0
}


if (args.d) {
  let file = fs.readFileSync(args.f)
  if (args.r)
    fs.unlink(args.f, (err) => {
      if (err) throw err
    })
  args.f = args.f.replace(`.enc`, '')
  decryptedFile = encrpytion.decrypt(file.toString(), password)
  fs.writeFileSync(args.f, decryptedFile)
  console.log(`File: "${args.f}" successfully decrypted!`)
  return 0
}

if (args.c) {
  let file = fs.readFileSync(args.f)
  var compressed = SnappyJS.compress(file);
  fs.writeFileSync(args.f + `.comp`, compressed)
  if (args.r)
    fs.unlink(args.f, (err) => {
      if (err) throw err
    })
  console.log(`File: "${args.f}" successfully compressed!`)
  return 0
}

if (args.x) {
  let file = fs.readFileSync(args.f)
  if (args.r)
    fs.unlink(args.f, (err) => {
      if (err) throw err
    })
  args.f = args.f.replace(`.comp`, '')
  var decompressed = SnappyJS.uncompress(file);
  fs.writeFileSync(args.f, decompressed)
  console.log(`File: "${args.f}" successfully decompressed!`)
  return 0
}
