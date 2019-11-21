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

const parseDir = (dir) => {
  const files = FS.readdirSync(dir);
  for(const f of files) {
    if(IGNORE.includes(f)) continue;
    const stat = FS.statSync(PATH.resolve(dir, f));
    if(stat.isDirectory()) parseDir(PATH.resolve(dir, f));
  }
  CP.exec('go doc -u -all', {cwd: dir}, (err, stdout, stderr) => {
    const PARTS = stdout.trim().match(/(package([\s\S](?!CONSTANTS|FUNCTIONS|TYPES))*)?(CONSTANTS([\s\S](?!FUNCTIONS|TYPES))*)?(FUNCTIONS([\s\S](?!TYPES))*)?(TYPES[\s\S]*)?/);
    console.log(dir, PARTS);
    console.log({
      line: null,
      file: null,
      package: PARTS[1],
      constants: PARTS[2],
      functions: PARTS[2],
      types: PARTS[3],
      private: false,
    });
  });
}
parseDir(PATH.resolve(__dirname, '../'));

const findFileAndLine = (cmd, dir) => {
  CP.exec(`grep -n -H "${cmd}" *.go`, {cwd: dir}, (err, stdout, stderr) => {
    console.log('findFileAndLine', {stdout});
  })
}
