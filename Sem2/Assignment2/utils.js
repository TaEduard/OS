const minimist = require('minimist')

const mandatoryCheck = (mandatArgs, args) => {
  for (const arg in mandatArgs)
    if ((args[mandatArgs[arg]] === 'undefined') || (args[mandatArgs[arg]] === '')) {
      console.log(`-${mandatArgs[arg]} is a mandatory argument`)
      console.log(`Use -h for help`)
      return 0
    }
  return 1
}
const mandatoryIfPopulated = (mandatArgs, args) => {
  for (let arg = 0; arg < mandatArgs.length - 1; arg++)
    if ((args[mandatArgs[arg]] === '') && (args[mandatArgs[arg + 1]] !== '') ||
      (args[mandatArgs[arg]] !== '') && (args[mandatArgs[arg + 1]] === '')) {
      console.log(`${mandatArgs} are mandatory arguments if one of them is populated`)
      console.log(`Use -h for help`)
      return 0
    }
  return 1
}
const helpMessage = (filename) => {
  console.log(`
  Usage: minitool -f path/to/file -x -r -p oneGoodPassword
  
  Mandatory:
  -f                       file
  
  One Operation:
  -e                       encrypt
  -d                       decrypt
  -c                       compress
  -x                       decompress
  
  Optional:
  -r                       remove initial file after operation
  -p                       password

`)
  return 0
}

module.exports = {
  helpMessage,
  minimist,
  mandatoryIfPopulated,
  mandatoryCheck
}