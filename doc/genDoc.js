const CP = require('child_process');
const PATH = require('path');
const FS = require('fs');
const IGNORE = [
  'doc',
  'test',
  'unused',
  'experimental',
  '.git'
];

const main = async () => {
  const doc = await parseDir(PATH.resolve(__dirname, '../'));
  console.log(doc);
}
setImmediate(main);

const parseDir = async (rootDir, workDir = "") => {
  const tail = await checkForOtherDir(rootDir, workDir);

  const stdout = await execPromise('go doc -u -all', PATH.resolve(rootDir, workDir));
  const parts = stdout.trim().split(/\n\n\n(?=[A-Z]+\n)/);

  const package = parsePackage(parts.findPop(a => a.startsWith('package')));
  const constants = await parseConstants(parts.findPop(a => a.startsWith('CONSTANTS')), package, rootDir, workDir);
  const functions = await parseFunctions(parts.findPop(a => a.startsWith('FUNCTIONS')), package, rootDir, workDir);
  const types = await parseTypes(parts.findPop(a => a.startsWith('TYPES')), package, rootDir, workDir);
  if(parts.length) {
    FS.writeFileSync(PATH.resolve(__dirname, './unknown.json'), JSON.stringify(parts));
    throw Error(`unknown types in "${PATH.resolve(rootDir, workDir)}"\nplease share "${PATH.resolve(__dirname, './unknown.json')}" to help improve this package`);
  }
  tail.push(
    ...constants,
    ...functions,
    ...types,
  );
  return tail;
}

Array.prototype.findPop = function(callback) {
  const obj = this.find(callback);
  if(obj === undefined) return null;
  this.splice(this.indexOf(obj) , 1);
  return obj;
}

const checkForOtherDir = async (rootDir, workDir) => {
  const tail = [];
  const files = FS.readdirSync(PATH.resolve(rootDir, workDir));

  for(const file of files) {
    const relativePos = PATH.relative(rootDir, PATH.resolve(rootDir, workDir, file));

    const isIgnored = IGNORE.some(pos => {
      return PATH.resolve(rootDir, pos) === PATH.resolve(rootDir, relativePos);
    });
    if(isIgnored) continue;
    const stat = FS.statSync(PATH.resolve(rootDir, relativePos));
    if(stat.isDirectory()) {
      const d = await parseDir(rootDir, relativePos);
      tail.push(...d);
    };
  }
  return tail;
}

const parsePackage = (content) => {
  if(!content) return null;
  return content.match(/package ([a-zA-Z0-9_]+)/)[1];
}
const parseConstants = async (content, package, rootDir, workDir) => {
  if(!content) return [];
  const constants = content.split('\n\n').slice(1);
  const parsedConstants = constants.map(async a => {
    const delimiter = a.indexOf('\n');

    const code = a.substr(0, delimiter);
    const definition = code.match(/(const\s)?([a-zA-Z0-9_]+)(\s[a-zA-Z0-9]+)? = ([\s\S]+)/);
    const { file, line } = await findFileAndLine(rootDir, workDir, code);

    return {
      package,
      file: PATH.relative(rootDir, PATH.resolve(rootDir, workDir, file)),
      line,
      code,
      exported: (/^[A-Z]/).test(definition[1]),
      name: nullsaveTrim(definition[1]),
      type: nullsaveTrim(definition[2]),
      value: nullsaveTrim(definition[3]),
      description: nullsaveTrim(a.substr(delimiter)),
    }
  });
  return await Promise.all(parsedConstants);
}
const parseFunctions = async (content, rootDir, workDir) => {
  return [];
}
const parseTypes = async (content, rootDir, workDir) => {
  return [];
}

const findFileAndLine = async (rootDir, workDir, code) => {
  const stdout = await execPromise(`grep -n -H "${code}" *.go`, PATH.resolve(rootDir, workDir));

  const file = workDir + stdout.split(':')[0];
  const line = stdout.split(':')[1];
  return { file, line };
}

const nullsaveTrim = (arg) => {
  if(typeof arg === 'string' && arg.length) return arg.trim();
  return null;
}

const execPromise = (cmd, cwd) => new Promise((resolve, reject) => {
  CP.exec(cmd, { cwd }, (err, stdout, stderr) => {
    if(err) reject(err);
    else resolve(stdout);
  });
});
