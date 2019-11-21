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

const parseDir = async (rootDir, workDir = "") => {
  const tail = await checkForOtherDir(rootDir, workDir);

  const stdout = await execPromise('go doc -u -all', PATH.resolve(rootDir, workDir));
  const parts = stdout.trim().split(/\n\n\n(?=[A-Z]+\n)/);

  const package = PARTS.find(a => a.startsWith('package')) || null;
  const constants = await parseConstants(parts.find(a => a.startsWith('CONSTANTS')), package, rootDir, workDir);
  const functions = await parseFunctions(parts.find(a => a.startsWith('FUNCTIONS')), package, rootDir, workDir);
  const types = await parseTypes(parts.find(a => a.startsWith('TYPES')), package, rootDir, workDir);
  tail.push(
    ...constants,
    ...functions,
    ...types,
  );
  return tail;
}
setImmediate(parseDir, PATH.resolve(__dirname, '../'));

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

const parseConstants = async (content, package, rootDir, workDir) => {
  if(!content) return [];
  const constants = content.split('\n\n').slice(1);
  return constants.map(async a => {
    const delimiter = a.indexOf('\n');

    const code = a.substr(0, delimiter);
    const definition = code.match(/(const\s)?([a-zA-Z0-9_]+)(\s[a-zA-Z0-9]+)? = ([\s\S]+)/);
    const { file, line } = await findFileAndLine(rootDir, workDir, code)

    return {
      package,
      file: PATH.relative(rootDir, PATH.resolve(rootDir, workDir, file)),
      line,
      code,
      exported: definition[1].test(/^[A-Z]/),
      name: nullsaveTrim(definition[1]),
      type: nullsaveTrim(definition[2]),
      value: nullsaveTrim(definition[3]),
      description: nullsaveTrim(a.substr(delimiter)),
    }
  });
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
